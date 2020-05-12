package translate

type Translator interface {
	Translate(from, to, text string) string
}
