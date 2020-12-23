package main

import (
	"fmt"
	"log"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"
)

var categories = map[int]string{
	9:  "general",
	10: "books",
	11: "film",
	12: "music",
	13: "theatre",
	14: "television",
	15: "games",
	16: "games",
	17: "science",
	18: "computers",
	19: "math",
	20: "mythology",
	21: "sports",
	22: "geography",
	23: "history",
	24: "politics",
	25: "art",
	26: "celebrities",
	27: "animals",
	28: "vehicles",
	29: "comics",
	30: "computers",
	31: "anime",
	32: "cartoon",
}

func main() {
	db, err := storage.Open(os.Getenv("MYSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	tdb, err := opentdb.New()
	if err != nil {
		log.Fatal(err)
	}

	var polls []storage.Poll
	for cat, catName := range categories {
		stats, err := tdb.Stats(cat)
		if err != nil {
			log.Fatal(err)
		}

		amount, left := 50, stats.TotalCount
		for {
			if amount > left {
				amount = left
			}

			trivias, err := tdb.Trivias(cat, amount)
			if err != nil {
				log.Fatal(err)
			}

			for _, tr := range trivias {
				polls = append(polls, storage.Poll{
					Category:    catName,
					Difficulty:  tr.Difficulty,
					QuestionEng: tr.Question,
					CorrectEng:  tr.CorrectAnswer,
					AnswersEng:  append(tr.IncorrectAnswers, tr.CorrectAnswer),
				})
			}

			left -= amount
			if left == 0 {
				break
			}
		}

		fmt.Println("DONE", catName)
	}

	insert := sq.Insert("polls").
		Columns(
			"category",
			"difficulty",
			"question_eng",
			"correct_eng",
			"answers_eng",
			"question",
			"correct",
			"answers",
		)

	for _, poll := range polls {
		insert = insert.Values(
			poll.Category,
			poll.Difficulty,
			poll.QuestionEng,
			poll.CorrectEng,
			poll.AnswersEng,
			"", "", "[]",
		)
	}

	_, err = insert.RunWith(db).Exec()
	if err != nil {
		log.Fatal(err)
	}
}
