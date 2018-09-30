package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bpiddubnyi/lottery/cmd/lotteryd/game"
	"github.com/bpiddubnyi/lottery/cmd/lotteryd/server"
)

var (
	timeout   = 5
	workers   = uint(runtime.NumCPU())
	showHelp  bool
	addr      = ":9876"
	container = "stack"
)

func init() {
	flag.UintVar(&workers, "w", workers, "number of workers")
	flag.BoolVar(&showHelp, "h", false, "show this help and exit")
	flag.IntVar(&timeout, "t", timeout, "connection timeout in seconds")
	flag.StringVar(&addr, "a", addr, "listen address")
	flag.StringVar(&container, "c", container, "lucky pair container type (stack, ring)")
}

func main() {
	flag.Parse()
	if showHelp {
		flag.Usage()
		return
	}

	con, err := getPairContainer(container)
	if err != nil {
		fmt.Printf("failed to initialize game lucky pair container: %s\n", err)
		flag.Usage()
		os.Exit(1)
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

	s := server.New(con)

	s.Timeout = time.Duration(timeout) * time.Second
	s.Workers = workers

	if err := s.Listen(ctx, addr); err != nil {
		log.Printf("error: server failed: %s", err)
	}
}

func getPairContainer(s string) (game.PairStack, error) {
	switch strings.ToLower(s) {
	case "stack":
		return game.NewWinStack()
	case "ring":
		return game.NewWinRing()
	default:
		return nil, fmt.Errorf("invalid value \"%s\"", s)
	}
}
