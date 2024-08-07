package caller_test

import (
	"fmt"
	"testing"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/caller"
)

func TestCaller_String(t *testing.T) {
	c := caller.NewCaller()
	fmt.Println(c.String())
}
