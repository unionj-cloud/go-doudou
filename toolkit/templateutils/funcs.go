package templateutils

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
)

func formatTime(t time.Time) string {
	return t.Format(constants.FORMAT)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Eval(t *template.Template) func(string, interface{}) (string, error) {
	return func(name string, v interface{}) (string, error) {
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf, name, v)
		return buf.String(), err
	}
}

func trimSuffix(suffix, v string) string {
	return strings.TrimSuffix(strings.TrimSpace(v), suffix)
}

func hasPrefix(v, prefix string) bool {
	return strings.HasPrefix(strings.TrimSpace(v), prefix)
}
