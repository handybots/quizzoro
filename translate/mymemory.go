package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MyMemoryService struct {
}

func (srv *MyMemoryService) Translate(from, to, text string) (string, error) {
	params := url.Values{}
	params.Set("q", text)
	params.Set("langpair", from+"|"+to)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.mymemory.translated.net/get?"+params.Encode(),
		nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("translate: mymemory: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("translate: mymemory: response code is %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Text string `json:"translatedText"`
		} `json:"responseData"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	if result.Data.Text == "" {
		return "", errors.New("translate: mymemory: no translations")
	}

	return result.Data.Text, nil
}
