package main

import (
	"flag"
	"log"

	"github.com/bpiddubnyi/lottery"
	"github.com/bpiddubnyi/lottery/cmd/lotteryc/game"
)

var (
	addr     = "127.0.0.1:9876"
	showHelp bool
	fee      uint64 = 150
)

func init() {
	flag.StringVar(&addr, "a", addr, "server address")
	flag.BoolVar(&showHelp, "h", false, "show this help and exit")
	flag.Uint64Var(&fee, "f", fee, "fee value")
}

func main() {
	flag.Parse()
	if showHelp {
		flag.Usage()
		return
	}

	c := game.NewClient(addr)
	resp, err := c.Play(fee)
	if err != nil {
		log.Fatalf("fatal: play failed: %s", err)
	}
	if resp.Type == lottery.Win {
		log.Printf("You won %d!", resp.Jackpot)
	}
}
