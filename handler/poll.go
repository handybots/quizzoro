package handler

import (
	"database/sql"

	"github.com/demget/quizzorobot/handler/tracker"
	"github.com/demget/quizzorobot/storage"
	tb "github.com/demget/telebot"
)

func (h Handler) OnPollAnswer(pa *tb.PollAnswer) {
	if err := h.onPollAnswer(pa); err != nil {
		h.OnError(pa, err)
	}
}

func (h Handler) onPollAnswer(pa *tb.PollAnswer) error {
	if len(pa.Options) == 0 {
		return nil
	}

	var (
		chatID int64
	)
	user, err := h.db.Users.ByPollID(pa.PollID)
	if err == sql.ErrNoRows {
		chatID = int64(pa.User.ID)
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
		return h.db.Users.AddPoll(int64(pa.User.ID), poll)
	}

	tracker.Data.Del(chatID)
	return h.sendQuiz(&pa.User, cache.LastCategory)
}
