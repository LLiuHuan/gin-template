// Package validation
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 11:09
package validation

import (
	"errors"
	"fmt"

	"github.com/LLiuHuan/gin-template/configs"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	zhTranslation "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	lang := configs.Get().Project.Local
	if lang == configs.ZhCN {
		trans, _ = ut.New(zh.New()).GetTranslator("zh")
		if err := zhTranslation.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), trans); err != nil {
			fmt.Println("validator zh translation error", err)
		}
	}

	if lang == configs.EnUS {
		trans, _ = ut.New(en.New()).GetTranslator("en")
		if err := enTranslation.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), trans); err != nil {
			fmt.Println("validator en translation error", err)
		}
	}
}

func validationError(err error) string {
	var message string
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err.Error()
	}
	for _, e := range validationErrors {
		message += e.Translate(trans) + ";"
	}
	return message
}

func ErrorE(err error) error {
	return errors.New(validationError(err))
}

func Error(err error) (message string) {
	return validationError(err)
}
