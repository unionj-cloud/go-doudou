package name

import (
	"bytes"
	"fmt"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/name/strategies"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"time"
)

type Name struct {
	File      string
	Strategy  string
	Omitempty bool
}

func (receiver Name) Exec() {
	if stringutils.IsEmpty(receiver.File) {
		log.Fatal("file flag should not be empty")
	}

	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, receiver.File, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	ast.Walk(&sc, root)
	fmt.Println(sc.Structs)

	marshalers := strings.TrimSuffix(receiver.File, ".go") + "_marshaller.go"
	f, err := os.Create(marshalers)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var sqlBuf bytes.Buffer
	strategies.Registry[receiver.Strategy].Execute(&sqlBuf, struct {
		StructCollector astutils.StructCollector
		Timestamp       time.Time
		Omitempty       bool
	}{
		StructCollector: sc,
		Timestamp:       time.Now(),
		Omitempty:       receiver.Omitempty,
	})

	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), marshalers)
}
