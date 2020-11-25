package handler

import (
	"database/sql"

	"github.com/handybots/quizzoro/handler/tracker"
	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h Handler) OnPollAnswer(c tele.Context) error {
	return h.onPollAnswer(c)
}

func (h Handler) onPollAnswer(c tele.Context) error {
	pa := c.PollAnswer()

	if len(pa.Options) == 0 {
		return nil
	}

	var (
		chatID int64
	)
	user, err := h.db.Users.ByPollID(pa.PollID)
	if err == sql.ErrNoRows {
		chatID = int64(pa.Sender.ID)
	} else if err == nil {
		chatID = user.ID
	} else {
		return err
	}

	state, err := h.db.Users.State(chatID)
	if err != nil {
		return err
	}
	if state == storage.StateDefault {
		return nil
	}

	cache, err := h.db.Users.Cache(chatID)
	if err != nil {
		return err
	}

	correct, err := h.db.Polls.CorrectAnswer(cache.OrigPollID)
	if err != nil {
		return err
	}

	has, err := h.db.Users.HasPoll(chatID, cache.OrigPollID)
	if err != nil {
		return err
	}

	poll := storage.PassedPoll{
		PollID:  cache.OrigPollID,
		Correct: pa.Options[0] == correct,
	}
	if !has {
		if err := h.db.Users.AddPoll(chatID, poll); err != nil {
			return err
		}
	}
	if fromGroup(chatID) {
		return h.db.Users.AddPoll(int64(pa.Sender.ID), poll)
	}

	tracker.Data.Del(chatID)
	return h.sendQuiz(c, cache.LastCategory)
}
