package lottery

import (
	"fmt"

	"github.com/google/uuid"
)

// Pair represents players guess for a lottery winning combination
type Pair [2]byte

func (p Pair) String() string {
	return fmt.Sprintf("%d:%d", p[0], p[1])
}

// Request is a client request message to the lottery game server
type Request struct {
	UUID  uuid.UUID
	Fee   uint64
	Guess Pair
}

func (r Request) String() string {
	return fmt.Sprintf("uuid: %s: fee: %d guess: %s",
		r.UUID.String(), r.Fee, r.Guess.String())
}

// ResponseType is a servers response message enum type
type ResponseType int

// Response types
const (
	NoWin ResponseType = iota
	Win
	Bonus
)

func (t ResponseType) String() string {
	switch t {
	case NoWin:
		return "no win"
	case Win:
		return "win"
	case Bonus:
		return "bonus"
	default:
		return "unknown"
	}
}

// Response is a server response message to the client
type Response struct {
	Type    ResponseType
	Jackpot uint64
}

func (r Response) String() string {
	switch r.Type {
	case Win:
		return fmt.Sprintf("%s: %d", r.Type, r.Jackpot)
	default:
		return r.Type.String()
	}
}
