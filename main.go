package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/handler"
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"

	"github.com/demget/clickrus"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/tucnak/telebot.v3"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	hook, err := clickrus.NewHook(clickrus.Config{
		Addr:    os.Getenv("CLICKHOUSE_URL"),
		Table:   "quizzoro.logs",
		Columns: []string{"date", "time", "level", "message", "event", "user_id", "chat_id"},
	})
	if err != nil {
		log.Fatal(err)
	}

	logrus.AddHook(hook)
	logrus.SetOutput(os.Stdout)

	tmpl := &tele.TemplateText{
		Dir:        "data",
		DelimLeft:  "${",
		DelimRight: "}",
	}

	pref, err := tele.NewSettingsYAML("bot.yml", tmpl)
	if err != nil {
		log.Fatal(err)
	}
	pref.Token = os.Getenv("TOKEN")

	b, err := tele.NewBot(pref)
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

	tdb, err := opentdb.Load()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(handler.Config{
		Conf: conf,
		Bot:  b,
		DB:   db,
		TDB:  tdb,
	})

	b.OnError = h.OnError

	b.Handle("/start", h.OnStart)
	b.Handle("/settings", h.OnSettings)
	b.Handle("/stop", h.OnStop)
	b.Handle(tele.OnAddedToGroup, h.OnStart)
	b.Handle(tele.OnPollAnswer, h.OnPollAnswer)
	b.Handle(b.Button("start"), h.OnCategories)
	b.Handle(b.Button("stats"), h.OnStats)
	b.Handle(b.Button("skip"), h.OnSkip)
	b.Handle(b.Button("stop"), h.OnStop)
	b.Handle(b.InlineButton("privacy"), h.OnPrivacy)
	b.Handle(b.InlineButton("category"), h.OnCategory)
	b.Handle(b.InlineButton("bad_quiz"), h.OnBadQuiz)
	b.Handle(b.InlineButton("bad_answers"), h.OnBadAnswers)

	// b.Poller = tele.NewMiddlewarePoller(b.Poller, h.Middleware) // todo:
	b.Start()
}
