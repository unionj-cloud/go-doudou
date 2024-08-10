package form

import (
	"reflect"
	"testing"

	. "github.com/go-playground/assert/v2"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//
//
// go test -cpuprofile cpu.out
// ./validator.test -test.bench=. -test.cpuprofile=cpu.prof
// go tool pprof validator.test cpu.prof
//
//
// go test -memprofile mem.out

func TestDecoderMultipleSimultaniousParseStructRequests(t *testing.T) {

	sc := newStructCacheMap()

	type Struct struct {
		Array []int
	}

	proceed := make(chan struct{})

	var test Struct

	sv := reflect.ValueOf(test)
	typ := sv.Type()

	for i := 0; i < 200; i++ {
		go func() {
			<-proceed
			s := sc.parseStruct(ModeImplicit, sv, typ, "form")
			NotEqual(t, s, nil)
		}()
	}

	close(proceed)
}
