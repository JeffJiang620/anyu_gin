package trans

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
)

var cannotSetNilTrans = errors.New("can not set nil trans")

var trans ut.Translator

func Trans() ut.Translator {
	return trans
}

func SetTrans(t ut.Translator) error {

	if t == nil {
		return cannotSetNilTrans
	}

	trans = t
	return nil
}
