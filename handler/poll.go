package handler

import (
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

	state, err := h.db.Users.State(pa.User.ID)
	if err != nil {
		return err
	}
	if state == storage.StateDefault {
		return nil
	}

	cache, err := h.db.Users.Cache(pa.User.ID)
	if err != nil {
		return err
	}

	correct, err := h.db.Polls.CorrectAnswer(pa.PollID)
	if err != nil {
		return err
	}

	poll := storage.PassedPoll{
		PollID:  pa.PollID,
		Correct: pa.Options[0] == correct,
	}
	if err := h.db.Users.AddPoll(pa.User.ID, poll); err != nil {
		return err
	}

	return h.sendQuiz(&pa.User, cache.LastCategory)
}
