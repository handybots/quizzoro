package handler

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/demget/quizzorobot/yandextr"
	"go.uber.org/atomic"
)

var currentSID atomic.String

func init() {
	go func() {
		for {
			sid, err := yandextr.ParseSID()
			if err != nil {
				log.Println(err)
				return
			}
			currentSID.Store(sid)
			time.Sleep(24 * time.Hour)
		}
	}()
}

func translateText(text string) (string, error) {
	sid := currentSID.Load()
	if sid == "" {
		return "", errors.New("sid is empty")
	}
	result, err := yandextr.Translate(currentSID.Load(), text)
	if err != nil {
		return "", err
	}
	return strings.Join(result.Text, ""), nil
}
