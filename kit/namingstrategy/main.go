package main

import (
	"github.com/unionj-cloud/papilio/kit/astutils"
	"github.com/unionj-cloud/papilio/kit/namingstrategy/strategies"
	"github.com/unionj-cloud/papilio/kit/stringutils"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"time"
)

var file = flag.String("file", "/Users/wubin1989/workspace/cloud/papilio/kit/namingstrategy/example/vo/vos.go", "name of file")
var strategy = flag.String("strategy", "lowerCaseNamingStrategy", "name of strategy")

func main() {
	flag.Parse()
	log.Println(*file)
	log.Println(*strategy)
	if stringutils.IsEmpty(*file) {
		log.Fatal("file flag should not be empty")
	}

	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, *file, nil, 0)
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	ast.Walk(&sc, root)
	fmt.Println(sc.Structs)

	marshalers := strings.TrimRight(*file, ".go") + "_marshaller.go"
	f, err := os.Create(marshalers)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	strategies.Registry[*strategy].Execute(f, struct {
		StructCollector astutils.StructCollector
		Timestamp       time.Time
	}{
		StructCollector: sc,
		Timestamp:       time.Now(),
	})
}
