package esutils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/constants"
	"testing"
	"time"
)

func TestEs_SaveOrUpdate(t *testing.T) {
	es := setupSubTest("test_saveorupdate")

	id, _ := es.SaveOrUpdate(context.Background(), map[string]interface{}{
		"id":       "9seTXHoBNx091WJ2QCh8",
		"createAt": time.Now().UTC().Format(constants.FORMATES),
		"text":     "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
	})
	assert.EqualValues(t, "9seTXHoBNx091WJ2QCh8", id)
}
