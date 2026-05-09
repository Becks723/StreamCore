package chat

import (
	"context"
	"log"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/hertz-contrib/websocket"
	"github.com/redis/go-redis/v9"
)

const wsPubSubChannel = "ws:push"

type WSService struct {
	hub         *connectionHub
	redisClient *redis.Client
}

var (
	wsServiceOnce sync.Once
	wsServiceInst *WSService
)

func InitWSService(rdb *redis.Client) *WSService {
	wsServiceOnce.Do(func() {
		wsServiceInst = &WSService{
			hub:         newConnectionHub(),
			redisClient: rdb,
		}
		go wsServiceInst.listenRedisPubSub()
	})
	return wsServiceInst
}

func GlobalWSService() *WSService {
	return wsServiceInst
}

func (s *WSService) RegisterClient(uid uint, conn *websocket.Conn) {
	s.hub.register(uid, conn)
}

func (s *WSService) UnregisterClient(uid uint, conn *websocket.Conn) {
	s.hub.unregister(uid, conn)
}

func (s *WSService) PushToUser(uid uint, msgType string, data any) error {
	buf, err := sonic.Marshal(&wsPubSubMsg{
		UID:  uid,
		Type: msgType,
		Data: mustMarshalRaw(data),
	})
	if err != nil {
		return err
	}
	return s.redisClient.Publish(context.Background(), wsPubSubChannel, buf).Err()
}

func (s *WSService) listenRedisPubSub() {
	pubsub := s.redisClient.Subscribe(context.Background(), wsPubSubChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var m wsPubSubMsg
		if err := sonic.Unmarshal([]byte(msg.Payload), &m); err != nil {
			log.Printf("ws pubsub unmarshal failed: %v", err)
			continue
		}
		s.pushToLocal(m.UID, m.Type, m.Data)
	}
}

func (s *WSService) pushToLocal(uid uint, msgType string, data []byte) {
	buf, err := sonic.Marshal(&MsgWrapper{Type: msgType, Data: data})
	if err != nil {
		log.Printf("ws marshal wrapper failed uid=%d: %v", uid, err)
		return
	}
	conn, ok := s.hub.get(uid)
	if !ok {
		return
	}
	conn.writeMu.Lock()
	defer conn.writeMu.Unlock()

	if err = conn.conn.WriteMessage(websocket.TextMessage, buf); err != nil {
		s.hub.unregister(uid, conn.conn)
		_ = conn.conn.Close()
		log.Printf("ws write failed uid=%d: %v", uid, err)
	}
}

func mustMarshalRaw(v any) []byte {
	buf, err := sonic.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return buf
}
