package translate

import (
	"os"
)

type Translator interface {
	Translate(from, to, text string) (string, error)
}

var (
	Google = GoogleService{}

	Yandex = YandexService{}

	MyMemory = MyMemoryService{}

	JustTranslate = JustTranslateService{}

	DeepL = DeepLService{
		key: os.Getenv("DEEPL_AUTHKEY"),
	}
)
