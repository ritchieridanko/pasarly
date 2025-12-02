package mailer

import (
	"github.com/ritchieridanko/pasarly/backend/services/notification/configs"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

func Init(cfg *configs.Mailer, l *zap.Logger) *gomail.Dialer {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass)

	l.Sugar().Infof("✅ [MAILER] initialized (host=%s, port=%d, from=%s)", cfg.Host, cfg.Port, cfg.From)
	return d
}
