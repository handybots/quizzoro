package handler

import (
	"log"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnError(err error, c tele.Context) {
	if c != nil {
		log.Println(c.Sender().Recipient(), err)
	} else {
		log.Println(err)
	}
}
