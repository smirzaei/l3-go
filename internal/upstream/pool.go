package upstream

import (
	"github.com/smirzaei/13-go/internal/config"
	"go.uber.org/zap"
)

type Request struct {
	buffer []byte
	msgLen uint32
	done   chan<- error
}

type Pool struct {
	l     zap.Logger
	c     *config.Config
	queue chan Request
}

func NewPool(logger zap.Logger, conf *config.Config) *Pool {
	p := Pool{
		l:     logger,
		c:     conf,
		queue: make(chan Request, 100),
	}

	return &p
}

func (p *Pool) Enqueue(buffer []byte, msgLen uint32, done chan<- error) {
	panic("not implemented")
}
