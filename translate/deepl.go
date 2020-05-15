package translate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DeepLService struct {
	API          string
	translateURL string
}

func (srv *DeepLService) Translate(from, to, text string) (string, error) {
	var data url.Values
	data.Set("auth_key", srv.API)
	data.Set("source_lang", from)
	data.Set("target_lang", to)
	data.Set("text", text)

	req, _ := http.NewRequest("POST",
		srv.translateURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("translate: deepl: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("translate: response code is %d", resp.StatusCode)
	}

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(jsonResp, &result); err != nil {
		return "", err
	}
	return result.Text, nil
}
