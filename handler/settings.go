package handler

import tb "github.com/demget/telebot"

func (h Handler) OnSettings(m *tb.Message) {
	if m.FromGroup() {
		return
	}
	if err := h.onSettings(m); err != nil {
		h.OnError(m, err)
	}
}

func (h Handler) OnPrivacy(c *tb.Callback) {
	if err := h.onPrivacy(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onSettings(m *tb.Message) error {
	privacy, err := h.db.Users.Privacy(m.Chat.ID)
	if err != nil {
		return err
	}

	_, err = h.b.Send(
		m.Sender,
		h.b.Text("privacy"),
		h.b.InlineMarkup("privacy", privacy),
		tb.ModeHTML)
	return err
}

func (h Handler) onPrivacy(c *tb.Callback) error {
	defer h.b.Respond(c, &tb.CallbackResponse{
		Text: h.b.String("privacy"),
	})

	privacy, err := h.db.Users.InvertPrivacy(c.Message.Chat.ID)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(c.Message,
		h.b.InlineMarkup("privacy", privacy))
	return err
}
