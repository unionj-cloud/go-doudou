package caches

import (
	"sync/atomic"
	"testing"

	"github.com/wubin1989/gorm"
)

func TestQueryTask_GetId(t *testing.T) {
	task := &queryTask{
		id: "myId",
		db: nil,
		queryCb: func(db *gorm.DB) {
		},
	}

	if task.GetId() != "myId" {
		t.Error("GetId on queryTask returned an unexpected value")
	}
}

func TestQueryTask_Run(t *testing.T) {
	var inc int32
	task := &queryTask{
		id: "myId",
		db: nil,
		queryCb: func(db *gorm.DB) {
			atomic.AddInt32(&inc, 1)
		},
	}

	task.Run()

	if atomic.LoadInt32(&inc) != 1 {
		t.Error("Run on queryTask was expected to execute the callback specified once")
	}
}
