package astutils

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
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

// RewriteJSONTag overwrites json tag by convert function and return formatted source code
func RewriteJSONTag(file string, omitempty bool, convert func(old string) string) (string, error) {
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return "", errors.Wrap(err, "call ParseFile() error")
	}
	re := regexp.MustCompile(`json:"(.*?)"`)
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
			tag := convert(field.Names[0].Name)
			if omitempty {
				tag += ",omitempty"
			}
			tag = fmt.Sprintf(`json:"%s"`, tag)
			if field.Tag != nil {
				if re.MatchString(field.Tag.Value) {
					if extractJsonPropName(field.Tag.Value) != "-" {
						field.Tag.Value = re.ReplaceAllLiteralString(field.Tag.Value, tag)
					}
				} else {
					lastindex := strings.LastIndex(field.Tag.Value, "`")
					if lastindex < 0 {
						panic(errors.New("call LastIndex() error"))
					}
					field.Tag.Value = field.Tag.Value[:lastindex] + fmt.Sprintf(" %s`", tag)
				}
			} else {
				field.Tag = &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("`%s`", tag),
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
