package handler

import (
	"math/rand"
	"strconv"

	tele "gopkg.in/tucnak/telebot.v3"
)

// randCategory returns random opentdb category code
// by given string alias.
func randCategory(s string) int {
	category := categories[s]
	return category[rand.Intn(len(category))]
}

// shuffleStrings shuffles slice of strings in one line.
func shuffleStrings(s []string) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

// shuffleWithCorrect shuffles strings and returns index
// of the matched correct answer.
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

// parseChatID converts Recipient string to the integer ID.
func parseChatID(to tele.Recipient) (n int64) {
	n, _ = strconv.ParseInt(to.Recipient(), 10, 64)
	return
}

// fromGroup checks if the given Recipient or chat ID is negative.
// Supports int64 and tb.Recipient, in other cases, returns false.
func fromGroup(to interface{}) bool {
	switch to := to.(type) {
	case int64:
		return to < 0
	case tele.Recipient:
		return to.Recipient()[0] == '-'
	default:
		return false
	}
}
