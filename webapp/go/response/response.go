package response

import (
	"log"
	"net/http"
	"strconv"

	"github.com/isucon/isucon9-qualify/webapp/go/compress"
)

func SendBody(w http.ResponseWriter, r *http.Request, data []byte, status int) error {
	// compress
	var encode string
	aes := r.Header.Get("accept-encoding")
	data, encode, err := compress.EncoderFactory(len(data), aes).Encode(data)
	if err != nil {
		log.Printf("response: failed to compress response data\n")
		return err
	}
	if encode != "" {
		w.Header().Set("content-encoding", encode)
	}

	// write
	w.Header().Set("content-length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	_, err := w.Write(data)
	return err
}
