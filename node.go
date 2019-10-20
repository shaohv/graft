package graft

type RaftNode struct {
	raft      *Raft
	LocalAddr string
	OtherNode []string
}

func NewRaftNode(addr string, others []string) *RaftNode {
	return &RaftNode{
		raft:      NewRaft(),
		LocalAddr: addr,
		OtherNode: others,
	}
}

func connectPeers(others []string) {

}

func (r *RaftNode) Put(k string, v string) error {
	return nil
}

func (r *RaftNode) Get(k string) (string, error) {
	return "", nil
}
