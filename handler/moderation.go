package handler

import (
	"strconv"

	tb "github.com/demget/telebot"
)

func (h Handler) OnBadQuiz(c *tb.Callback) {
	if err := h.onBadQuiz(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onBadQuiz(c *tb.Callback) error {
	if err := h.db.Polls.Delete(c.Data); err != nil {
		return err
	}
	return h.b.Delete(c.Message)
}

func (h Handler) OnBadAnswers(c *tb.Callback) {
	if err := h.onBadAnswers(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onBadAnswers(c *tb.Callback) error {
	cached, err := h.db.Polls.ByID(c.Data)
	if err != nil {
		return err
	}

	answers := cached.AnswersEng
	correct := shuffleWithCorrect(answers, cached.CorrectEng)

	msg, err := h.sendPoll(cached.Question, answers, correct)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(msg,
		h.b.InlineMarkup("moderation", msg.Poll.ID))
	if err != nil {
		return err
	}

	if err := h.db.Polls.Delete(cached.ID); err != nil {
		return err
	}

	// NOTE:
	// Here, we're reassigning current poll's id to the new one,
	// so the user will receive the poll again, but with
	// proper english answers. For now, I decided to keep
	// such behaviour to avoid conflicts with unique PollID
	// manipulations in OnPollAnswer handler.

	cached.ID = msg.Poll.ID
	cached.MessageID = strconv.Itoa(msg.ID)

	if err := h.db.Polls.Create(cached); err != nil {
		return err
	}

	return h.b.Delete(c.Message)
}
