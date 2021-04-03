package templateutils

import (
	"bytes"
	"github.com/unionj-cloud/go-doudou/kit/constants"
	"strings"
	"text/template"
	"time"
)

func FormatTime(t time.Time) string {
	return t.Format(constants.FORMAT)
}

func BoolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func Eval(t *template.Template) func(string, interface{}) (string, error) {
	return func(name string, v interface{}) (string, error) {
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf, name, v)
		return buf.String(), err
	}
}

func TrimSuffix(suffix, v string) string {
	return strings.TrimSuffix(strings.TrimSpace(v), suffix)
}
