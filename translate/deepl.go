package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DeepLService struct {
	Key string
	url string
}

func (srv *DeepLService) Translate(from, to, text string) (string, error) {
	params := url.Values{}
	params.Set("auth_key", srv.Key)
	params.Set("source_lang", from)
	params.Set("target_lang", to)
	params.Set("text", text)

	req, err := http.NewRequest(http.MethodPost,
		srv.url+"?auth_key="+srv.Key,
		strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("translate: deepl: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("translate: deepl: response code is %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	if len(result.Translations) == 0 {
		return "", errors.New("translate: deepl: no translations")
	}

	return result.Translations[0].Text, nil
}
