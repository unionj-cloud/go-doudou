package memberlist_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	memmock "github.com/unionj-cloud/go-doudou/framework/memberlist/mock"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testCountingWriter struct {
	t        *testing.T
	numCalls *int32
}

func (tw testCountingWriter) Write(p []byte) (n int, err error) {
	atomic.AddInt32(tw.numCalls, 1)
	if !strings.Contains(string(p), "memberlist: Error accepting TCP connection") {
		tw.t.Error("did not receive expected log message")
	}
	tw.t.Log("countingWriter:", string(p))
	return len(p), nil
}

func TestAddress_String(t *testing.T) {
	addr := &memberlist.Address{
		Addr: "localhost:7946",
		Name: "testNode",
	}
	require.Equal(t, "testNode (localhost:7946)", addr.String())

	addr = &memberlist.Address{
		Addr: "localhost:7946",
	}
	require.Equal(t, "localhost:7946", addr.String())
}

func Test_shimNodeAwareTransport_WriteToAddress(t *testing.T) {
	Convey("", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		nat := memmock.NewMockTransport(ctrl)
		addr := memberlist.Address{
			Addr: "localhost:7946",
			Name: "testNode",
		}
		now := time.Now()
		msg := []byte("test message")
		nat.
			EXPECT().
			WriteTo(msg, addr.Addr).
			AnyTimes().
			Return(now, nil)

		s := memberlist.NewshimNodeAwareTransport(nat)
		sendAt, err := s.WriteToAddress(msg, addr)
		So(err, ShouldBeNil)
		So(sendAt, ShouldResemble, now)
	})
}

type emptyReadNetConn struct {
	net.Conn
}

func (c *emptyReadNetConn) Read(b []byte) (n int, err error) {
	return 0, io.EOF
}

func (c *emptyReadNetConn) Close() error {
	return nil
}

func Test_shimNodeAwareTransport_DialAddressTimeout(t1 *testing.T) {
	Convey("", t1, func() {
		ctrl := gomock.NewController(t1)
		defer ctrl.Finish()
		nat := memmock.NewMockTransport(ctrl)
		addr := memberlist.Address{
			Addr: "localhost:7946",
			Name: "testNode",
		}
		nat.
			EXPECT().
			DialTimeout(addr.Addr, 10*time.Second).
			AnyTimes().
			Return(&emptyReadNetConn{}, nil)

		s := memberlist.NewshimNodeAwareTransport(nat)
		_, err := s.DialAddressTimeout(addr, 10*time.Second)
		So(err, ShouldBeNil)
	})
}
