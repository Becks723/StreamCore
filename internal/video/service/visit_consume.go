package service

import (
	"context"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/mq/model"
	"StreamCore/pkg/mq"
	"github.com/bytedance/sonic"
)

func (s *VideoService) consumeVisit(ctx context.Context) {
	c, err := s.mq.Consumer()
	if err != nil {
		// TODO: log consumer init error
		return
	}
	queueSize := constants.VideoVisitQueueSize
	batchSize := constants.VideoVisitBatchSize
	msgQueue := make(chan *mq.Message, queueSize)
	batch := newVisitBatch(batchSize, c)
	ticker := time.NewTicker(constants.VideoVisitFlushInterval)
	defer ticker.Stop()
	go func() {
		for {
			msg, err := c.Receive()
			if err != nil { // eof
				break
			}
			msgQueue <- msg
		}
	}()
	for {
		select {
		case msg, ok := <-msgQueue:
			if !ok {
				s.flushVisitToDB(context.Background(), batch)
				return
			}
			batch.Add(msg)
			if batch.Length() >= batchSize {
				s.flushVisitToDB(ctx, batch)
			}
		case <-ticker.C:
			if batch.Length() > 0 {
				s.flushVisitToDB(ctx, batch)
			}
		case <-ctx.Done():
			// TODO: log ctx done

			s.flushVisitToDB(context.Background(), batch) // clean batch
			draining := true
			for draining {
				select {
				case msg, ok := <-msgQueue:
					if !ok {
						draining = false
						break
					}
					batch.Add(msg)
					if batch.Length() >= batchSize {
						s.flushVisitToDB(context.Background(), batch)
					}
				default:
					draining = false // queue empty
				}

			}
			s.flushVisitToDB(context.Background(), batch) // final batch
			return
		}
	}
}

func (s *VideoService) flushVisitToDB(ctx context.Context, batch *visitBatch) {
	c, msgs := batch.Flush()
	err := s.db.BatchUpdateVisits(ctx, c)
	if err != nil {
		// TODO: log
		return
	}

	for _, msg := range msgs {
		batch.consumer.Ack(msg)
	}
}

type visitBatch struct {
	c         map[uint]int64 // vid -> visitCount
	batchSize int
	msgs      []*mq.Message
	consumer  mq.Consumer
}

func newVisitBatch(batchSize int, consumer mq.Consumer) *visitBatch {
	return &visitBatch{
		c:         make(map[uint]int64, batchSize),
		batchSize: batchSize,
		consumer:  consumer,
	}
}

func (batch *visitBatch) Add(msg *mq.Message) {
	var ev model.VisitEvent
	err := sonic.Unmarshal(msg.Body, &ev)
	if err != nil {
		batch.consumer.Ack(msg) // throw away err msg
		return
	}
	batch.c[ev.Vid]++
	batch.msgs = append(batch.msgs, msg)
}

func (batch *visitBatch) Length() int {
	return len(batch.msgs)
}

func (batch *visitBatch) Flush() (map[uint]int64, []*mq.Message) {
	c := batch.c
	msgs := batch.msgs
	batch.c = make(map[uint]int64, batch.batchSize)
	batch.msgs = nil
	return c, msgs
}
