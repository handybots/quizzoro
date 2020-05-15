package translate

import (
	gt "github.com/bas24/googletranslatefree"
)

type GoogleService struct{}

func (srv *GoogleService) Translate(from, to, text string) (string, error) {
	return gt.Translate(text, from, to)
}
