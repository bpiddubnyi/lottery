package game

import (
	"crypto/rand"

	"github.com/bpiddubnyi/lottery"
)

const (
	stackLen  = 100
	stackSize = stackLen * 2
)

type winRing struct {
	data [stackSize]byte
	cur  int
}

func (s *winRing) Pop() (lottery.Pair, error) {
	if s.cur+2 > stackSize {
		s.cur = 0
	}

	var win lottery.Pair
	copy(win[:], s.data[s.cur:s.cur+2])

	_, err := rand.Read(s.data[s.cur : s.cur+2])
	if err == nil {
		s.cur += 2
	}

	return win, err
}

func newWinRing() (*winRing, error) {
	s := &winRing{}
	_, err := rand.Read(s.data[:])
	if err != nil {
		return nil, err
	}

	return s, nil
}
