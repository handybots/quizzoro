package translate

import (
	"fmt"

	gt "github.com/bas24/googletranslatefree"
)

type GoogleService struct{}

func (srv *GoogleService) Translate(from, to, text string) (string, error) {
	result, err := gt.Translate(text, from, to)
	if err != nil {
		return "", fmt.Errorf("translate: google: %v", err)
	}
	return result, nil
}
