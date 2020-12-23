package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h handler) OnSettings(c tele.Context) error {
	if c.Message().FromGroup() {
		return nil
	}

	privacy, err := h.db.Users.Privacy(c.Chat().ID)
	if err != nil {
		return err
	}

	return c.Send(
		h.lt.Text(c, "settings"),
		h.lt.Markup(c, "privacy", privacy),
		tele.ModeHTML,
	)
}

func (h handler) OnPrivacy(c tele.Context) error {
	defer c.Respond(&tele.CallbackResponse{
		Text: h.lt.Text(c, "privacy"),
	})

	privacy, err := h.db.Users.InvertPrivacy(c.Chat().ID)
	if err != nil {
		return err
	}

	markup := h.lt.Markup(c, "privacy", privacy)
	_, err = h.b.EditReplyMarkup(c.Message(), markup)
	return err
}
