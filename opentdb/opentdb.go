package opentdb

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
)

const (
	Multiple  = "multiple"
	TrueFalse = "boolean"
)

type Trivia struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

func RandomTrivia(cat int) (*Trivia, error) {
	url := "https://opentdb.com/api.php?amount=%d&category=%d"
	const amount = 1

	resp, err := http.DefaultClient.Get(fmt.Sprintf(url, amount, cat))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code   int      `json:"response_code"`
		Trivia []Trivia `json:"results"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("opentdb: response code is %d", result.Code)
	}

	trivia := result.Trivia[0]
	trivia.Question = html.UnescapeString(trivia.Question)
	trivia.CorrectAnswer = html.UnescapeString(trivia.CorrectAnswer)

	for i, a := range trivia.IncorrectAnswers {
		trivia.IncorrectAnswers[i] = html.UnescapeString(a)
	}

	return &trivia, nil
}
