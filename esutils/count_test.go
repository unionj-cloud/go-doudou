package esutils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEs_Count(t *testing.T) {
	es, terminator := setupSubTest()
	defer terminator()
	count, _ := es.Count(context.Background(), nil)
	assert.EqualValues(t, 3, count)
}
