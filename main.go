package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/demget/quizzorobot/bot"
	"github.com/demget/quizzorobot/handler"
	"github.com/demget/quizzorobot/storage"

	"github.com/bshuster-repo/logrus-logstash-hook"
	tb "github.com/demget/telebot"
	"github.com/sirupsen/logrus"
)

func init() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal(err)
	}

	f := logrustash.DefaultFormatter(logrus.Fields{"app": "quizzorobot"})
	logrus.AddHook(logrustash.New(conn, f))

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

	db, err := storage.Connect(os.Getenv("MYSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
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
	b.Handle(b.Button("stats"), h.OnStats)
	b.Handle(b.InlineButton("category"), h.OnCategory)
	b.Handle(b.InlineButton("bad_quiz"), h.OnBadQuiz)

	b.Poller = tb.NewMiddlewarePoller(b.Poller, h.Middleware)
	b.Start()
}
