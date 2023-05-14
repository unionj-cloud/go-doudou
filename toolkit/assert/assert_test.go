package assert_test

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"testing"
)

func init() {
	zlogger.Logger = zlogger.Logger.With().Caller().Logger()
}

func TestTrue(t *testing.T) {
	var m map[string]interface{}
	assert.True(m != nil, "m should not be nil")
}
