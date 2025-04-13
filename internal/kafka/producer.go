package kafka

import (
	"cloth-mini-app/internal/config"
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(cfg config.Kafka) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBroker),
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
		Logger:   kafka.LoggerFunc(logf),
	}

	return &Producer{
		w: writer,
	}
}

func (p *Producer) WriteMesage(ctx context.Context, payload []byte) error {
	err := p.w.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})
	if err != nil {
		return err
	}

	return nil
}

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}
