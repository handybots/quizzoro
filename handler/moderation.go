package handler

import (
	"strconv"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnBadQuiz(c tele.Context) error {
	if err := h.db.Polls.Delete(c.Data()); err != nil {
		return err
	}
	return c.Delete()
}

func (h handler) OnBadAnswers(c tele.Context) error {
	cached, err := h.db.Polls.ByID(c.Data())
	if err != nil {
		return err
	}

	answers := cached.AnswersEng
	correct := shuffleWithCorrect(answers, cached.CorrectEng)

	msg, err := h.sendPoll(h.lt.ChatID("quizzes_chat"), cached.Question, answers, correct)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(msg, h.lt.Markup(c, "moderation", msg.Poll.ID))
	if err != nil {
		return err
	}

	// NOTE:
	// Here, I reassign current poll's id to the new one,
	// so the user will receive the poll again, but with
	// proper english answers. For now, I decided to keep
	// such behaviour to avoid conflicts with unique PollID
	// manipulations in OnPollAnswer handler.

	cached.PollID = msg.Poll.ID
	cached.IsEng = true
	cached.MessageID = strconv.Itoa(msg.ID)

	if err := h.db.Polls.Update(cached); err != nil {
		return err
	}

	return c.Delete()
}
