package astutils

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFixImport(t *testing.T) {
	code := `package main

import (
	"fmt"
"encoding/json"
)

type UserVo struct {
					Id    int
					Name  string
	Phone string
	Dept  string
}

type Page struct {
	PageNo int
Size   int
Items  []UserVo
}

func main() {
	page := Page{
	PageNo: 10,
	Size:   30,
}
b, _ := json.Marshal(page)
fmt.Println(string(b))
}
`
	expect := `package main

import (
	"encoding/json"
	"fmt"
)

type UserVo struct {
	Id    int
	Name  string
	Phone string
	Dept  string
}

type Page struct {
	PageNo int
	Size   int
	Items  []UserVo
}

func main() {
	page := Page{
		PageNo: 10,
		Size:   30,
	}
	b, _ := json.Marshal(page)
	fmt.Println(string(b))
}
`
	file := pathutils.Abs("testdata/output.go")
	FixImport([]byte(code), file)
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte(expect)) {
		t.Error("somewhat bad happen")
	}
}

func TestFixImportPanic(t *testing.T) {
	code := `package main

import (
	"fmt"
"encoding/json"
)

type UserVo
`
	assert.Panics(t, func() {
		FixImport([]byte(code), "")
	})
}

func TestMethodMeta_String(t *testing.T) {
	type fields struct {
		Recv     string
		Name     string
		Params   []FieldMeta
		Results  []FieldMeta
		Comments []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		panic  bool
	}{
		{
			name: "",
			fields: fields{
				Recv: "handler",
				Name: "HandleEvent",
				Params: []FieldMeta{
					{
						Name:     "ctx",
						Type:     "context.Context",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
					{
						Name:     "etype",
						Type:     "int",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
					{
						Name:     "uid",
						Type:     "string",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Results: []FieldMeta{
					{
						Name:     "",
						Type:     "bool",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
					{
						Name:     "",
						Type:     "error",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Comments: nil,
			},
			want: "func (receiver handler) HandleEvent(ctx context.Context, etype int, uid string) (bool, error)",
		},
		{
			name: "",
			fields: fields{
				Recv: "",
				Name: "HandleEvent",
				Params: []FieldMeta{
					{
						Name:     "ctx",
						Type:     "context.Context",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
					{
						Name:     "etype",
						Type:     "int",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
					{
						Name:     "uid",
						Type:     "string",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Results: []FieldMeta{
					{
						Name:     "",
						Type:     "error",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Comments: nil,
			},
			want: "func HandleEvent(ctx context.Context, etype int, uid string) error",
		},
		{
			name: "",
			fields: fields{
				Recv: "",
				Name: "",
				Params: []FieldMeta{
					{
						Name:     "etype",
						Type:     "int",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
					{
						Name:     "uid",
						Type:     "string",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Results: []FieldMeta{
					{
						Name:     "",
						Type:     "error",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Comments: nil,
			},
			want: "func(etype int, uid string) error",
		},
		{
			name: "",
			fields: fields{
				Recv: "PlaceHolder",
				Name: "",
				Params: []FieldMeta{
					{
						Name:     "etype",
						Type:     "int",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
					{
						Name:     "uid",
						Type:     "string",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Results: []FieldMeta{
					{
						Name:     "",
						Type:     "error",
						Tag:      "",
						Comments: nil,
						IsExport: false,
						DocName:  "",
					},
				},
				Comments: nil,
			},
			want:  "",
			panic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mm := MethodMeta{
				Recv:     tt.fields.Recv,
				Name:     tt.fields.Name,
				Params:   tt.fields.Params,
				Results:  tt.fields.Results,
				Comments: tt.fields.Comments,
			}
			if tt.panic {
				assert.Panics(t, func() {
					mm.String()
				})
			} else {
				if got := mm.String(); stringutils.IsNotEmpty(tt.want) && got != tt.want {
					t.Errorf("String() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestVisit(t *testing.T) {
	testDir := pathutils.Abs("./testdata")
	vodir := filepath.Join(testDir, "vo")
	var files []string
	err := filepath.Walk(vodir, Visit(&files))
	if err != nil {
		logrus.Panicln(err)
	}
	assert.Len(t, files, 1)
}

func TestGetMod(t *testing.T) {
	testDir := pathutils.Abs("./testdata")
	_ = os.Chdir(testDir)
	assert.Equal(t, "testdata", GetMod())
}

func TestGetImportPath(t *testing.T) {
	testDir := pathutils.Abs("./testdata")
	_ = os.Chdir(testDir)
	assert.Equal(t, "testdata/vo", GetImportPath(testDir+"/vo"))
}

func TestNewMethodMeta(t *testing.T) {
	file := pathutils.Abs("testdata/cat.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
}

func TestGetImportStatements(t *testing.T) {
	input := `import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
	"{{.VoPackage}}"
	"github.com/pkg/errors"
)`
	ret := GetImportStatements([]byte(input))
	fmt.Println(string(ret))
}

func TestAppendImportStatements(t *testing.T) {
	input := `import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"

	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
	"{{.VoPackage}}"
)

type UsersvcHandlerImpl struct {
	usersvc service.Usersvc
}
`
	ret := AppendImportStatements([]byte(input), []byte(`
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
`))

	ret = AppendImportStatements(ret, []byte(`
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
`))

	ret = AppendImportStatements(ret, []byte(`
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
`))

	expected := `import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"

	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
	"{{.VoPackage}}"

	"github.com/pkg/errors"
)

type UsersvcHandlerImpl struct {
	usersvc service.Usersvc
}
`

	fmt.Println(string(ret))
	require.Equal(t, expected, string(ret))
}

func TestAppendImportStatements1(t *testing.T) {
	input := `import ()`
	ret := AppendImportStatements([]byte(input), []byte(`
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
`))
	expected := `import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)`

	fmt.Println(string(ret))
	require.Equal(t, expected, string(ret))
}

func TestAppendImportStatements2(t *testing.T) {
	input := `import ()`
	ret := AppendImportStatements([]byte(input), []byte(`

`))
	expected := input

	fmt.Println(string(ret))
	require.Equal(t, expected, string(ret))
}

func ExampleGetAnnotations() {
	ret := GetAnnotations(`// <b style="color: red">NEW</b> 删除数据接口（不删数据文件）@role(SUPER_ADMIN)@permission(create,update)这是几个注解@sss()`)
	fmt.Println(ret)
	// Output:
	// [{@role [SUPER_ADMIN]} {@permission [create update]} {@sss []}]
}
