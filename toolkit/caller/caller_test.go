package caller_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"testing"
)

func TestCaller_String(t *testing.T) {
	c := caller.NewCaller()
	require.Equal(t, "called from github.com/unionj-cloud/go-doudou/toolkit/caller_test.TestCaller_String on /Users/wubin1989/workspace/cloud/go-doudou/toolkit/caller/caller_test.go#10",
		c.String())
}
