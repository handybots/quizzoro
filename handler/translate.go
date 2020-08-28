package handler

import (
	"log"

	"github.com/demget/quizzorobot/translate"
)

// func init() {
// 	go func() {
// 		for {
// 			err := translate.Yandex.UpdateSID()
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
//
// 			time.Sleep(24 * time.Hour)
// 		}
// 	}()
// }

func translateText(input string) (output string, err error) {
	output, err = translate.MyMemory.Translate("en", "ru", input)
	if err == nil {
		return output, nil
	}
	log.Println(err)

	output, err = translate.Google.Translate("en", "ru", input)
	if err == nil {
		return output, nil
	}
	log.Println(err)

	output, err = translate.DeepL.Translate("en", "ru", input)
	if err != nil {
		return "", err
	}

	return output, nil
}
