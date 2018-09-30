package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/bpiddubnyi/lottery"
	"github.com/bpiddubnyi/lottery/cmd/lotteryd/game"
	"github.com/bpiddubnyi/lottery/encoding"
	"github.com/bpiddubnyi/lottery/encoding/plain"
)

const (
	defaultTimeout = 10 * time.Second
	defaultWorkers = 10
)

var (
	defaultProtocol encoding.Server = plain.Server{}
)

type Server struct {
	// Protocol implementation
	Proto encoding.Server
	// Connection timeout
	Timeout time.Duration
	// Number of worker routines
	Workers uint

	game  *game.Game
	gameL sync.Mutex
}

func New(stack game.PairStack) *Server {
	return &Server{
		Timeout: defaultTimeout,
		Workers: defaultWorkers,
		Proto:   defaultProtocol,
		game:    game.New(stack),
	}
}

func (s *Server) Listen(ctx context.Context, addr string) error {
	var wg sync.WaitGroup

	lCtx, lCancel := context.WithCancel(ctx)
	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		lCancel()
		return err
	}

	wg.Add(1)
	go func() {
		<-lCtx.Done()
		l.Close()
		wg.Done()
	}()

	connC := make(chan net.Conn)
	for i := uint(0); i < s.Workers; i++ {
		wg.Add(1)
		go func() {
			s.work(connC)
			wg.Done()
		}()
	}

	var c net.Conn
theLoop:
	for {
		c, err = l.Accept()
		if err != nil {
			break
		}

		select {
		case connC <- c:
		case <-ctx.Done():
			break theLoop
		}
	}

	close(connC)
	wg.Wait()
	lCancel()
	return err
}

func (s *Server) play(fee uint64, bet lottery.Pair) (*lottery.Response, error) {
	s.gameL.Lock()
	defer s.gameL.Unlock()

	return s.game.Play(fee, bet)
}

func (s *Server) match(c net.Conn) (*lottery.Response, error) {
	dec := s.Proto.GetRequestDecoder(c)
	enc := s.Proto.GetResponseEncoder(c)
	req := lottery.Request{}

	remote := c.RemoteAddr().String()
	if err := dec.Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode initial request: %s", err)
	}
	log.Printf("info: %s: request: %s", remote, req.String())

	resp, err := s.play(req.Fee, req.Guess)
	if err != nil {
		return nil, fmt.Errorf("game failed: %s", err)
	}
	log.Printf("info: %s: response: %s", remote, resp.String())

	if err := enc.Encode(resp); err != nil {
		return nil, fmt.Errorf("failed to send response: %s", err)
	}

	return resp, nil
}

func (s *Server) handleConn(c net.Conn) error {
	defer c.Close()

	resp, err := s.match(c)
	if err != nil {
		return err
	}

	if resp.Type != lottery.Bonus {
		return nil
	}

	_, err = s.match(c)
	return err
}

func (s *Server) work(connC <-chan net.Conn) {
	for c := range connC {
		c.SetDeadline(time.Now().Add(s.Timeout))
		remote := c.RemoteAddr().String()
		log.Printf("info: %s: new connection", remote)
		if err := s.handleConn(c); err != nil {
			log.Printf("error: %s: failed to handle connection: %s", remote, err)
		}
		log.Printf("info: current jackpot: %d", s.game.Jackpot)
	}
}
