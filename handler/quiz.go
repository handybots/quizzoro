package handler

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/demget/quizzorobot/opentdb"
	"github.com/demget/quizzorobot/storage"

	tb "github.com/demget/telebot"
)

func (h Handler) OnSkip(m *tb.Message) {
	if err := h.onSkip(m); err != nil {
		h.OnError(m, err)
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
		return err
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
		h.OnError(m, err)
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

	if err := h.sendCategories(m.Sender); err != nil {
		return err
	}

	return h.db.Users.Update(m.Sender.ID, storage.User{
		State: storage.StateDefault,
	})
}

func (h Handler) sendQuiz(user *tb.User, category string) error {
	avail, err := h.db.Polls.Available(user.ID, category)
	if err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		var correct int
		shuffleStrings(avail.Answers)
		for i, a := range avail.Answers {
			if a == avail.Correct {
				correct = i
				break
			}
		}

		_, err = h.sendTriviaPoll(avail.Question, avail.Answers, correct)
		return err
	}

	trivia, err := opentdb.RandomTrivia(categories[category])
	if err != nil {
		return err
	}

	cached, err := h.db.Polls.ByQuestion(category, trivia.Question)
	if err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		// TODO: check if user had already passed the quiz
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

	question, err := translateText(trivia.Question)
	if err != nil {
		return err
	}

	msg, err := h.sendTriviaPoll(question, answers, correct)
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

func (h Handler) sendTriviaPoll(q string, a []string, i int) (*tb.Message, error) {
	poll := &tb.Poll{
		Type:          tb.PollQuiz,
		CorrectOption: i,
		Question:      q,
	}
	poll.AddOptions(a...)

	// TODO: Replace with tb.ChatID
	msg, err := h.b.Send(&tb.Chat{ID: h.conf.QuizzesChat}, poll)
	if err != nil {
		return nil, err
	}

	_, err = h.b.EditReplyMarkup(msg, h.b.InlineMarkup("moderation", msg.Poll.ID))
	if err != nil {
		return nil, err
	}

	return msg, err
}
