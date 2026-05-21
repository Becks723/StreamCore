package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 运行前确保 RabbitMQ 已启动: make env-up
// go test -bench=. -benchtime=5s ./pkg/mq/

func BenchmarkRabbitMQPublish(b *testing.B) {
	conn, err := amqp.Dial("amqp://streamcore:streamcore@localhost:5672/")
	if err != nil {
		b.Skipf("RabbitMQ not available, skip benchmark: %v", err)
	}
	defer conn.Close()

	sender, err := NewRabbitSender(conn, "benchmark_queue")
	if err != nil {
		b.Fatal(err)
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":  12345,
		"video_id": 67890,
		"action":   1,
		"time":     time.Now().Unix(),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := sender.Send(context.Background(), payload); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRabbitMQConsume(b *testing.B) {
	conn, err := amqp.Dial("amqp://streamcore:streamcore@localhost:5672/")
	if err != nil {
		b.Skipf("RabbitMQ not available, skip benchmark: %v", err)
	}
	defer conn.Close()

	queueName := "benchmark_consume_queue"

	// pre-fill
	sender, _ := NewRabbitSender(conn, queueName)
	payload, _ := json.Marshal(map[string]string{"test": "data"})
	for i := 0; i < b.N+100; i++ {
		sender.Send(context.Background(), payload)
	}

	consumer, err := NewRabbitConsumer(conn, queueName)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg, err := consumer.Receive()
		if err != nil {
			b.Fatal(err)
		}
		consumer.Ack(msg)
	}
}

// 容量测算辅助函数：测出当前环境下的极限吞吐
func TestCapacityReport(b *testing.T) {
	conn, err := amqp.Dial("amqp://streamcore:streamcore@localhost:5672/")
	if err != nil {
		b.Skipf("RabbitMQ not available: %v", err)
	}
	defer conn.Close()

	queueName := "capacity_test_queue"
	sender, _ := NewRabbitSender(conn, queueName)

	payload := []byte(`{"test":"data"}`)
	duration := 3 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	count := 0
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			elapsed := time.Since(start).Seconds()
			qps := float64(count) / elapsed
			fmt.Printf("\n=== RabbitMQ 容量测算 ===\n")
			fmt.Printf("单队列单连接: %.0f msg/s\n", qps)
			return
		default:
			if err := sender.Send(ctx, payload); err == nil {
				count++
			}
		}
	}
}
