package handler

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/demget/quizzorobot/opentdb"
	"github.com/demget/quizzorobot/storage"

	gt "github.com/bas24/googletranslatefree"
	tb "github.com/demget/telebot"
)

func (h Handler) OnSkip(m *tb.Message) {
	if err := h.onSkip(m); err != nil {
		log.Println(err)
	}
}

func (h Handler) onSkip(m *tb.Message) error {
	state, err := h.db.Users.State(m.Sender.ID)
	if err != nil {
		return err
	}
	if state != storage.StateQuiz {
		return nil
	}

	cache, err := h.db.Users.Cache(m.Sender.ID)
	if err != nil {
		log.Println(err)
	}

	msg := tb.StoredMessage{
		MessageID: cache.MessageID,
		ChatID:    m.Chat.ID,
	}
	if err := h.b.Delete(msg); err != nil {
		return err
	}

	return h.sendQuiz(m.Sender, cache.Category)
}

func (h Handler) OnStop(m *tb.Message) {
	if err := h.onStop(m); err != nil {
		log.Println(err)
	}
}

func (h Handler) onStop(m *tb.Message) error {
	state, err := h.db.Users.State(m.Sender.ID)
	if err != nil {
		return err
	}
	if state != storage.StateQuiz {
		return nil
	}

	cache, err := h.db.Users.Cache(m.Sender.ID)
	if err != nil {
		return err
	}

	_ = h.b.Delete(tb.StoredMessage{
		MessageID: cache.MessageID,
		ChatID:    m.Chat.ID,
	})

	return h.db.Users.Update(m.Sender.ID, storage.User{
		State: storage.StateDefault,
	})
}

func (h Handler) sendQuiz(user *tb.User, category string) error {
	trivia, err := opentdb.RandomTrivia(categories[category])
	if err != nil {
		return err
	}

	cached, err := h.db.Polls.ByQuestion(category, trivia.Question)
	if err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		_, err := h.forward(user, cached)
		return err
	}

	var correct int
	var answers []string

	if trivia.Type == opentdb.Multiple {
		answers = []string{trivia.CorrectAnswer}
		answers = append(answers, trivia.IncorrectAnswers...)
		shuffleStrings(answers)

		for i, a := range answers {
			if a == trivia.CorrectAnswer {
				correct = i
			}

			tr, err := translateText(a)
			if err != nil {
				return err
			}

			answers[i] = strings.Title(tr)
		}
	} else {
		answers = trueFalseAnswers
		if trivia.CorrectAnswer == "False" {
			correct = 1
		}
	}

	question, err := gt.Translate(trivia.Question, "en", "ru")
	if err != nil {
		return err
	}

	poll := &tb.Poll{
		Type:          tb.PollQuiz,
		CorrectOption: correct,
		Question:      question,
	}
	poll.AddOptions(answers...)

	// TODO: Replace with tb.ChatID
	chat := &tb.Chat{ID: h.conf.QuizzesChat}
	msg, err := h.b.Send(chat, poll)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(msg,
		h.b.InlineMarkup("moderation", msg.Poll.ID))
	if err != nil {
		return err
	}

	quiz := storage.Poll{
		ID:          msg.Poll.ID,
		MessageID:   strconv.Itoa(msg.ID),
		ChatID:      h.conf.QuizzesChat,
		Category:    category,
		Difficulty:  trivia.Difficulty,
		Question:    question,
		Correct:     answers[correct],
		Answers:     answers,
		QuestionEng: trivia.Question,
		CorrectEng:  trivia.CorrectAnswer,
		AnswersEng:  append(trivia.IncorrectAnswers, trivia.CorrectAnswer),
	}
	if err := h.db.Polls.Create(quiz); err != nil {
		return err
	}

	msg, err = h.b.Forward(user, msg)
	if err != nil {
		return err
	}

	cache := storage.UserCache{
		MessageID: strconv.Itoa(msg.ID),
		Category:  category,
	}
	return h.db.Users.Update(user.ID, storage.User{
		State:     storage.StateQuiz,
		UserCache: cache,
	})
}
