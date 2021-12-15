package cache

import pb "cache/cache/cachepb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
//找节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}


// PeerGetter is the interface that must be implemented by a peer.
//找到节点，找值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}