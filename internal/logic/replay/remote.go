package replay

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func UnpackRemote(link string) (*UnpackedReplay, error) {
	resp, err := http.DefaultClient.Get(link)
	if err != nil {
		return nil, ErrInvalidReplayFile
	}
	defer resp.Body.Close()

	fmt.Printf("%s\n", link)
	fmt.Printf("%s\n", resp.Status)
	fmt.Printf("%d\n", resp.StatusCode)
	fmt.Printf("%s\n", resp.Header.Get("Content-Length"))

	// Convert 10 MB to bytes
	const maxFileSize = 10 * 1024 * 1024

	// Check the Content-Length header
	contentLengthStr := resp.Header.Get("Content-Length")
	if contentLengthStr == "" {
		log.Warn().Msg("Content-Length header is missing on remote replay file")
		return nil, ErrInvalidReplayFile
	}

	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil {
		return nil, err
	}
	if contentLength > maxFileSize {
		return nil, ErrInvalidReplayFile
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return Unpack(bytes.NewReader(data), contentLength)
}
