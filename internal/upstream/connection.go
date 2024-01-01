package upstream

import (
	"context"
	"io"
	"net"

	"go.uber.org/zap"
)

type Connection struct {
	l      zap.Logger
	stream io.ReadWriteCloser
}

func NewConnection(logger zap.Logger, host string) (*Connection, error) {
	stream, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	c := Connection{
		l:      logger,
		stream: stream,
	}

	return &c, nil
}

func (c *Connection) Serve(ctx context.Context, queue <-chan Request) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case req := <-queue:
			_ = req
			panic("not implemented")
		}
	}
}
