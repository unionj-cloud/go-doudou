package caller_test

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"testing"
)

func TestCaller_String(t *testing.T) {
	c := caller.NewCaller()
	fmt.Println(c.String())
}
