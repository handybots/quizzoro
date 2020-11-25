package handler

import (
	"strconv"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h Handler) OnBadQuiz(c tele.Context) error {
	return h.onBadQuiz(c)
}

func (h Handler) onBadQuiz(c tele.Context) error {
	if err := h.db.Polls.Delete(c.Data()); err != nil {
		return err
	}
	return c.Delete()
}

func (h Handler) OnBadAnswers(c tele.Context) error {
	return h.onBadAnswers(c)
}

func (h Handler) onBadAnswers(c tele.Context) error {
	cached, err := h.db.Polls.ByID(c.Data())
	if err != nil {
		return err
	}

	answers := cached.AnswersEng
	correct := shuffleWithCorrect(answers, cached.CorrectEng)

	msg, err := h.sendPoll(h.conf.QuizzesChat, cached.Question, answers, correct)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(msg,
		h.lt.Markup(c, "moderation", msg.Poll.ID))
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
	cached.IsEng = true
	cached.MessageID = strconv.Itoa(msg.ID)

	if err := h.db.Polls.Create(cached); err != nil {
		return err
	}

	return c.Delete()
}
