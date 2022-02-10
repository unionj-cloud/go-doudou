package esutils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEs_Count(t *testing.T) {
	es := setupSubTest("test_count")
	count, _ := es.Count(context.Background(), nil)
	assert.EqualValues(t, 3, count)
}
