package daemon

import (
	"context"

	"github.com/smirzaei/13-go/internal/config"
	"github.com/smirzaei/13-go/internal/downstream"
	"github.com/smirzaei/13-go/internal/upstream"
	"go.uber.org/zap"
)

type Daemon struct {
	l *zap.Logger
	c *config.Config
}

func NewDaemon(logger *zap.Logger, conf *config.Config) *Daemon {
	d := Daemon{
		l: logger,
		c: conf,
	}

	return &d
}

func (d *Daemon) Run(ctx context.Context) error {
	d.l.Info("running the deamon")

	upstreamCtx, upstreamCancel := context.WithCancel(context.Background())
	downstreamCtx, downstreamCancel := context.WithCancel(context.Background())
	defer func() {
		d.l.Info("stopping the daemon")

		upstreamCancel()
		downstreamCancel()
	}()

	upstreamPool := upstream.NewPool(d.l, d.c)
	upstreamPool.Start(upstreamCtx)

	downstreamServer := downstream.NewServer(d.l, &d.c.Service, upstreamPool)
	err := downstreamServer.Listen(downstreamCtx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
