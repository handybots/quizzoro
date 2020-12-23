package opentdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	Multiple  = "multiple"
	TrueFalse = "boolean"
)

type Session struct {
	token string
}

type Trivia struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

type Stats struct {
	Category    int
	TotalCount  int `json:"total_question_count"`
	EasyCount   int `json:"total_easy_question_count"`
	MediumCount int `json:"total_medium_question_count"`
	HardCount   int `json:"total_hard_question_count"`
}

// Load loads session from disk, if it exists, otherwise creates a new one.
// Use Load() instead of New() to automatically save created session.
func Load() (*Session, error) {
	const path = "opentdb.session"

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		session, err := New()
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(path, []byte(session.token), os.ModePerm)
		return session, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	session := &Session{token: string(data)}
	go session.startPingWorker(time.Hour)
	return session, nil
}

// New requests a new session.
func New() (*Session, error) {
	const url = "https://opentdb.com/api_token.php?command=request"

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code    int    `json:"response_code"`
		Message string `json:"response_message"`
		Token   string `json:"token"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("opentdb: %s (%d)", result.Message, result.Code)
	}

	session := &Session{token: result.Token}
	go session.startPingWorker(time.Hour)
	return session, nil
}

func (s Session) Stats(category int) (*Stats, error) {
	const url = "https://opentdb.com/api_count.php?category=%d"

	resp, err := http.DefaultClient.Get(fmt.Sprintf(url, category))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Stats Stats `json:"category_question_count"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	result.Stats.Category = category
	return &result.Stats, nil
}

func (s Session) Trivia(category int) (*Trivia, error) {
	trivias, err := s.Trivias(category, 1)
	if err != nil {
		return nil, err
	}
	return &trivias[0], nil
}

func (s Session) Trivias(category, amount int) ([]Trivia, error) {
	const url = "https://opentdb.com/api.php?amount=%d&category=%d&token=%s"

	resp, err := http.DefaultClient.Get(fmt.Sprintf(url, amount, category, s.token))
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
	if len(result.Trivia) == 0 {
		return nil, errors.New("opentdb: no trivia was found")
	}

	for i, trivia := range result.Trivia {
		trivia.Question = html.UnescapeString(trivia.Question)
		trivia.CorrectAnswer = html.UnescapeString(trivia.CorrectAnswer)

		for i, a := range trivia.IncorrectAnswers {
			trivia.IncorrectAnswers[i] = html.UnescapeString(a)
		}

		result.Trivia[i] = trivia
	}

	return result.Trivia, nil
}

func (s Session) startPingWorker(d time.Duration) {
	t := time.NewTicker(d)
	for range t.C {
		_, err := s.Trivia(9) // general
		if err != nil {
			log.Println("opentdb/ping:", err)
		}
	}
}
