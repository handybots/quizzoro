// +build ignore

package handler

import (
	"log"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (h handler) OnError(v interface{}, err error) {
	var f Fields
	if s, ok := v.(string); ok {
		f = Fields{Event: s}
	} else {
		f = NewFields(v)
	}

	logrus.WithFields(f.Fields()).Error(err)
	log.Printf("%+v\n", errors.WithStack(err))
}
