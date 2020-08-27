package translate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type GoogleService struct{}

func (srv *GoogleService) Translate(from, to, text string) (string, error) {
	q := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=" +
		from + "&tl=" + to + "&dt=t&q=" + url.QueryEscape(text)

	resp, err := http.Get(q)
	if err != nil {
		return "", fmt.Errorf("translate: google: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("translate: google: %v", err)
	}

	bad := strings.Contains(string(body), `Error 400 (Bad Request)`)
	if bad {
		return "", fmt.Errorf("translate: google: bad request")
	}

	var (
		raw    []interface{}
		result []string
	)
	if err := json.Unmarshal(body, &raw); err != nil {
		log.Println("translate: google:", string(body))
		return "", fmt.Errorf("translate: google: %v", err)
	}

	if len(raw) > 0 {
		inner := raw[0]
		for _, slice := range inner.([]interface{}) {
			for _, t := range slice.([]interface{}) {
				result = append(result, fmt.Sprintf("%v", t))
				break
			}
		}
		return strings.Join(result, ""), nil
	}

	return "", fmt.Errorf("translate: google: no result")
}
