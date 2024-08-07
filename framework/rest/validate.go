package rest

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
)

var validate = validator.New()
var translator ut.Translator

func GetValidate() *validator.Validate {
	return validate
}

func GetTranslator() ut.Translator {
	return translator
}

func SetTranslator(trans ut.Translator) {
	translator = trans
}

func handleValidationErr(err error) error {
	if err == nil {
		return nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	translations := errs.Translate(translator)
	var errmsgs []string
	for _, v := range translations {
		errmsgs = append(errmsgs, v)
	}
	return errors.New(strings.Join(errmsgs, ", "))
}

func ValidateStruct(value interface{}) error {
	return handleValidationErr(validate.Struct(value))
}

func ValidateVar(value interface{}, tag, param string) error {
	if stringutils.IsNotEmpty(param) {
		return errors.Wrap(handleValidationErr(validate.Var(value, tag)), param)
	}
	return handleValidationErr(validate.Var(value, tag))
}
