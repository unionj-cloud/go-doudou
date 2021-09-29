package templateutils

import (
	"bytes"
	"github.com/unionj-cloud/go-doudou/constants"
	"strings"
	"text/template"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format(constants.FORMAT)
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func eval(t *template.Template) func(string, interface{}) (string, error) {
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
