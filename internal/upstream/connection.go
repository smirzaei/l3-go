package upstream

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/smirzaei/13-go/internal/frame"
	"go.uber.org/zap"
)

var (
	ErrFrameRead = fmt.Errorf("failed to read the frame")
	ErrWrite     = fmt.Errorf("failed to write the response")
)

type Connection struct {
	l      *zap.Logger
	stream net.Conn
}

func NewConnection(logger *zap.Logger, host string) (*Connection, error) {
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
			err := c.stream.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return err
			}

			n, err := c.stream.Write(req.buffer[:req.msgLen])
			if err != nil {
				return err
			}

			if n != int(req.msgLen) {
				return ErrWrite
			}

			err = c.stream.SetReadDeadline(time.Now().Add(30 * time.Second))
			if err != nil {
				return err
			}

			n, err = c.stream.Read(req.buffer[:8])
			if err != nil {
				return err
			}

			if n != 8 {
				return ErrFrameRead
			}

			msgFrame, err := frame.Build([8]byte(req.buffer[:8]))
			if err != nil {
				return err
			}

			n, err = io.ReadFull(c.stream, req.buffer[:msgFrame.MessageLength])
			if err != nil {
				return err
			}

			req.done <- n
		}
	}
}
