package encoding

import (
	"io"

	"github.com/bpiddubnyi/lottery"
)

// RequestEncoder is a common interface for
// different request encoding implementations
type RequestEncoder interface {
	Encode(*lottery.Request) error
}

// RequestDecoder is a common interface for
// different request decoding implementations
type RequestDecoder interface {
	Decode(*lottery.Request) error
}

// ResponseEncoder is a common interface for
// different response encoding implementations
type ResponseEncoder interface {
	Encode(*lottery.Response) error
}

// ResponseDecoder is a common interface for
// different response decoding implementations
type ResponseDecoder interface {
	Decode(*lottery.Response) error
}

// Server is an interface for getting server-side
// encoder-decoder pair
type Server interface {
	GetRequestDecoder(r io.Reader) RequestDecoder
	GetResponseEncoder(w io.Writer) ResponseEncoder
}

// Client is an interface for getting client-side
// encoder-decoder pair
type Client interface {
	GetRequestEncoder(w io.Writer) RequestEncoder
	GetResponseDecoder(r io.Reader) ResponseDecoder
}
