package main

import (
	"flag"
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

var file = flag.String("file", "", "absolute path of vo file")
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
