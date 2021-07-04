package esutils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEs_GetByID(t *testing.T) {
	es, terminator := setupSubTest()
	defer terminator()
	doc, _ := es.GetByID(context.Background(), "9seTXHoBNx091WJ2QCh5")
	assert.NotZero(t, doc)
}
