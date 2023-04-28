package serversets

import (
	"reflect"
	"testing"
)

const TestServer = "localhost"

// This is the big run through a typical use case of add and remove and make sure it works.
func TestServerSetAddAndRemove(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	watch, err := set.Watch()
	if err != nil {
		panic(err)
	}

	// add first endpoint
	ep1, err := set.RegisterEndpoint("localhost", 1, nil)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D1%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// add second endpoint
	ep2, err := set.RegisterEndpoint("localhost", 2, nil)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D1%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D2%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// add third endpoint
	ep3, err := set.RegisterEndpoint("localhost", 3, nil)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D1%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D2%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D3%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// remove second endpoint
	ep2.Close()

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D1%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D3%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// add fourth endpoint
	ep4, err := set.RegisterEndpoint("localhost", 4, nil)
	if err != nil {
		t.Fatalf("registration failure: %v", err)
	}

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D1%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D3%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D4%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// close and reopen watch
	if watch.EventCount != 5 {
		t.Errorf("event count incorrect, got %d", watch.EventCount)
	}

	watch.Close()

	watch, err = set.Watch()
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D1%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D3%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D4%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// remove third endpoint
	ep3.Close()

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D1%29%2F%25%21s%28%3Cnil%3E%29%3F", "%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D4%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// remove first endpoint
	ep1.Close()

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{"%25%21s%28%3Cnil%3E%29%3A%2F%2Flocalhost%3A%25%21s%28float64%3D4%29%2F%25%21s%28%3Cnil%3E%29%3F"}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	// remove fourth endpoint
	ep4.Close()

	<-watch.Event()
	if !reflect.DeepEqual(watch.Endpoints(), []string{}) {
		t.Errorf("server list incorrect, got %v", watch.Endpoints())
	}

	if watch.EventCount != 3 {
		t.Errorf("event count incorrect, got %d", watch.EventCount)
	}
	watch.Close()
}

func TestBaseZnodePath(t *testing.T) {
	// to verify nothing happens to the default
	path := BaseZnodePath(Test, "gotest")
	if path != "/registry/gotest/providers" {
		t.Errorf("baseznodepath incorrect, got %v", path)
	}
}

func TestServerSetDirectoryPath(t *testing.T) {
	set := New(Test, "gotest", []string{TestServer})
	path := set.directoryPath()

	// should just be a pass through to BaseZnodePath
	if path != BaseZnodePath(Test, "gotest") {
		t.Errorf("directory path incorrect, got %v", path)
	}
}

func TestSplitPaths(t *testing.T) {
	path := "/discovery/test/gotest"
	parts := splitPaths(path)
	if !reflect.DeepEqual(parts, []string{"/discovery", "/discovery/test", "/discovery/test/gotest"}) {
		t.Errorf("split not correct, got %v", parts)
	}

	path = "/discovery/test/"
	parts = splitPaths(path)
	if !reflect.DeepEqual(parts, []string{"/discovery", "/discovery/test"}) {
		t.Errorf("split not correct, got %v", parts)
	}

	path = "/"
	parts = splitPaths(path)
	if !reflect.DeepEqual(parts, []string{}) {
		t.Errorf("split not correct, got %v", parts)
	}
}
