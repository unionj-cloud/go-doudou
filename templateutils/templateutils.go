package templateutils

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
	"text/template"
)

func String(tmpl string, data interface{}) (string, error) {
	var (
		sqlBuf bytes.Buffer
		err    error
		tpl    *template.Template
	)
	tpl = template.Must(template.New(filepath.Base(tmpl)).ParseFiles(tmpl))
	if err = tpl.Execute(&sqlBuf, data); err != nil {
		return "", errors.Wrap(err, "error returned from calling tpl.Execute")
	}
	return strings.TrimSpace(sqlBuf.String()), nil
}

func StringBlock(tmpl string, block string, data interface{}) (string, error) {
	var (
		sqlBuf bytes.Buffer
		err    error
		tpl    *template.Template
	)
	tpl = template.Must(template.New(filepath.Base(tmpl)).ParseFiles(tmpl))
	if err = tpl.ExecuteTemplate(&sqlBuf, block, data); err != nil {
		return "", errors.Wrap(err, "error returned from calling tpl.ExecuteTemplate")
	}
	return strings.TrimSpace(sqlBuf.String()), nil
}

func StringBlockMysql(tmpl string, block string, data interface{}) (string, error) {
	var (
		sqlBuf  bytes.Buffer
		err     error
		tpl     *template.Template
		funcMap map[string]interface{}
	)
	tpl = template.New(filepath.Base(tmpl))
	funcMap = make(map[string]interface{})
	funcMap["FormatTime"] = FormatTime
	funcMap["BoolToInt"] = BoolToInt
	funcMap["Eval"] = Eval(tpl)
	funcMap["TrimSuffix"] = TrimSuffix
	tpl = template.Must(tpl.Funcs(funcMap).Funcs(template.FuncMap(sprig.FuncMap())).ParseFiles(tmpl))
	if err = tpl.ExecuteTemplate(&sqlBuf, block, data); err != nil {
		return "", errors.Wrap(err, "error returned from calling tpl.ExecuteTemplate")
	}
	return strings.TrimSpace(sqlBuf.String()), nil
}
