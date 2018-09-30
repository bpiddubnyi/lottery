package game

import "github.com/bpiddubnyi/lottery"

// PairStack is a comon interface for different lucky pair stack implementations
type PairStack interface {
	Pop() (lottery.Pair, error)
}

// Game describes lottery game logic
type Game struct {
	Jackpot uint64
	Stack   PairStack
}

// New creates new Game instance
func New(stack PairStack) *Game {
	return &Game{Stack: stack}
}

// Play checks if player's bet metches to a win pair from lucky pairs stack and
// returns a match result
func (g *Game) Play(fee uint64, bet lottery.Pair) (*lottery.Response, error) {
	win, err := g.Stack.Pop()
	if err != nil {
		return nil, err
	}

	r := &lottery.Response{Type: lottery.NoWin}
	if win == bet {
		if g.Jackpot != 0 {
			r.Type = lottery.Win
			r.Jackpot = g.Jackpot + fee
			g.Jackpot = 0
		} else {
			r.Type = lottery.Bonus
			g.Jackpot = fee
		}
	} else {
		g.Jackpot += fee
	}
	return r, nil
}
