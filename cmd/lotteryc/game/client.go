package game

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"

	"github.com/bpiddubnyi/lottery"
	"github.com/bpiddubnyi/lottery/encoding"
	"github.com/bpiddubnyi/lottery/encoding/plain"
)

var (
	defaultProto encoding.Client = plain.Client{}
)

type Client struct {
	Proto encoding.Client
	addr  string
}

func NewClient(addr string) *Client {
	return &Client{Proto: defaultProto, addr: addr}
}

func (cli *Client) Play(fee uint64) (*lottery.Response, error) {
	c, err := net.Dial("tcp", cli.addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %s", err)
	}
	defer c.Close()

	req, err := genInitRequest(fee)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial request: %s", err)
	}

	addr := c.LocalAddr().String()
	log.Printf("info: %s request: %s", addr, req.String())

	enc := cli.Proto.GetRequestEncoder(c)
	if err = enc.Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode request: %s", err)
	}

	resp := &lottery.Response{}
	dec := cli.Proto.GetResponseDecoder(c)
	if err = dec.Decode(resp); err != nil {
		return nil, fmt.Errorf("failed to decode initial response: %s", err)
	}

	log.Printf("info: %s response: %s", addr, resp.String())

	if resp.Type != lottery.Bonus {
		return resp, nil
	}

	if err = genBonusRequest(req); err != nil {
		return nil, fmt.Errorf("failed to create bonus request: %s", err)
	}
	if err = enc.Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode bonus request: %s", err)
	}
	if err = dec.Decode(resp); err != nil {
		return nil, fmt.Errorf("failed to decode bonus response: %s", err)
	}

	log.Printf("info: %s response: %s", addr, resp.String())

	return resp, nil
}

func genInitRequest(fee uint64) (*lottery.Request, error) {
	var err error
	req := &lottery.Request{Fee: fee}

	req.UUID, err = uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	_, err = rand.Read(req.Guess[:])
	if err != nil {
		return nil, err
	}

	return req, nil
}

func genBonusRequest(req *lottery.Request) error {
	req.Fee = 0

	_, err := rand.Read(req.Guess[:])
	return err
}
