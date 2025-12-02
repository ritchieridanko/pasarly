package subscriber

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/notification/configs"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Subscriber struct {
	config *configs.Broker
	reader *kafka.Reader
	logger *logger.Logger
}

func NewSubscriber(cfg *configs.Broker, r *kafka.Reader, l *logger.Logger) *Subscriber {
	return &Subscriber{config: cfg, reader: r, logger: l}
}

func (s *Subscriber) Listen(ctx context.Context, handler func(context.Context, kafka.Message) error) error {
	for {
		m, err := s.reader.FetchMessage(ctx)
		if err != nil {
			t := s.reader.Config().Topic
			p := s.reader.Config().Partition

			if !isRetryable(err) {
				return fmt.Errorf("failed to fetch message (topic=%s): %w", t, err)
			}

			s.logger.Sugar().Warnf("failed to fetch message (topic=%s, partition=%d): %s", t, p, err.Error())
			continue
		}

		c := context.Background()
		if err := s.process(c, m, handler); err != nil {
			s.logger.Base().Error(
				"PROCESS_FAILED",
				zap.String("topic", m.Topic),
				zap.Int("partition", m.Partition),
				zap.Int64("offset", m.Offset),
				zap.String("key", string(m.Key)),
				zap.String("error_detail", err.Error()),
			)

			continue
		}
		if err := s.commit(ctx, m); err != nil {
			s.logger.Base().Error(
				"COMMIT_FAILED",
				zap.String("topic", m.Topic),
				zap.Int("partition", m.Partition),
				zap.Int64("offset", m.Offset),
				zap.String("key", string(m.Key)),
				zap.String("error_detail", err.Error()),
			)
		}
	}
}

func (s *Subscriber) process(ctx context.Context, m kafka.Message, handler func(context.Context, kafka.Message) error) error {
	var e error
	for attempt := 0; attempt < s.config.MaxAttempts; attempt++ {
		err := handler(ctx, m)
		if err == nil {
			return nil
		}

		e = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, s.config.BaseDelay, attempt); err != nil {
			return fmt.Errorf("failed to process message: %w", err)
		}
	}

	return e
}

func (s *Subscriber) commit(ctx context.Context, m kafka.Message) error {
	var e error
	for attempt := 0; attempt < s.config.MaxAttempts; attempt++ {
		err := s.reader.CommitMessages(ctx, m)
		if err == nil {
			return nil
		}

		e = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, s.config.BaseDelay, attempt); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}
	}

	return fmt.Errorf("failed to commit message: %w", e)
}
