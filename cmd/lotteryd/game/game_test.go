package game

import (
	"reflect"
	"testing"

	"github.com/bpiddubnyi/lottery"
)

type stackMockOnes struct{}

func (stackMockOnes) Pop() (lottery.Pair, error) {
	return lottery.Pair{1, 1}, nil
}

func TestGame_Play(t *testing.T) {
	type fields struct {
		Jackpot uint64
		Stack   PairStack
	}
	type args struct {
		fee uint64
		bet lottery.Pair
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		want             *lottery.Response
		wantJackPotAfter uint64
		wantErr          bool
	}{
		{
			name: "win",
			fields: fields{
				Jackpot: 58,
				Stack:   stackMockOnes{},
			},
			args: args{
				fee: 42,
				bet: lottery.Pair{1, 1},
			},
			want: &lottery.Response{
				Type:    lottery.Win,
				Jackpot: 100,
			},
			wantErr:          false,
			wantJackPotAfter: 0,
		},
		{
			name: "nowin",
			fields: fields{
				Jackpot: 58,
				Stack:   stackMockOnes{},
			},
			args: args{
				fee: 42,
				bet: lottery.Pair{1, 2},
			},
			want: &lottery.Response{
				Type:    lottery.NoWin,
				Jackpot: 0,
			},
			wantErr:          false,
			wantJackPotAfter: 100,
		},
		{
			name: "bonus",
			fields: fields{
				Jackpot: 0,
				Stack:   stackMockOnes{},
			},
			args: args{
				fee: 42,
				bet: lottery.Pair{1, 1},
			},
			want: &lottery.Response{
				Type:    lottery.Bonus,
				Jackpot: 0,
			},
			wantErr:          false,
			wantJackPotAfter: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Jackpot: tt.fields.Jackpot,
				Stack:   tt.fields.Stack,
			}
			got, err := g.Play(tt.args.fee, tt.args.bet)
			if (err != nil) != tt.wantErr {
				t.Errorf("Game.Play() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Game.Play() = %v, want %v", got, tt.want)
			}
			if g.Jackpot != tt.wantJackPotAfter {
				t.Errorf("Game.Play() jackpot = %d, wantJackpot %d", g.Jackpot, tt.wantJackPotAfter)
			}
		})
	}
}
