package game

import (
	"container/list"
	"crypto/rand"
	"errors"

	"github.com/bpiddubnyi/lottery"
)

type winStack struct {
	l *list.List
}

func newWinStack() (*winStack, error) {
	l := list.New()
	for i := 0; i < stackLen; i++ {
		var p lottery.Pair

		_, err := rand.Read(p[:])
		if err != nil {
			return nil, err
		}

		l.PushBack(p)
	}
	return &winStack{l: l}, nil
}

func (s *winStack) Pop() (lottery.Pair, error) {
	var p lottery.Pair

	// Sanity checks. None of this should ever really happen
	e := s.l.Front()
	if e == nil {
		return p, errors.New("win stack is empty")
	}

	p, ok := e.Value.(lottery.Pair)
	if !ok {
		return p, errors.New("wrong value type in win stack, *LuckPair is expected")
	}

	var newP lottery.Pair
	_, err := rand.Read(newP[:])
	if err != nil {
		return p, err
	}

	s.l.Remove(e)
	s.l.PushBack(p)

	return p, nil
}
