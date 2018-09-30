package game

import "github.com/bpiddubnyi/lottery"

type PairStack interface {
	Pop() (lottery.Pair, error)
}

// game describes lottery game logic
type game struct {
	jackpot uint64
	stack   PairStack
}

// newGame creates new Game instance
func newGame() (*game, error) {
	g := &game{}

	var err error
	g.stack, err = newWinStack()
	if err != nil {
		return nil, err
	}

	return g, nil
}

// play checks if player's bet metches to a win pair from lucky pairs stack and
// returns a match result
func (g *game) play(fee uint64, bet lottery.Pair) (*lottery.Response, error) {
	win, err := g.stack.Pop()
	if err != nil {
		return nil, err
	}

	r := &lottery.Response{Type: lottery.NoWin}
	if win == bet {
		if g.jackpot != 0 {
			r.Type = lottery.Win
		} else {
			r.Type = lottery.Bonus
		}
	}

	g.jackpot += fee
	r.Jackpot = g.jackpot
	return r, nil
}
