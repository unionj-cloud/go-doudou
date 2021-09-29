package templateutils

import (
	"bytes"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
	"text/template"
)

// String return result of calling template Execute as string
func String(tmplname, tmpl string, data interface{}) (string, error) {
	var (
		sqlBuf bytes.Buffer
		err    error
		tpl    *template.Template
	)
	tpl = template.Must(template.New(tmplname).Parse(tmpl))
	if err = tpl.Execute(&sqlBuf, data); err != nil {
		return "", errors.Wrap(err, "error returned from calling tpl.Execute")
	}
	return strings.TrimSpace(sqlBuf.String()), nil
}

// StringBlock return result of calling template Execute as string
func StringBlock(tmplname, tmpl string, block string, data interface{}) (string, error) {
	var (
		sqlBuf bytes.Buffer
		err    error
		tpl    *template.Template
	)
	tpl = template.Must(template.New(tmplname).Parse(tmpl))
	if err = tpl.ExecuteTemplate(&sqlBuf, block, data); err != nil {
		return "", errors.Wrap(err, "error returned from calling tpl.ExecuteTemplate")
	}
	return strings.TrimSpace(sqlBuf.String()), nil
}

// StringBlockMysql return result of calling template Execute as string from template file
func StringBlockMysql(tmpl string, block string, data interface{}) (string, error) {
	var (
		sqlBuf  bytes.Buffer
		err     error
		tpl     *template.Template
		funcMap map[string]interface{}
	)
	tpl = template.New(filepath.Base(tmpl))
	funcMap = make(map[string]interface{})
	funcMap["FormatTime"] = formatTime
	funcMap["BoolToInt"] = boolToInt
	funcMap["Eval"] = eval(tpl)
	funcMap["TrimSuffix"] = trimSuffix
	funcMap["isNil"] = func(t interface{}) bool {
		return t == nil
	}
	tpl = template.Must(tpl.Funcs(funcMap).ParseFiles(tmpl))
	if err = tpl.ExecuteTemplate(&sqlBuf, block, data); err != nil {
		return "", errors.Wrap(err, "error returned from calling tpl.ExecuteTemplate")
	}
	return strings.TrimSpace(sqlBuf.String()), nil
}

// BlockMysql return result of calling template Execute as string from template file
func BlockMysql(tmplname, tmpl string, block string, data interface{}) (string, error) {
	var (
		sqlBuf  bytes.Buffer
		err     error
		tpl     *template.Template
		funcMap map[string]interface{}
	)
	tpl = template.New(tmplname)
	funcMap = make(map[string]interface{})
	funcMap["FormatTime"] = formatTime
	funcMap["BoolToInt"] = boolToInt
	funcMap["Eval"] = eval(tpl)
	funcMap["TrimSuffix"] = trimSuffix
	funcMap["isNil"] = func(t interface{}) bool {
		return t == nil
	}
	tpl = template.Must(tpl.Funcs(funcMap).Parse(tmpl))
	if err = tpl.ExecuteTemplate(&sqlBuf, block, data); err != nil {
		return "", errors.Wrap(err, "error returned from calling tpl.ExecuteTemplate")
	}
	return strings.TrimSpace(sqlBuf.String()), nil
}
