package downstream

import (
	"context"

	"github.com/smirzaei/13-go/internal/config"
	"go.uber.org/zap"
)

type UpstreamQueuer interface {
	Enqueue(buffer []byte, msgLen uint, done chan<- error)
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
	panic("not implemented")
}
