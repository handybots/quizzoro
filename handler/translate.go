package handler

import (
	"log"
	"time"

	"github.com/demget/quizzorobot/translate"
)

func init() {
	go func() {
		for {
			err := translate.Yandex.UpdateSID()
			if err != nil {
				log.Println(err)
				return
			}

			time.Sleep(24 * time.Hour)
		}
	}()
}

func translateText(input string) (string, error) {
	var (
		output string
		err    error
	)

	output, err = translate.Google.Translate("en", "ru", input)
	if err == nil {
		return output, nil
	}
	log.Println(err)

	output, err = translate.DeepL.Translate("en", "ru", input)
	if err == nil {
		return output, nil
	}
	log.Println(err)

	output, err = translate.Yandex.Translate("en", "ru", input)
	if err != nil {
		return "", err
	}

	return output, nil
}
