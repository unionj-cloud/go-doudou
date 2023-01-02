package memberlist_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	memmock "github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist/mock"
	"reflect"
	"testing"
)

func TestMemberlistBroadcast_Invalidates(t *testing.T) {
	m1 := memberlist.NewMemberlistBroadcast("test", nil, nil)
	m2 := memberlist.NewMemberlistBroadcast("foo", nil, nil)

	if m1.Invalidates(m2) || m2.Invalidates(m1) {
		t.Fatalf("unexpected invalidation")
	}

	if !m1.Invalidates(m1) {
		t.Fatalf("expected invalidation")
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	b := memmock.NewMockBroadcast(ctrl)
	ok := m1.Invalidates(b)
	require.False(t, ok)
}

func TestMemberlistBroadcast_Message(t *testing.T) {
	m1 := memberlist.NewMemberlistBroadcast("test", []byte("test"), nil)
	msg := m1.Message()
	if !reflect.DeepEqual(msg, []byte("test")) {
		t.Fatalf("messages do not match")
	}
}
