package subscriber

import (
	"strings"
	"time"

	"github.com/ritchieridanko/pasarly/backend/services/user/configs"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Init(cfg *configs.Broker, topic string, l *zap.Logger) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(cfg.Brokers, ","),
		GroupID:        "user-service",
		Topic:          topic,
		MaxBytes:       cfg.MaxBytes,
		CommitInterval: time.Second,
		MaxAttempts:    cfg.MaxAttempts,
	})

	l.Sugar().Infof("✅ [SUBSCRIBER] initialized (topic=%s, brokers=%s)", topic, cfg.Brokers)
	return r
}
