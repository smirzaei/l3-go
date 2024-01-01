package downstream

import (
	"context"
	"fmt"
	"net"

	"github.com/smirzaei/13-go/internal/config"
	"go.uber.org/zap"
)

type UpstreamQueuer interface {
	Enqueue(buffer []byte, msgLen uint32, done chan<- int)
}

type Server struct {
	l      zap.Logger
	c      *config.Service
	queuer UpstreamQueuer
}

func NewServer(logger zap.Logger, conf *config.Service, queuer UpstreamQueuer) *Server {
	s := Server{
		l:      logger,
		c:      conf,
		queuer: queuer,
	}

	return &s
}

func (s *Server) Listen(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.c.Port))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.l.Warn("failed to accept downstream connection", zap.Error(err))
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// TODO: Handle the graceful shutdown of the clients using the context below
	client := NewClient(context.TODO(), s.l, s.c, s.queuer, conn)
	err := client.Serve()
	if err != nil {
		s.l.Warn("downstream connection error", zap.Error(err))
	}
}
