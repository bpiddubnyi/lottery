package plain

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/bpiddubnyi/lottery"
	"github.com/bpiddubnyi/lottery/encoding"
	"go4.org/strutil"
)

const (
	fieldSeparator byte = ' '
)

type RequestEncoder struct {
	w io.Writer
}

func NewRequestEncoder(w io.Writer) *RequestEncoder {
	return &RequestEncoder{w: w}
}

func (enc *RequestEncoder) Encode(r *lottery.Request) error {
	data, _ := r.UUID.MarshalText()
	data = append(data, fieldSeparator)
	data = strconv.AppendUint(data, r.Fee, 10)
	data = append(data, fieldSeparator)
	data = append(data, r.Guess[:]...)

	_, err := enc.w.Write(data)
	return err
}

type RequestDecoder struct {
	r io.Reader
}

func NewRequestDecoder(r io.Reader) *RequestDecoder {
	return &RequestDecoder{r: r}
}

func (dec *RequestDecoder) Decode(r *lottery.Request) error {
	// len(UUID): 36
	// len(MaxUInt64): 20
	// len(Pair): 2
	// => max token len = 36 + 1 (separator)
	buf := make([]byte, 37)

	// Read and parse UUID
	_, err := io.ReadFull(dec.r, buf)
	if err != nil {
		return err
	}
	err = r.UUID.UnmarshalText(buf[:36])
	if err != nil {
		return err
	}

	// Read and parse fee
	n, err := readUntil(dec.r, buf, fieldSeparator)
	if err != nil {
		return err
	}
	fee, err := strutil.ParseUintBytes(buf[:n], 10, 64)
	if err != nil {
		return err
	}
	r.Fee = fee

	// Read lucky pair
	_, err = io.ReadFull(dec.r, r.Guess[:])
	if err != nil {
		return err
	}
	return nil
}

type ResponseEncoder struct {
	w io.Writer
}

func NewResponseEncoder(w io.Writer) *ResponseEncoder {
	return &ResponseEncoder{w: w}
}

var (
	noWinB = []byte("nowin")
	winB   = []byte("win")
	bonusB = []byte("bonus")
)

func marshalResponseType(t lottery.ResponseType) ([]byte, error) {
	var (
		res = make([]byte, 5)
		c   []byte
	)

	switch t {
	case lottery.NoWin:
		c = noWinB
	case lottery.Win:
		c = winB
	case lottery.Bonus:
		c = bonusB
	default:
		return nil, fmt.Errorf("invalid value: '%d'", t)
	}

	copy(res, c)
	res = res[:len(c)]

	return res, nil
}

func (enc *ResponseEncoder) Encode(r *lottery.Response) error {
	data, err := marshalResponseType(r.Type)
	if err != nil {
		return err
	}

	data = append(data, fieldSeparator)
	if r.Type == lottery.Win {
		data = strconv.AppendUint(data, r.Jackpot, 10)
		data = append(data, fieldSeparator)
	}

	_, err = enc.w.Write(data)
	return err
}

type ResponseDecoder struct {
	r io.Reader
}

func NewResponseDecoder(r io.Reader) *ResponseDecoder {
	return &ResponseDecoder{r: r}
}

func unmarshalResponseType(data []byte) (lottery.ResponseType, error) {
	if len(data) < 3 {
		return lottery.NoWin, errors.New("string too short")
	}

	if bytes.Equal(data, noWinB) {
		return lottery.NoWin, nil
	} else if bytes.Equal(data, winB) {
		return lottery.Win, nil
	} else if bytes.Equal(data, bonusB) {
		return lottery.Bonus, nil
	} else {
		return lottery.NoWin, fmt.Errorf("invalid string: '%s'", data)
	}
}

func (dec *ResponseDecoder) Decode(r *lottery.Response) error {
	buf := make([]byte, 20)
	n, err := readUntil(dec.r, buf, fieldSeparator)
	if err != nil {
		return err
	}

	r.Type, err = unmarshalResponseType(buf[:n])
	if err != nil {
		return fmt.Errorf("failed to parse response type: %s", err)
	}
	if r.Type != lottery.Win {
		return nil
	}

	n, err = readUntil(dec.r, buf, fieldSeparator)
	if err != nil {
		return err
	}

	r.Jackpot, err = strutil.ParseUintBytes(buf[:n], 10, 64)
	return err
}

var (
	errNoDelimFound = errors.New("no delimiter found")
)

func readUntil(r io.Reader, buf []byte, delim byte) (int, error) {
	for i := range buf {
		_, err := r.Read(buf[i : i+1])

		if buf[i] == delim {
			return i, nil
		}
		if err != nil {
			return i, err
		}
	}
	return 0, errNoDelimFound
}

type Server struct{}

func (Server) GetRequestDecoder(r io.Reader) encoding.RequestDecoder {
	return NewRequestDecoder(r)
}

func (Server) GetResponseEncoder(w io.Writer) encoding.ResponseEncoder {
	return NewResponseEncoder(w)
}

type Client struct{}

func (Client) GetRequestEncoder(w io.Writer) encoding.RequestEncoder {
	return NewRequestEncoder(w)
}

func (Client) GetResponseDecoder(r io.Reader) encoding.ResponseDecoder {
	return NewResponseDecoder(r)
}
