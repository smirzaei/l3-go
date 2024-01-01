package downstream

import (
	"context"
	"io"

	"github.com/smirzaei/13-go/internal/config"
	"go.uber.org/zap"
)

type Client struct {
	l      zap.Logger
	c      *config.Service
	ctx    context.Context
	queuer UpstreamQueuer
	stream io.ReadWriteCloser
}

func NewClient(ctx context.Context, logger zap.Logger, conf *config.Service, queuer UpstreamQueuer, stream io.ReadWriteCloser) *Client {
	c := Client{
		l:      logger,
		c:      conf,
		ctx:    ctx,
		queuer: queuer,
		stream: stream,
	}

	return &c
}

func (c *Client) Serve() error {
	panic("not implemented")
}
