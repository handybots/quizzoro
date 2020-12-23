package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/handybots/quizzoro/storage"
	"github.com/handybots/quizzoro/translate"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
)

var tr = translate.JustTranslate

func main() {
	lt, err := layout.New("bot.yml")
	if err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(lt.Settings())
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.Open(os.Getenv("MYSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	var polls []storage.Poll
	const query = `SELECT * FROM polls WHERE poll_id=''`
	if err := db.Select(&polls, query); err != nil {
		log.Fatal(err)
	}

	for _, poll := range polls {
		var (
			correct    int
			moderation = "moderation_en"
			answers    = make([]string, len(poll.AnswersEng))
		)
		if len(poll.AnswersEng) > 2 { // opentdb.Multiple
			for i, a := range poll.AnswersEng {
				if a == poll.CorrectEng {
					correct = i
				}

				answer, err := tr.Translate("en", "ru", a)
				if err != nil {
					log.Fatal(err)
				}

				answers[i] = strings.Title(answer)
			}
		} else { // opentdb.TrueFalse
			moderation = "moderation"
			answers = lt.Strings("true_false")
			if poll.CorrectEng == "False" {
				correct = 1
			}
		}

		question, err := tr.Translate("en", "ru", poll.QuestionEng)
		if err != nil {
			log.Fatal(err)
		}

		msg, err := b.Send(
			lt.ChatID("quizzes_chat"),
			telePoll(question, answers, correct),
			lt.MarkupLocale("ru", moderation, poll.ID),
		)
		if err != nil {
			log.Fatal(err)
		}

		poll.PollID = msg.Poll.ID
		poll.MessageID = strconv.Itoa(msg.ID)
		poll.ChatID = msg.Chat.ID
		poll.Question = question
		poll.Answers = answers
		poll.Correct = answers[correct]

		if err := db.Polls.Update(poll); err != nil {
			log.Fatal(err)
		}

		log.Println(poll.ID, "DONE")
	}
}

func telePoll(q string, a []string, i int) *tele.Poll {
	poll := &tele.Poll{
		Type:          tele.PollQuiz,
		CorrectOption: i,
		Question:      q,
	}
	poll.AddOptions(a...)
	return poll
}
