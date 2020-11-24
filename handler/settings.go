package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h Handler) OnSettings(c tele.Context) error {
	if m.FromGroup() {
		return
	}
	if err := h.onSettings(m); err != nil {
		h.OnError(m, err)
	}
}

func (h Handler) OnPrivacy(c tele.Context) error {
	if err := h.onPrivacy(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onSettings(c tele.Context) error {
	privacy, err := h.db.Users.Privacy(c.Chat().ID)
	if err != nil {
		return err
	}

	return c.Send(
		h.b.Text("privacy"),
		h.b.InlineMarkup("privacy", privacy),
		tb.ModeHTML,
	)
}

func (h Handler) onPrivacy(c tele.Context) error {
	defer h.b.Respond(c.Callback(), &tele.CallbackResponse{
		Text: h.b.String("privacy"),
	})

	privacy, err := h.db.Users.InvertPrivacy(c.Chat().ID)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(c.Message(),
		h.b.InlineMarkup("privacy", privacy))
	return err
}
