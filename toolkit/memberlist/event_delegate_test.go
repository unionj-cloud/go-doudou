package memberlist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelEventDelegate_NotifyJoin(t *testing.T) {
	eventCh := make(chan NodeEvent, 1)
	c := &ChannelEventDelegate{
		Ch: eventCh,
	}
	c.NotifyJoin(&Node{
		Name: "",
		Addr: "test",
		Port: 7946,
	})
	ne := <-eventCh
	require.Equal(t, NodeJoin, ne.Event)
}

func TestChannelEventDelegate_NotifyLeave(t *testing.T) {
	eventCh := make(chan NodeEvent, 1)
	c := &ChannelEventDelegate{
		Ch: eventCh,
	}
	c.NotifyLeave(&Node{
		Name: "",
		Addr: "test",
		Port: 7946,
	})
	ne := <-eventCh
	require.Equal(t, NodeLeave, ne.Event)
}

func TestChannelEventDelegate_NotifyUpdate(t *testing.T) {
	eventCh := make(chan NodeEvent, 1)
	c := &ChannelEventDelegate{
		Ch: eventCh,
	}
	c.NotifyUpdate(&Node{
		Name: "",
		Addr: "test",
		Port: 7946,
	})
	ne := <-eventCh
	require.Equal(t, NodeUpdate, ne.Event)
}

func TestChannelEventDelegate_NotifyWeight(t *testing.T) {
	eventCh := make(chan NodeEvent, 1)
	c := &ChannelEventDelegate{
		Ch: eventCh,
	}
	c.NotifyWeight(&Node{
		Name: "",
		Addr: "test",
		Port: 7946,
	})
	ne := <-eventCh
	require.Equal(t, NodeWeight, ne.Event)
}

func TestChannelEventDelegate_NotifySuspectSateChange(t *testing.T) {
	eventCh := make(chan NodeEvent, 1)
	c := &ChannelEventDelegate{
		Ch: eventCh,
	}
	c.NotifySuspectSateChange(&Node{
		Name: "",
		Addr: "test",
		Port: 7946,
	})
	ne := <-eventCh
	require.Equal(t, NodeSuspect, ne.Event)
}
