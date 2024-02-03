package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cufee/aftermath-core/internal/core/utils"
)

var replayUploadUrl string

func init() {
	replayUploadUrl = utils.MustGetEnv("WOT_INSPECTOR_REPLAYS_URL")
}

func GetReplayData(replayUrl string) (*Replay, error) {
	link := url.Values{}
	link.Set("upload_url", replayUrl)
	link.Set("title", "Aftermath Replay Upload")
	// link.Set("private", "true") // Does not work -- tried true and 1

	req, err := http.NewRequest("POST", replayUploadUrl, strings.NewReader(link.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := insecureClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		var data map[string][]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to get replay data: %s", res.Status)
		}
		for _, v := range data {
			if len(v) > 0 {
				return nil, errors.New(strings.ToLower(v[0]))
			}
		}
		return nil, errors.New("failed to get replay data: " + res.Status)
	}

	var data replayData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Private: %t\n", data.Private)

	return data.Replay(), nil
}
