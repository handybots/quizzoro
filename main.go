package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/demget/quizzorobot/bot"
	"github.com/demget/quizzorobot/handler"
	"github.com/demget/quizzorobot/storage"

	tb "github.com/demget/telebot"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	tmpl := &tb.TemplateText{
		Dir:        "data",
		DelimLeft:  "${",
		DelimRight: "}",
	}

	pref, err := tb.NewSettingsYAML("bot.yml", tmpl)
	if err != nil {
		log.Fatal(err)
	}
	pref.Token = os.Getenv("TOKEN")

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	var conf bot.Config
	if err := b.Vars(&conf); err != nil {
		log.Fatal(err)
	}

	db, err := storage.Connect(
		os.Getenv("MONGODB_NAME"),
		os.Getenv("MONGODB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(handler.Config{
		Conf: conf,
		Bot:  b,
		DB:   db,
	})

	b.Handle("/start", h.OnStart)
	b.Handle("/skip", h.OnSkip)
	b.Handle("/stop", h.OnStop)
	b.Handle(tb.OnPollAnswer, h.OnPollAnswer)
	b.Handle(b.Button("start"), h.OnCategories)
	b.Handle(b.InlineButton("category"), h.OnCategory)
	b.Handle(b.InlineButton("bad_quiz"), h.OnBadQuiz)

	b.Poller = tb.NewMiddlewarePoller(b.Poller, h.Middleware)
	b.Start()
}
