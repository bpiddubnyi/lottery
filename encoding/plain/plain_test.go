package plain

import (
	"bytes"
	"math"
	"testing"

	"github.com/bpiddubnyi/lottery"
	"github.com/google/uuid"
)

func TestRequestEncoder_Encode(t *testing.T) {
	type fields struct {
		w *bytes.Buffer
	}
	type args struct {
		r *lottery.Request
	}

	id, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		buf     []byte
	}{
		{
			name: "!#",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Request{
					UUID:  id,
					Fee:   42,
					Guess: lottery.Pair{33, 35}, // ASCII 33: !, 35: #
				},
			},
			wantErr: false,
			buf:     []byte("550e8400-e29b-41d4-a716-446655440000 42 !#"),
		},
		{
			name: "MaxUint64",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Request{
					UUID:  id,
					Fee:   math.MaxUint64,
					Guess: lottery.Pair{33, 35}, // ASCII 33: !, 35: #
				},
			},
			wantErr: false,
			buf:     []byte("550e8400-e29b-41d4-a716-446655440000 18446744073709551615 !#"),
		},
		{
			name: "zero",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Request{
					UUID:  id,
					Fee:   0,
					Guess: lottery.Pair{33, 35}, // ASCII 33: !, 35: #
				},
			},
			wantErr: false,
			buf:     []byte("550e8400-e29b-41d4-a716-446655440000 0 !#"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &RequestEncoder{
				w: tt.fields.w,
			}
			if err := enc.Encode(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("RequestEncoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && !bytes.Equal(tt.buf, tt.fields.w.Bytes()) {
				t.Errorf("RequestEncoder.Encode() {%s} != {%s}", string(tt.buf), string(tt.fields.w.Bytes()))
			}
		})
	}
}

func TestResponseEncoder_Encode(t *testing.T) {
	type fields struct {
		w *bytes.Buffer
	}
	type args struct {
		r *lottery.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		buf     []byte
	}{
		{
			name: "win",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Response{
					Type:    lottery.Win,
					Jackpot: 42,
				},
			},
			wantErr: false,
			buf:     []byte("win 42 "),
		},
		{
			name: "nowin",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Response{
					Type:    lottery.NoWin,
					Jackpot: 0,
				},
			},
			wantErr: false,
			buf:     []byte("nowin "),
		},
		{
			name: "bonus",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Response{
					Type:    lottery.Bonus,
					Jackpot: 0,
				},
			},
			wantErr: false,
			buf:     []byte("bonus "),
		},
		{
			name: "MaxUint64",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Response{
					Type:    lottery.Win,
					Jackpot: math.MaxUint64,
				},
			},
			wantErr: false,
			buf:     []byte("win 18446744073709551615 "),
		},
		{
			name: "wrong type",
			fields: fields{
				w: &bytes.Buffer{},
			},
			args: args{
				r: &lottery.Response{
					Type:    lottery.ResponseType(42),
					Jackpot: math.MaxUint64,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &ResponseEncoder{
				w: tt.fields.w,
			}
			if err := enc.Encode(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("ResponseEncoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && !bytes.Equal(tt.buf, tt.fields.w.Bytes()) {
				t.Errorf("RequestEncoder.Encode() {%s} != {%s}", string(tt.buf), string(tt.fields.w.Bytes()))
			}
		})
	}
}

func TestRequestDecoder_Decode(t *testing.T) {
	type fields struct {
		r *bytes.Reader
	}
	type args struct {
		r *lottery.Request
	}

	id, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		res     lottery.Request
	}{
		{
			name: "42",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-41d4-a716-446655440000 42 !#")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: false,
			res: lottery.Request{
				UUID:  id,
				Fee:   42,
				Guess: lottery.Pair{33, 35},
			},
		},
		{
			name: "MaxUint64",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-41d4-a716-446655440000 18446744073709551615 !#")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: false,
			res: lottery.Request{
				UUID:  id,
				Fee:   math.MaxUint64,
				Guess: lottery.Pair{33, 35},
			},
		},
		{
			name: "bad UUID",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-a716-446655440000 18446744073709551615 !#")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: true,
		},
		{
			name: "bad fee",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-41d4-a716-446655440000 bad !#")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: true,
		},
		{
			name: "short guess",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-41d4-a716-446655440000 18446744073709551615 !")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: true,
		},
		{
			name: "no guess",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-41d4-a716-446655440000 18446744073709551615 ")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: true,
		},
		{
			name: "short fee",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-41d4-a716-446655440000 18446744073")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: true,
		},
		{
			name: "short uuid",
			fields: fields{
				r: bytes.NewReader([]byte("550e8400-e29b-41d4-a716-446655440000")),
			},
			args: args{
				r: &lottery.Request{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &RequestDecoder{
				r: tt.fields.r,
			}
			if err := dec.Decode(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("RequestDecoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && tt.res != *tt.args.r {
				t.Errorf("RequestDecoder.Decode() {%v} != {%v}", tt.res, *tt.args.r)
			}
		})
	}
}

func TestResponseDecoder_Decode(t *testing.T) {
	type fields struct {
		r *bytes.Reader
	}
	type args struct {
		r *lottery.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		res     lottery.Response
	}{
		{
			name: "win",
			fields: fields{
				r: bytes.NewReader([]byte("win 42 ")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: false,
			res: lottery.Response{
				Type:    lottery.Win,
				Jackpot: 42,
			},
		},
		{
			name: "nowin",
			fields: fields{
				r: bytes.NewReader([]byte("nowin ")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: false,
			res: lottery.Response{
				Type:    lottery.NoWin,
				Jackpot: 0,
			},
		},
		{
			name: "bonus",
			fields: fields{
				r: bytes.NewReader([]byte("bonus ")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: false,
			res: lottery.Response{
				Type:    lottery.Bonus,
				Jackpot: 0,
			},
		},
		{
			name: "win_short_jackpot",
			fields: fields{
				r: bytes.NewReader([]byte("win 42")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: true,
		},
		{
			name: "win_no_jackpot",
			fields: fields{
				r: bytes.NewReader([]byte("win ")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: true,
		},
		{
			name: "short type",
			fields: fields{
				r: bytes.NewReader([]byte("w ")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: true,
		},
		{
			name: "wrong type",
			fields: fields{
				r: bytes.NewReader([]byte("winwin 43434")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: true,
		},
		{
			name: "short message",
			fields: fields{
				r: bytes.NewReader([]byte("win")),
			},
			args: args{
				r: &lottery.Response{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &ResponseDecoder{
				r: tt.fields.r,
			}
			if err := dec.Decode(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("ResponseDecoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && tt.res != *tt.args.r {
				t.Errorf("RequestDecoder.Decode() {%v} != {%v}", tt.res, *tt.args.r)
			}
		})
	}
}
