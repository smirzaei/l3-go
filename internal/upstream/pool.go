package upstream

import (
	"context"
	"time"

	"github.com/smirzaei/13-go/internal/config"
	"go.uber.org/zap"
)

type Request struct {
	buffer []byte
	msgLen uint32
	done   chan<- int
}

type Pool struct {
	l     *zap.Logger
	c     *config.Config
	queue chan Request
}

func NewPool(logger *zap.Logger, conf *config.Config) *Pool {
	p := Pool{
		l:     logger,
		c:     conf,
		queue: make(chan Request, 100),
	}

	return &p
}

func (p *Pool) Start(ctx context.Context) {
	for _, host := range p.c.Upstream.Hosts {
		for i := uint(0); i < p.c.Upstream.Connections; i++ {
			go p.handleConnection(host, time.Duration(i)*time.Millisecond*10)
		}
	}
}

func (p *Pool) handleConnection(host string, waitFor time.Duration) {
	time.Sleep(waitFor)
	c, err := NewConnection(p.l, host)
	if err != nil {
		// TODO: This needs to be improved. Right now it's only retrying
		//  but it's a recipe for disaster
		//  Need to implement a better connection pool by the amount of work
		//  pilled up in the queue

		p.handleConnection(host, waitFor+time.Second)
		return
	}

	waitFor = 0
	// TODO: use this context for graceful shutdown
	err = c.Serve(context.TODO(), p.queue)
	if err != nil {
		// TODO: Same as above. Retrying like this is not a good idea
		p.handleConnection(host, waitFor+time.Second)
	}
}

func (p *Pool) Enqueue(buffer []byte, msgLen uint32, done chan<- int) {
	req := Request{
		buffer: buffer,
		msgLen: msgLen,
		done:   done,
	}

	// TODO: This is a good place to check the queue and increase/decrease the
	//	connection pool as necessary
	p.queue <- req
}
