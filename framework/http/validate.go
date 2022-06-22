package ddhttp

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"strings"
)

var validate = validator.New()

func GetValidate() *validator.Validate {
	return validate
}

func handleValidationErr(err error) error {
	if err == nil {
		return nil
	}
	var sb strings.Builder
	for i, errItem := range err.(validator.ValidationErrors) {
		sb.WriteString(fmt.Sprintf(`%d. `, i+1))
		if stringutils.IsNotEmpty(errItem.Namespace()) {
			sb.WriteString(fmt.Sprintf(`%s `, errItem.Namespace()))
		}
		sb.WriteString(fmt.Sprintf(`"%s" didn't satisfy the validation rule: "%s`, errItem.Value(), errItem.ActualTag()))
		if stringutils.IsNotEmpty(errItem.Param()) {
			sb.WriteString(fmt.Sprintf(`=%s"`, errItem.Param()))
		} else {
			sb.WriteString(`"`)
		}
		sb.WriteString("\n")
	}
	return errors.New(sb.String())
}

func ValidateStruct(value interface{}) error {
	return handleValidationErr(validate.Struct(value))
}

func ValidateVar(value interface{}, tag string) error {
	return handleValidationErr(validate.Var(value, tag))
}
