package middleware

import (
	"sync"
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
)

func RateLimit(d time.Duration) tele.MiddlewareFunc {
	var (
		mu sync.Mutex
		dm = make(map[string]int64)
	)

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			rec := c.Recipient().Recipient()

			mu.Lock()
			last, ok := dm[rec]
			mu.Unlock()

			if ok && time.Now().Sub(time.Unix(last, 0)) < d {
				return nil
			}

			mu.Lock()
			dm[rec] = time.Now().Unix()
			mu.Unlock()

			return next(c)
		}
	}
}
