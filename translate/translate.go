package translate

import (
	"net/url"
	"os"

	"go.uber.org/atomic"
)

type Translator interface {
	Translate(from, to, text string) (string, error)
}

var (
	Google = GoogleService{}
	Yandex = YandexService{
		atomic.NewString(""),

		"https://translate.yandex.ru",
		"https://translate.yandex.net/api",
		"/v1/tr.json/translate?",
		url.Values{
			"srv":    {"tr-text"},
			"lang":   {"en-ru"},
			"reason": {"auto"},
			"format": {"text"},
		},
	}
	DeepL = DeepLService{
		API:          os.Getenv("DEEPL_AUTHKEY"),
		translateURL: "https://api.deepl.com/v2/translate",
	}
)
