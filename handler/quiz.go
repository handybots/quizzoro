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

func (h Handler) OnStop(m *tb.Message) {
	if err := h.onStop(m); err != nil {
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

	_ = h.b.Delete(tb.StoredMessage{
		MessageID: cache.LastMessageID,
		ChatID:    m.Chat.ID,
	})

	return h.sendQuiz(m.Sender, cache.LastCategory)
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
		MessageID: cache.LastMessageID,
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
	privacy, err := h.db.Users.Privacy(user.ID)
	if err != nil {
		return err
	}

	avail, err := h.db.Polls.Available(user.ID, category)
	if err == nil {
		var (
			msg *tb.Message
		)
		if privacy {
			answers := avail.Answers
			correct := shuffleWithCorrect(answers, avail.Correct)
			msg, err = h.sendPoll(user, avail.Question, answers, correct)
		} else {
			msg, err = h.b.Forward(user, avail)
		}
		if err != nil {
			return err
		}

		cache := storage.UserCache{
			LastPollID:    avail.ID,
			LastMessageID: strconv.Itoa(msg.ID),
			LastCategory:  category,
		}
		return h.db.Users.Update(user.ID, storage.User{
			State:     storage.StateQuiz,
			UserCache: cache,
		})
	}
	if err != sql.ErrNoRows {
		return err
	}

TRIVIA:
	trivia, err := opentdb.RandomTrivia(categories[category])
	if err != nil {
		return err
	}

	_, err = h.db.Polls.ByQuestion(category, trivia.Question)
	if err == nil {
		goto TRIVIA
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

	msg, err := h.sendPoll(h.conf.QuizzesChat, question, answers, correct)
	if err != nil {
		return err
	}

	pollID := msg.Poll.ID
	_, err = h.b.EditReplyMarkup(msg, h.b.InlineMarkup(moderation, pollID))
	if err != nil {
		return err
	}

	poll := storage.Poll{
		ID:          pollID,
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
	if err := h.db.Polls.Create(poll); err != nil {
		return err
	}

	if privacy {
		msg, err = h.sendPoll(user, question, answers, correct)
	} else {
		msg, err = h.b.Forward(user, msg)
	}
	if err != nil {
		return err
	}

	cache := storage.UserCache{
		LastPollID:    pollID,
		LastMessageID: strconv.Itoa(msg.ID),
		LastCategory:  category,
	}
	return h.db.Users.Update(user.ID, storage.User{
		State:     storage.StateQuiz,
		UserCache: cache,
	})
}

func (h Handler) sendPoll(to tb.Recipient,
	q string, a []string, i int) (*tb.Message, error) {

	poll := &tb.Poll{
		Type:          tb.PollQuiz,
		CorrectOption: i,
		Question:      q,
	}

	poll.AddOptions(a...)
	return h.b.Send(to, poll)
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
