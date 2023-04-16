package serversets

import (
	"sort"
	"testing"
)

func TestWatchSortEndpoints(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})

	watch, err := set.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer watch.Close()

	ep1, err := set.RegisterEndpoint("localhost", 1002, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ep1.Close()
	<-watch.Event()

	ep2, err := set.RegisterEndpoint("localhost", 1001, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ep2.Close()
	<-watch.Event()

	ep3, err := set.RegisterEndpoint("localhost", 1003, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ep3.Close()
	<-watch.Event()

	endpoints := watch.Endpoints()
	if len(endpoints) != 3 {
		t.Errorf("should have 3 endpoint, got %v", endpoints)
	}

	if !sort.StringsAreSorted(endpoints) {
		t.Errorf("endpoint list should be sorted, got %v", endpoints)
	}
}

func TestWatchUpdateEndpoints(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})

	watch, err := set.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer watch.Close()

	ep1, err := set.RegisterEndpoint("localhost", 1002, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ep1.Close()
	<-watch.Event()

	conn, _, err := set.connectToZookeeper()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	eps, err := watch.updateEndpoints(conn, []string{MemberPrefix + "random"})
	if err != nil {
		t.Fatalf("should not have error, got %v", err)
	}

	if len(eps) != 0 {
		t.Errorf("should not have any endpoints, got %v", eps)
	}
}

func TestWatchIsClosed(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	watch, err := set.Watch()
	if err != nil {
		t.Fatal(err)
	}

	watch.Close()

	if watch.IsClosed() == false {
		t.Error("should say it's closed right after we close it")
	}
}

func TestWatchMultipleClose(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	watch, err := set.Watch()
	if err != nil {
		t.Fatal(err)
	}

	watch.Close()
	watch.Close()
	watch.Close()
}

func TestWatchTriggerEvent(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	watch, err := set.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer watch.Close()

	watch.triggerEvent()
	watch.triggerEvent()
	watch.triggerEvent()
	watch.triggerEvent()
}
