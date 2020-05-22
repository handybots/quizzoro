package handler

import (
	"database/sql"
	"math/rand"
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
	if state == storage.StateDefault {
		return nil
	}

	cache, err := h.db.Users.Cache(m.Sender.ID)
	if err != nil {
		return err
	}

	// TODO: remember skipped poll to avoid duplicates

	_ = h.b.Delete(tb.StoredMessage{
		MessageID: cache.MessageID,
		ChatID:    m.Chat.ID,
	})

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
	if state == storage.StateDefault {
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
	if err == nil {
		_, err := h.b.Forward(user, avail)
		return err
	}
	if err != sql.ErrNoRows {
		return err
	}

	trivia, err := opentdb.RandomTrivia(categories[category])
	if err != nil {
		return err
	}

	cached, err := h.db.Polls.ByQuestion(category, trivia.Question)
	if err == nil {
		// TODO: check if user had already passed the quiz
		_, err := h.b.Forward(user, cached)
		return err
	}
	if err != sql.ErrNoRows {
		return err
	}

	var (
		correct    int
		answers    []string
		moderation = "moderation_en"
	)

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
		moderation = "moderation"
		answers = trueFalseAnswers
		if trivia.CorrectAnswer == "False" {
			correct = 1
		}
	}

	question, err := translateText(trivia.Question)
	if err != nil {
		return err
	}

	msg, err := h.sendPoll(question, answers, correct)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(msg, h.b.InlineMarkup(moderation, msg.Poll.ID))
	if err != nil {
		return err
	}

	quiz := storage.Poll{
		ID:          msg.Poll.ID,
		MessageID:   strconv.Itoa(msg.ID),
		ChatID:      int64(h.conf.QuizzesChat),
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

func (h Handler) sendPoll(q string, a []string, i int) (*tb.Message, error) {
	poll := &tb.Poll{
		Type:          tb.PollQuiz,
		CorrectOption: i,
		Question:      q,
	}
	poll.AddOptions(a...)

	msg, err := h.b.Send(h.conf.QuizzesChat, poll)
	if err != nil {
		return nil, err
	}

	return msg, err
}

func shuffleStrings(s []string) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func shuffleWithCorrect(s []string, correct string) (ind int) {
	shuffleStrings(s)
	for i, a := range s {
		if a == correct {
			ind = i
			break
		}
	}
	return
}
