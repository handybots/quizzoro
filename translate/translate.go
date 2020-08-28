package translate

import (
	"os"

	"go.uber.org/atomic"
)

type Translator interface {
	Translate(from, to, text string) (string, error)
}

var (
	Google   = GoogleService{}
	MyMemory = MyMemoryService{}

	DeepL = DeepLService{
		key: os.Getenv("DEEPL_AUTHKEY"),
	}

	Yandex = YandexService{
		sid:    atomic.NewString(""),
		urlSID: "https://translate.yandex.ru",
		urlAPI: "https://translate.yandex.net/api/v1/tr.json/translate",
	}
)
