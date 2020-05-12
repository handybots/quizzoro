package handler

import (
	"log"

	"github.com/demget/quizzorobot/storage"

	tb "github.com/demget/telebot"
)

func (h Handler) OnPollAnswer(pa *tb.PollAnswer) {
	if err := h.onPollAnswer(pa); err != nil {
		log.Println(err)
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
	if state != storage.StateQuiz {
		return nil
	}

	correct, err := h.db.Quizzes.CorrectAnswer(pa.PollID)
	if err != nil {
		return err
	}

	poll := storage.UserPoll{
		ID:      pa.PollID,
		Correct: pa.Options[0] == correct,
	}
	if err := h.db.Users.AddPoll(pa.User.ID, poll); err != nil {
		return err
	}

	cache, err := h.db.Users.Cache(pa.User.ID)
	if err != nil {
		log.Println(err)
	}

	return h.sendQuiz(&pa.User, cache.Category)
}
