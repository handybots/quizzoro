package bot

type Config struct {
	// TODO: Replace with tb.ChatID
	//  (tucnak/telebot@v2.2)
	QuizzesChat int64 `json:"quizzes_chat"`
}

type Random struct {
	Value    int
	Category string
}
