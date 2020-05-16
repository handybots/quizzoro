package translate

import (
	"os"

	"go.uber.org/atomic"
)

type Translator interface {
	Translate(from, to, text string) (string, error)
}

var (
	Google = GoogleService{}
	Yandex = YandexService{
		sid:    atomic.NewString(""),
		urlSID: "https://translate.yandex.ru",
		urlAPI: "https://translate.yandex.net/api/v1/tr.json/translate",
	}
	DeepL = DeepLService{
		Key: os.Getenv("DEEPL_AUTHKEY"),
		url: "https://api.deepl.com/v2/translate",
	}
)
