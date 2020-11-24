package handler

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h Handler) OnSkip(c tele.Context) error {
	return h.onSkip(c)
}

func (h Handler) OnStop(c tele.Context) error {
	return h.onStop(c)
}

func (h Handler) onSkip(c tele.Context) error {
	m := c.Message()

	if m.FromGroup() {
		return nil
	}

	state, err := h.db.Users.State(m.Chat.ID)
	if err != nil {
		return err
	}
	if state == storage.StateDefault {
		return h.sendNotStarted(m.Chat)
	}

	cache, err := h.db.Users.Cache(m.Chat.ID)
	if err != nil {
		return err
	}

	_ = h.b.Delete(tb.StoredMessage{
		MessageID: cache.LastMessageID,
		ChatID:    m.Chat.ID,
	})

	return h.sendQuiz(m.Chat, cache.LastCategory)
}

func (h Handler) onStop(c tele.Context) error {
	m := c.Message()

	state, err := h.db.Users.State(m.Chat.ID)
	if err != nil {
		return err
	}
	if state == storage.StateDefault {
		return h.sendNotStarted(m.Chat)
	}

	cache, err := h.db.Users.Cache(m.Chat.ID)
	if err != nil {
		return err
	}

	_ = h.b.Delete(tb.StoredMessage{
		MessageID: cache.LastMessageID,
		ChatID:    m.Chat.ID,
	})

	if err := c.Send(
		h.b.Text("start", m.Chat),
		h.b.Markup("menu"),
		tele.ModeHTML,
	); err != nil {
		return err
	}

	return h.db.Users.Update(m.Chat.ID, storage.User{
		State: storage.StateDefault,
	})
}

func (h Handler) sendNotStarted(to tele.Recipient) error {
	_, err := h.b.Send(to,
		h.b.Text("not_started"),
		h.b.Markup("menu"),
		tb.ModeHTML)
	return err
}

func (h Handler) sendQuiz(to tele.Recipient, category string) error {
	var (
		chatID  = parseChatID(to)
		privacy = true
	)

	if !fromGroup(chatID) {
		p, err := h.db.Users.Privacy(chatID)
		if err != nil {
			return err
		}
		privacy = p
	}

	avail, err := h.db.Polls.Available(chatID, category)
	if err == nil {
		var (
			msg     *tb.Message
			answers []string
			correct string
		)
		if privacy {
			if avail.IsEng {
				answers = avail.AnswersEng
				correct = avail.CorrectEng
			} else {
				answers = avail.Answers
				correct = avail.Correct
			}
			correct := shuffleWithCorrect(answers, correct)
			msg, err = h.sendPoll(to, avail.Question, answers, correct)
		} else {
			msg, err = h.b.Forward(to, avail)
		}
		if err != nil {
			return err
		}

		if fromGroup(chatID) {
			h.prepareGroupPoll(to)
		}

		cache := storage.UserCache{
			OrigPollID:    avail.ID,
			LastPollID:    msg.Poll.ID,
			LastMessageID: strconv.Itoa(msg.ID),
			LastCategory:  category,
		}
		return h.db.Users.Update(chatID, storage.User{
			State:     storage.StateQuiz,
			UserCache: cache,
		})
	}
	if err != sql.ErrNoRows {
		return err
	}

TRIVIA:
	trivia, err := h.tdb.Trivia(randCategory(category))
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
		answers = bot.TrueFalseAnswers
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
		msg, err = h.sendPoll(to, question, answers, correct)
	} else {
		msg, err = h.b.Forward(to, msg)
	}
	if err != nil {
		return err
	}

	if fromGroup(chatID) {
		h.prepareGroupPoll(to)
	}

	cache := storage.UserCache{
		OrigPollID:    pollID,
		LastPollID:    msg.Poll.ID,
		LastMessageID: strconv.Itoa(msg.ID),
		LastCategory:  category,
	}
	return h.db.Users.Update(chatID, storage.User{
		State:     storage.StateQuiz,
		UserCache: cache,
	})
}

func (h Handler) sendPoll(to tele.Recipient, q string, a []string, i int) (*tele.Message, error) {
	poll := &tele.Poll{
		Type:          tele.PollQuiz,
		CorrectOption: i,
		Question:      q,
	}

	if to != h.conf.QuizzesChat && fromGroup(to) {
		poll.OpenPeriod = int(h.conf.OpenPeriod)
	}

	poll.AddOptions(a...)
	return h.b.Send(to, poll)
}

func (h Handler) prepareGroupPoll(to tele.Recipient) {
	chatID := parseChatID(to)

	f := func() error {
		cache, err := h.db.Users.Cache(chatID)
		if err != nil {
			return err
		}

		has, err := h.db.Users.HasPoll(chatID, cache.OrigPollID)
		if err != nil {
			return err
		}
		if !has {
			return h.db.Users.Update(chatID, storage.User{
				State: storage.StateDefault,
			})
		}

		return h.sendQuiz(to, cache.LastCategory)
	}

	time.AfterFunc(h.conf.OpenPeriod*time.Second, func() {
		if err := f(); err != nil {
			h.OnError(fmt.Sprintf("sendGroupPoll(%d)", chatID), err)
		}
	})
}
