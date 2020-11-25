package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/demget/clickrus"
	"github.com/handybots/quizzoro/handler"
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
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

	lt, err := layout.New("bot.yml")
	if err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(lt.Settings())
	if err != nil {
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
		Layout: lt,
		Bot:    b,
		DB:     db,
		TDB:    tdb,
	})

	b.OnError = h.OnError
	b.Use(lt.Middleware("ru", h.LocaleFunc))

	b.Handle("/start", h.OnStart)
	b.Handle("/settings", h.OnSettings)
	b.Handle("/stop", h.OnStop)
	b.Handle(tele.OnAddedToGroup, h.OnStart)
	b.Handle(tele.OnPollAnswer, h.OnPollAnswer)
	b.Handle(lt.Callback("privacy"), h.OnPrivacy)
	b.Handle(lt.Callback("category"), h.OnCategory)
	b.Handle(lt.Callback("bad_quiz"), h.OnBadQuiz)
	b.Handle(lt.Callback("bad_answers"), h.OnBadAnswers)

	for _, loc := range []string{"ru"} {
		b.Handle(lt.ButtonLocale(loc, "start"), h.OnCategories)
		b.Handle(lt.ButtonLocale(loc, "stats"), h.OnStats)
		b.Handle(lt.ButtonLocale(loc, "skip"), h.OnSkip)
		b.Handle(lt.ButtonLocale(loc, "stop"), h.OnStop)
	}

	// b.Poller = tele.NewMiddlewarePoller(b.Poller, h.Middleware) // todo:
	b.Start()
}
