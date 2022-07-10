package astutils

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"regexp"
	"strings"
	"unicode"
)

func isExport(field string) bool {
	return unicode.IsUpper([]rune(field)[0])
}

func extractJsonPropName(tag string) string {
	re := regexp.MustCompile(`json:"(.*?)"`)
	if re.MatchString(tag) {
		subs := re.FindAllStringSubmatch(tag, -1)
		return strings.TrimSpace(strings.Split(subs[0][1], ",")[0])
	}
	return ""
}

func extractFormPropName(tag string) string {
	re := regexp.MustCompile(`form:"(.*?)"`)
	if re.MatchString(tag) {
		subs := re.FindAllStringSubmatch(tag, -1)
		return strings.TrimSpace(strings.Split(subs[0][1], ",")[0])
	}
	return ""
}

type RewriteTagConfig struct {
	File        string
	Omitempty   bool
	ConvertFunc func(old string) string
	Form        bool
}

// RewriteTag overwrites json tag by convert function and return formatted source code
func RewriteTag(config RewriteTagConfig) (string, error) {
	file, convert, omitempty, form := config.File, config.ConvertFunc, config.Omitempty, config.Form
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return "", errors.Wrap(err, caller.NewCaller().String())
	}
	re := regexp.MustCompile(`json:"(.*?)"`)
	reForm := regexp.MustCompile(`form:"(.*?)"`)
	astutil.Apply(root, func(cursor *astutil.Cursor) bool {
		return true
	}, func(cursor *astutil.Cursor) bool {
		structSpec, ok := cursor.Node().(*ast.StructType)
		if !ok {
			return true
		}
		for _, field := range structSpec.Fields.List {
			if field.Names == nil {
				continue
			}
			fname := field.Names[0].Name
			if !isExport(fname) {
				continue
			}
			tagValue := convert(field.Names[0].Name)
			jsonTagValue := tagValue
			if omitempty {
				jsonTagValue += ",omitempty"
			}
			jsonTag := fmt.Sprintf(`json:"%s"`, jsonTagValue)

			formTagValue := tagValue
			if omitempty {
				formTagValue += ",omitempty"
			}
			formTag := fmt.Sprintf(`form:"%s"`, formTagValue)

			if field.Tag != nil {
				if re.MatchString(field.Tag.Value) {
					if extractJsonPropName(field.Tag.Value) != "-" {
						field.Tag.Value = re.ReplaceAllLiteralString(field.Tag.Value, jsonTag)
					}
				} else {
					lastindex := strings.LastIndex(field.Tag.Value, "`")
					if lastindex < 0 {
						panic(errors.New("call LastIndex() error"))
					}
					field.Tag.Value = field.Tag.Value[:lastindex] + fmt.Sprintf(" %s`", jsonTag)
				}
				if form {
					if reForm.MatchString(field.Tag.Value) {
						if extractFormPropName(field.Tag.Value) != "-" {
							field.Tag.Value = re.ReplaceAllLiteralString(field.Tag.Value, formTag)
						}
					} else {
						lastindex := strings.LastIndex(field.Tag.Value, "`")
						if lastindex < 0 {
							panic(errors.New("call LastIndex() error"))
						}
						field.Tag.Value = field.Tag.Value[:lastindex] + fmt.Sprintf(" %s`", formTag)
					}
				}
			} else {
				if form {
					field.Tag = &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`%s %s`", jsonTag, formTag),
					}
				} else {
					field.Tag = &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`%s`", jsonTag),
					}
				}
			}
		}
		return true
	})
	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, root)
	if err != nil {
		return "", fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.String(), nil
}
