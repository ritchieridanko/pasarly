package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
)

type Cookie struct {
	config   *configs.Config
	domain   string
	isSecure bool
	httpOnly bool
}

func NewCookie(cfg *configs.Config, httpOnly bool) *Cookie {
	isProd := NormalizeString(cfg.App.Env) == "prod"

	domain := ""
	if isProd {
		domain = cfg.Server.Host
	}

	return &Cookie{config: cfg, domain: domain, isSecure: isProd, httpOnly: httpOnly}
}

func (u *Cookie) Set(ctx *gin.Context, name, value string, d time.Duration, path string) {
	ctx.SetCookie(name, value, int(d.Seconds()), path, u.domain, u.isSecure, u.httpOnly)
}

func (u *Cookie) Unset(ctx *gin.Context, name, path string) {
	ctx.SetCookie(name, "", -1, path, u.domain, u.isSecure, u.httpOnly)
}
