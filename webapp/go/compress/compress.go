package compress

import (
	"strings"

	"github.com/google/brotli/go/cbrotli"
)

const (
	// Threshold is the minimum size to be compressed
	Threshold = 1024 // byte
	Brotli    = "br"
	Gzip      = "gzip"
)

func parseAcceptEncoding(in string) []string {
	if in == "" {
		return nil
	}
	es := strings.Split(in, ",")
	ret := make([]string, 0, len(es))
	for _, e := range es {
		if e == "" {
			continue
		}
		v := strings.TrimSpace(strings.Split(e, ";")[0]) // ignore quality
		if v == "" || v == "*" {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

// EncoderFactory creates encoder and returns it.
func EncoderFactory(size int, acceptEncodingHeader string) Encoder {
	availableEncodes := parseAcceptEncoding(acceptEncodingHeader)
	if size < Threshold || len(availableEncodes) == 0 {
		return &noEncoder{}
	}
	for _, e := range availableEncodes {
		if e == Brotli {
			return &brEncoder{}
		}
	}
	return &noEncoder{}
}

type Encoder interface {
	Encode(data []byte) (out []byte, encode string, err error)
}

type noEncoder struct{}

func (e *noEncoder) Encode(data []byte) ([]byte, string, error) {
	return data, "", nil
}

type brEncoder struct{}

func (e *brEncoder) Encode(data []byte) ([]byte, string, error) {
	out, err := cbrotli.Encode(data, cbrotli.WriterOptions{Quality: 6, LGWin: 22})
	if err != nil {
		return nil, Brotli, err
	}
	return out, Brotli, nil
}
