package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/atomic"
)

type YandexService struct {
	sid atomic.String
}

func (srv *YandexService) Translate(from, to, text string) (string, error) {
	form := url.Values{}
	form.Set("text", text)
	form.Set("option", "4")

	params := url.Values{}
	params.Set("id", srv.sid.Load())
	params.Set("lang", from+"-"+to)
	params.Set("srv", "tr-text")
	params.Set("reason", "auto")
	params.Set("format", "text")

	req, err := http.NewRequest(
		http.MethodPost,
		"https://translate.yandex.net/api/v1/tr.json/translate?"+params.Encode(),
		strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("translate: yandex: response code is %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Text    []string `json:"text"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	if result.Code != 200 {
		return "", fmt.Errorf("translate: yandex: %s (code=%d)", result.Message, result.Code)
	}

	return strings.Join(result.Text, ""), nil
}

var (
	ErrNoSID = errors.New("translate: no sid found")
	reSID    = regexp.MustCompile(`sid: *'([^']+)'`)
)

func (srv *YandexService) UpdateSID() error {
	req, err := http.NewRequest(http.MethodGet, "https://translate.yandex.ru", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("translate: yandex: response code is %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	match := reSID.FindStringSubmatch(string(data))
	if len(match) != 2 {
		return ErrNoSID
	}

	group := strings.Split(match[1], ".")
	if len(group) != 3 {
		return ErrNoSID
	}

	group = []string{
		reverseString(group[0]),
		reverseString(group[1]),
		reverseString(group[2]),
	}

	sid := strings.Join(group, ".") + "-0-0"
	srv.sid.Store(sid)
	return nil
}

func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
