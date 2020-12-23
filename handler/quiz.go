package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	"unicode"

	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnSkip(c tele.Context) error {
	to := c.Chat()

	state, err := h.db.Users.State(to.ID)
	if err != nil {
		return err
	}
	if state == storage.StateDefault {
		return h.sendNotStarted(c)
	}

	cache, err := h.db.Users.Cache(to.ID)
	if err != nil {
		return err
	}

	has, err := h.db.Users.HasPoll(to.ID, cache.OrigPollID)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	_ = h.b.Delete(tele.StoredMessage{
		MessageID: cache.LastMessageID,
		ChatID:    to.ID,
	})

	return h.sendQuiz(c, cache.LastCategory)
}

func (h handler) OnStop(c tele.Context) error {
	to := c.Chat()

	state, err := h.db.Users.State(to.ID)
	if err != nil {
		return err
	}
	if state == storage.StateDefault {
		return h.sendNotStarted(c)
	}

	cache, err := h.db.Users.Cache(to.ID)
	if err != nil {
		return err
	}

	h.b.Delete(tele.StoredMessage{
		MessageID: cache.LastMessageID,
		ChatID:    to.ID,
	})

	if err := h.sendNotStarted(c); err != nil {
		return err
	}

	return h.db.Users.Update(to.ID, storage.User{
		State: storage.StateDefault,
	})
}

func (h handler) sendNotStarted(c tele.Context) error {
	_, err := h.b.Send(
		c.Chat(),
		h.lt.TextLocale("ru", "not_started"),
		h.menuMarkup(c),
		tele.ModeHTML,
	)
	return err
}

func (h handler) sendQuiz(c tele.Context, category string) error {
	var to tele.Recipient
	if c.PollAnswer() != nil {
		to = c.Sender()
	} else {
		to = c.Chat()
	}

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

	categories := h.lt.Get("categories").Strings(category)
	poll, err := h.db.Polls.Available(chatID, categories)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: No polls available
			return nil
		} else {
			return err
		}
	}

	if poll.PollID == "" {
		var (
			correct    int
			answers    = make([]string, len(poll.AnswersEng))
			moderation = "moderation_en"
		)

		if len(poll.AnswersEng) > 2 { // opentdb.Multiple
			shuffleStrings(poll.AnswersEng)
			for i, a := range poll.AnswersEng {
				if a == poll.CorrectEng {
					correct = i
				}

				if i != 0 {
					time.Sleep(500 * time.Millisecond)
				}

				tr, err := translateText(a)
				if err != nil {
					return err
				}

				rs := []rune(tr)
				rs[0] = unicode.ToUpper(rs[0])
				answers[i] = string(rs)
			}
		} else { // opentdb.TrueFalse
			moderation = "moderation"
			answers = h.lt.Strings("true_false")
			if poll.CorrectEng == "False" {
				correct = 1
			}
		}

		question, err := translateText(poll.QuestionEng)
		if err != nil {
			return err
		}

		msg, err := h.sendPoll(h.lt.ChatID("quizzes_chat"), question, answers, correct)
		if err != nil {
			return err
		}

		_, err = h.b.EditReplyMarkup(msg, h.lt.MarkupLocale("ru", moderation, msg.Poll.ID))
		if err != nil {
			h.OnError(err, c)
		}

		poll.PollID = msg.Poll.ID
		poll.MessageID = strconv.Itoa(msg.ID)
		poll.ChatID = h.lt.Int64("quizzes_chat")
		poll.Question = question
		poll.Answers = answers
		poll.Correct = answers[correct]

		if err := h.db.Polls.Update(poll); err != nil {
			return err
		}
	}

	var (
		msg     *tele.Message
		answers []string
		correct string
	)

	if privacy {
		if poll.IsEng {
			answers = poll.AnswersEng
			correct = poll.CorrectEng
		} else {
			answers = poll.Answers
			correct = poll.Correct
		}

		correct := shuffleWithCorrect(answers, correct)
		msg, err = h.sendPoll(to, poll.Question, answers, correct)
	} else {
		msg, err = h.b.Forward(to, poll)
	}
	if err != nil {
		return err
	}

	if fromGroup(chatID) {
		h.prepareGroupPoll(c)
	}

	cache := storage.UserCache{
		OrigPollID:    poll.ID,
		LastPollID:    msg.Poll.ID,
		LastMessageID: strconv.Itoa(msg.ID),
		LastCategory:  category,
	}

	return h.db.Users.Update(chatID, storage.User{
		State:     storage.StateQuiz,
		UserCache: cache,
	})
}

func (h handler) sendPoll(to tele.Recipient, q string, a []string, i int) (*tele.Message, error) {
	poll := &tele.Poll{
		Type:          tele.PollQuiz,
		CorrectOption: i,
		Question:      q,
	}

	if to != h.lt.ChatID("quizzes_chat") && fromGroup(to) {
		poll.OpenPeriod = int(h.lt.Duration("open_period").Seconds())
	}

	poll.AddOptions(a...)
	return h.b.Send(to, poll)
}

func (h handler) prepareGroupPoll(c tele.Context) {
	var (
		to     = c.Chat()
		chatID = parseChatID(to)
		stop   = errors.New("stop")
	)

	f := func() error {
		user, err := h.db.Users.ByID(chatID)
		if err != nil {
			return err
		}
		if user.State == storage.StateDefault {
			return stop
		}

		has, err := h.db.Users.HasPoll(chatID, user.OrigPollID)
		if err != nil {
			return err
		}
		if !has {
			if err := h.sendNotStarted(c); err != nil {
				return err
			}
			return h.db.Users.Update(chatID, storage.User{
				State: storage.StateDefault,
			})
		}

		return h.sendQuiz(c, user.LastCategory)
	}

	time.AfterFunc(h.lt.Duration("open_period"), func() {
		log.Println(chatID, "sendGroupPoll: next poll")
		if err := f(); err != nil && err != stop {
			h.OnError(fmt.Errorf("sendGroupPoll: %v", err), c)
		}
	})
}
