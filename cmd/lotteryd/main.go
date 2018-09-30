package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bpiddubnyi/lottery/cmd/lotteryd/game"
)

var (
	timeout  = 5
	workers  = uint(runtime.NumCPU())
	showHelp bool
	addr     = ":9876"
)

func init() {
	flag.UintVar(&workers, "w", workers, "number of workers")
	flag.BoolVar(&showHelp, "h", false, "show this help and exit")
	flag.IntVar(&timeout, "t", timeout, "connection timeout in seconds")
	flag.StringVar(&addr, "a", addr, "listen address")
}

func main() {
	flag.Parse()
	if showHelp {
		flag.Usage()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigC := make(chan os.Signal, 1)
	defer close(sigC)

	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigC
		log.Printf("info: signal received: %s", s)
		cancel()
	}()

	s, err := game.NewServer()
	if err != nil {
		log.Printf("error: failed to create game server: %s", err)
	}

	s.Timeout = time.Duration(timeout) * time.Second
	s.Workers = workers

	if err := s.Listen(ctx, addr); err != nil {
		log.Printf("error: server failed: %s", err)
	}
}
