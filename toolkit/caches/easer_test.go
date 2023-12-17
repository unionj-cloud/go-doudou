package caches

import (
	"sync"
	"testing"
	"time"
)

func TestEase(t *testing.T) {
	t.Run("same queries", func(t *testing.T) {
		queue := &sync.Map{}

		myTask := &mockTask{
			delay:  1 * time.Second,
			expRes: "expect-this",
			id:     "unique-id",
		}
		myDupTask := &mockTask{
			delay:  1 * time.Second,
			expRes: "not-this",
			id:     "unique-id",
		}

		wg := &sync.WaitGroup{}
		wg.Add(2)

		var (
			myTaskRes    *mockTask
			myDupTaskRes *mockTask
		)

		// Both queries will run at the same time, the second one will run half a second later
		go func() {
			myTaskRes = ease(myTask, queue).(*mockTask)
			wg.Done()
		}()
		go func() {
			time.Sleep(500 * time.Millisecond)
			myDupTaskRes = ease(myDupTask, queue).(*mockTask)
			wg.Done()
		}()
		wg.Wait()

		if myTaskRes.actRes != myTaskRes.expRes {
			t.Error("expected first query to be executed")
		}

		if myTaskRes.actRes != myDupTaskRes.actRes {
			t.Errorf("expected same result from both tasks, expected: %s, actual: %s",
				myTaskRes.actRes, myDupTaskRes.actRes)
		}
	})

	t.Run("different queries", func(t *testing.T) {
		queue := &sync.Map{}

		myTask := &mockTask{
			delay:  1 * time.Second,
			expRes: "expect-this",
			id:     "unique-id",
		}
		myDupTask := &mockTask{
			delay:  1 * time.Second,
			expRes: "not-this",
			id:     "other-unique-id",
		}

		wg := &sync.WaitGroup{}
		wg.Add(2)

		var (
			myTaskRes    *mockTask
			myDupTaskRes *mockTask
		)

		// Both queries will run at the same time, the second one will run half a second later
		go func() {
			myTaskRes = ease(myTask, queue).(*mockTask)
			wg.Done()
		}()
		go func() {
			time.Sleep(500 * time.Millisecond)
			myDupTaskRes = ease(myDupTask, queue).(*mockTask)
			wg.Done()
		}()
		wg.Wait()

		if myTaskRes.actRes != myTaskRes.expRes {
			t.Errorf("expected first query to be executed, expected: %s, actual: %s",
				myTaskRes.actRes, myTaskRes.expRes)
		}

		if myTaskRes.actRes == myDupTaskRes.actRes {
			t.Errorf("expected different result from both tasks, expected: %s, actual: %s",
				myTaskRes.actRes, myDupTaskRes.actRes)
		}

		if myDupTaskRes.actRes != myDupTaskRes.expRes {
			t.Error("expected second query to be executed")
		}
	})
}
