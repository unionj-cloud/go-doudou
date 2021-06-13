package astutils

import (
	"bytes"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"io/ioutil"
	"os"
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
	type args struct {
		src  []byte
		file string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				src:  []byte(code),
				file: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := pathutils.Abs("testfiles/output.go")
			FixImport(tt.args.src, file)
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
		})
	}
}