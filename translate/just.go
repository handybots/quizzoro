package translate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type JustTranslateService struct {
}

func (srv *JustTranslateService) Translate(from, to, text string) (string, error) {
	params := url.Values{}
	params.Set("text", text)
	params.Set("lang_from", from)
	params.Set("lang_to", to)

	req, err := http.NewRequest(
		http.MethodGet,
		"https://just-translated.p.rapidapi.com/?"+params.Encode(),
		nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-RapidApi-Host", "just-translated.p.rapidapi.com")
	req.Header.Add("X-RapidApi-Key", os.Getenv("RAPIDAPI_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("translate: just: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("translate: just: response code is %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Message string   `json:"message"`
		Text    []string `json:"text"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	if result.Message != "" {
		return "", fmt.Errorf("translate: just: %s", result.Message)
	}

	return result.Text[0], nil
}
