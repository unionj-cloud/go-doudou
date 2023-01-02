package memberlist

/*
The broadcast mechanism works by maintaining a sorted list of messages to be
sent out. When a message is to be broadcast, the retransmit count
is set to zero and appended to the queue. The retransmit count serves
as the "priority", ensuring that newer messages get sent first. Once
a message hits the retransmit limit, it is removed from the queue.

Additionally, older entries can be invalidated by new messages that
are contradictory. For example, if we send "{suspect M1 inc: 1},
then a following {alive M1 inc: 2} will invalidate that message
*/

type weightBroadcast struct {
	node string
	msg  []byte
}

func NewWeightBroadcast(node string, msg []byte) *weightBroadcast {
	return &weightBroadcast{
		node,
		msg,
	}
}

func (b *weightBroadcast) Invalidates(other Broadcast) bool {
	// Check if that broadcast is a weight type
	mb, ok := other.(*weightBroadcast)
	if !ok {
		return false
	}

	// Invalidates any message about the same node
	return b.node == mb.node
}

func (b *weightBroadcast) Message() []byte {
	return b.msg
}

func (b *weightBroadcast) Finished() {
}
