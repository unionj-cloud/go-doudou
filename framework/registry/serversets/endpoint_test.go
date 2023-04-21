package serversets

import (
	"errors"
	"testing"
	"time"
)

func TestEndpointSameName(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	watch, err := set.Watch()
	if err != nil {
		panic(err)
	}
	defer watch.Close()

	// add first endpoint
	ep1, err := set.RegisterEndpoint("localhost", 1, nil)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}
	defer ep1.Close()

	<-watch.Event()

	ep2, err := set.RegisterEndpoint("localhost", 1, nil)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}
	defer ep2.Close()

	<-watch.Event()

	if l := len(watch.Endpoints()); l != 2 {
		t.Errorf("should have 2 servers, got %v", watch.Endpoints())
	}
}

func TestEndpointPingInitiallyUp(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	watch, err := set.Watch()
	if err != nil {
		panic(err)
	}
	defer watch.Close()

	alive := false // will switch to up on first check
	pingFunction := func() error {
		alive = !alive

		if !alive {
			return errors.New("not alive")
		}

		return nil
	}

	// add endpoint
	ep, err := set.RegisterEndpoint("localhost", 1, pingFunction)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}
	defer ep.Close()

	<-watch.Event()

	if l := len(watch.Endpoints()); l != 1 {
		t.Errorf("should have one endpoint, got %v", watch.Endpoints())
	}

	// in a second should ping to false
	<-watch.Event()
	if l := len(watch.Endpoints()); l != 0 {
		t.Errorf("should have zero endpoints, got %v", watch.Endpoints())
	}

	// third ping should be true
	<-watch.Event()
	if l := len(watch.Endpoints()); l != 1 {
		t.Errorf("should have one endpoint, got %v", watch.Endpoints())
	}
}

func TestEndpointPingInitiallyDown(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	watch, err := set.Watch()
	if err != nil {
		panic(err)
	}
	defer watch.Close()

	alive := true // will switch to down on first check
	pingFunction := func() error {
		alive = !alive

		if !alive {
			return errors.New("not alive")
		}

		return nil
	}

	// add endpoint
	ep, err := set.RegisterEndpoint("localhost", 1, pingFunction)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}
	defer ep.Close()

	// down so doesn't register itself
	// <-watch.Event()

	if l := len(watch.Endpoints()); l != 0 {
		t.Errorf("should have zero endpoints, got %v", watch.Endpoints())
	}

	// in a second should ping to true
	<-watch.Event()
	if l := len(watch.Endpoints()); l != 1 {
		t.Errorf("should have one endpoint, got %v", watch.Endpoints())
	}

	// third ping should be false
	<-watch.Event()
	if l := len(watch.Endpoints()); l != 0 {
		t.Errorf("should have zero endpoints, got %v", watch.Endpoints())
	}
}

func TestEndpointClosePingRoutine(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})

	ping := 0
	ep, err := set.RegisterEndpoint("localhost", 1, func() error {
		ping++
		return nil
	})

	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}

	ep.Close()

	time.Sleep(3 * time.Second)
	if ping > 1 {
		t.Errorf("ping should be closed, called %d times", ping)
	}
}

func TestEndpointMultipleClose(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})

	ep, err := set.RegisterEndpoint("localhost", 1, nil)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}

	ep.Close()
	ep.Close()
	ep.Close()
	ep.Close()
}
