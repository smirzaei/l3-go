package downstream

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/smirzaei/13-go/internal/config"
	"github.com/smirzaei/13-go/internal/frame"
	"go.uber.org/zap"
)

var (
	ErrFrameRead       = fmt.Errorf("failed to read the whole frame")
	ErrPayloadRead     = fmt.Errorf("failed to read the payload")
	ErrUpstreamFailure = fmt.Errorf("failed to get a response from upstream")
	ErrWrite           = fmt.Errorf("failed to write the response")
)

type Client struct {
	l      zap.Logger
	c      *config.Service
	ctx    context.Context
	queuer UpstreamQueuer
	stream net.Conn
}

func NewClient(ctx context.Context, logger zap.Logger, conf *config.Service, queuer UpstreamQueuer, stream net.Conn) *Client {
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
	buffer := make([]byte, c.c.MaxMessageLength, c.c.MaxMessageLength)
	upstreamDone := make(chan int, 1)

	for {
		err := c.stream.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err != nil {
			return err
		}

		n, err := c.stream.Read(buffer[:8])
		if err != nil {
			return err
		}

		if n < 8 {
			return ErrFrameRead
		}

		msgFrame, err := frame.Build([8]byte(buffer[:8]))
		if err != nil {
			return err
		}

		err = c.stream.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return err
		}

		n, err = c.stream.Read(buffer[:msgFrame.MessageLength])
		if err != nil {
			return err
		}

		// Is this conversion correct?
		// a u32 should fit in i64
		if n != int(msgFrame.MessageLength) {
			return ErrPayloadRead
		}

		c.queuer.Enqueue(buffer, msgFrame.MessageLength, upstreamDone)
		upstreamPayloadLen := <-upstreamDone
		if upstreamPayloadLen == 0 {
			return err
		}

		err = c.stream.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return err
		}

		n, err = c.stream.Write(buffer[0:upstreamPayloadLen])
		if err != nil {
			return err
		}

		if n != upstreamPayloadLen {
			return ErrWrite
		}
	}
}
