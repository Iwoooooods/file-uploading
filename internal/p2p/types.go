package p2p

type Peer struct {
	ID string
	Address string
}

type Node struct {
	ID string
	Address string
	Peers map[string]*Peer
}

func NewNode(id, address string) *Node {
	return &Node{
		ID: id,
		Address: address,
		Peers: make(map[string]*Peer),
	}
}
