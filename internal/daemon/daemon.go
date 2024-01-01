package daemon

import (
	"github.com/smirzaei/13-go/internal/config"
	"go.uber.org/zap"
)

type Daemon struct {
	l zap.Logger
	c *config.Config
}

func NewDaemon(logger zap.Logger, conf *config.Config) *Daemon {
	d := Daemon{
		l: logger,
		c: conf,
	}

	return &d
}
