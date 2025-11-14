package publisher

import (
	"strings"

	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Init(cfg *configs.Broker, topic string, l *zap.Logger) (*kafka.Writer, error) {
	b := strings.Split(cfg.Brokers, ",")

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      b,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireAll),
		Async:        false,
		BatchTimeout: cfg.Timeout.Batch,
		MaxAttempts:  cfg.MaxAttempts,
	})

	l.Sugar().Infof("✅ [PUBLISHER] initialized (topic=%s, brokers=%s)", topic, cfg.Brokers)
	return w, nil
}
