package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h Handler) OnSettings(c tele.Context) error {
	if c.Message().FromGroup() {
		return nil
	}
	return h.onSettings(c)
}

func (h Handler) OnPrivacy(c tele.Context) error {
	return h.onPrivacy(c)
}

func (h Handler) onSettings(c tele.Context) error {
	privacy, err := h.db.Users.Privacy(c.Chat().ID)
	if err != nil {
		return err
	}

	return c.Send(
		h.lt.Text(c, "privacy"),
		h.lt.Markup(c, "privacy", privacy),
		tele.ModeHTML,
	)
}

func (h Handler) onPrivacy(c tele.Context) error {
	defer h.b.Respond(c.Callback(), &tele.CallbackResponse{
		Text: h.lt.String("privacy"),
	})

	privacy, err := h.db.Users.InvertPrivacy(c.Chat().ID)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(c.Message(),
		h.lt.Markup(c, "privacy", privacy))
	return err
}
