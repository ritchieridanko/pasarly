package publisher

import (
	"strings"

	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Init(cfg *configs.Broker, topic string, l *zap.Logger) *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      strings.Split(cfg.Brokers, ","),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireAll),
		Async:        false,
		BatchTimeout: cfg.Timeout.Batch,
		MaxAttempts:  cfg.MaxAttempts,
	})

	l.Sugar().Infof("✅ [PUBLISHER] initialized (topic=%s, brokers=%s)", topic, cfg.Brokers)
	return w
}
