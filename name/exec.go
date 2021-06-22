package name

import (
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
)

type Name struct {
	File      string
	Strategy  string
	Omitempty bool
}

const (
	lowerCamelStrategy = "lowerCamel"
	snakeStrategy      = "snake"
)

// https://github.com/iancoleman/strcase
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

	newcode, err := astutils.RewriteJsonTag(receiver.File, receiver.Omitempty, convert)
	if err != nil {
		panic(err)
	}
	astutils.FixImport([]byte(newcode), receiver.File)
}
