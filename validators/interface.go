package validators

import (
	"fmt"
	"github.com/JeffJiang620/anyu_gin/trans"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

type Validator interface {
	FailedText() string
	CallValidationEvenIfNull() bool
	TagName() string
	Validate(fl validator.FieldLevel) bool
}

func initValidateTrans(locale string, v *validator.Validate) (err error) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	zhCh := zh.New()
	enUs := en.New()
	uni := ut.New(zhCh, zhCh, enUs)
	var ok bool

	t, ok := uni.GetTranslator(locale)
	if !ok {
		return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
	}

	err = trans.SetTrans(t)
	if err != nil {
		return err
	}
	switch locale {
	case "zh":
		err = zhTrans.RegisterDefaultTranslations(v, trans.Trans())
	default:
		err = enTrans.RegisterDefaultTranslations(v, trans.Trans())
	}

	return
}

func RegisterValidator(locale string, validators ...Validator) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err = initValidateTrans(locale, v)
		if err != nil {
			return
		}

		for _, vl := range validators {
			if vl != nil {
				_ = v.RegisterValidation(vl.TagName(), vl.Validate, vl.CallValidationEvenIfNull())
				_ = v.RegisterTranslation(vl.TagName(), trans.Trans(), func(ut ut.Translator) error {
					return ut.Add(vl.TagName(), vl.FailedText(), true)
				}, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T(vl.TagName(), fe.Field())
					return t
				})
			}
		}

		return
	}
	return
}
