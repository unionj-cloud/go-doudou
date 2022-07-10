package name

import (
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
)

// Name wraps config properties for name command
type Name struct {
	File      string
	Strategy  string
	Omitempty bool
	Form      bool
}

const (
	lowerCamelStrategy = "lowerCamel"
	snakeStrategy      = "snake"
)

// Exec rewrites the json tag of each field of all structs in the file as snake case or lower camel case.
// Unexported or ignored fields will be skipped.
func (receiver Name) Exec() {
	if stringutils.IsEmpty(receiver.File) {
		panic(errors.New("file flag should not be empty"))
	}

	var convert func(string) string
	switch receiver.Strategy {
	case lowerCamelStrategy:
		convert = strcase.ToLowerCamel
	case snakeStrategy:
		convert = strcase.ToSnake
	default:
		panic(errors.New(`unknown strategy. currently only support "lowerCamel" and "snake"`))
	}

	newcode, err := astutils.RewriteTag(astutils.RewriteTagConfig{
		File:        receiver.File,
		Omitempty:   receiver.Omitempty,
		ConvertFunc: convert,
		Form:        receiver.Form,
	})
	if err != nil {
		panic(err)
	}
	astutils.FixImport([]byte(newcode), receiver.File)
}
