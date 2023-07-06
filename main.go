package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/handybots/quizzoro/handler"
	"github.com/handybots/quizzoro/handler/middleware"
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"

	"github.com/demget/clickrus"
	"github.com/sirupsen/logrus"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	tdb, err := opentdb.Load()
	if err != nil {
		log.Fatal(err)
	}

	hook, err := clickrus.NewHook(clickHouseConfig)
	if err != nil {
		log.Fatal(err)
	}

	logrus.AddHook(hook)
	logrus.SetOutput(os.Stdout)

	h := handler.New(handler.Handler{
		Layout: lt,
		Bot:    b,
		DB:     db,
		TDB:    tdb,
	})

	b.OnError = h.OnError
	b.Use(middleware.Logger(logrus.StandardLogger(), h.LoggerFields))
	b.Use(lt.Middleware("uk"))

	b.Handle("/start", h.OnStart)
	b.Handle("/settings", h.OnSettings)
	b.Handle(tele.OnAddedToGroup, h.OnStart)
	b.Handle(tele.OnPollAnswer, h.OnPollAnswer)
	b.Handle(lt.Callback("privacy"), h.OnPrivacy)
	b.Handle(lt.Callback("category"), h.OnCategory)
	b.Handle(lt.Callback("start"), h.OnCategories)
	b.Handle(lt.Callback("stats"), h.OnStats)

	rl := b.Group()
	rl.Use(middleware.RateLimit(10 * time.Second))
	rl.Handle("/skip", h.OnSkip)
	rl.Handle(lt.Callback("skip"), h.OnSkip)

	b.Handle("/stop", h.OnStop)
	b.Handle(lt.Callback("stop"), h.OnStop)
	b.Handle(lt.Callback("bad_quiz"), h.OnBadQuiz)
	b.Handle(lt.Callback("bad_answers"), h.OnBadAnswers)

	b.Start()
}

var clickHouseConfig = clickrus.Config{
	Addr:    os.Getenv("CLICKHOUSE_URL"),
	Table:   "quizzoro.logs",
	Columns: []string{"event", "user_id"},
}
