package graft

import (
	"log"

	"google.golang.org/grpc"
)

type RaftNode struct {
	raft      *Raft
	LocalAddr string
	OtherNode []string
	RPCServer *grpc.Server
}

func NewRaftNode(addr string, others []string) *RaftNode {
	if len(others) < 2 || (len(others)%2) != 0 {
		log.Println("node num not right!")
		return nil
	}
	var r = &RaftNode{
		LocalAddr: addr,
		OtherNode: others,
	}

	// Start rpc server
	var servC = make(chan *grpc.Server, 1)
	go StartGrpcServer(addr, servC)
	r.RPCServer = <-servC

	// load data ...

	// Connect other node
	var peers = make([]*PeerCli, len(others)+1)
	for _, o := range others {
		cli, err := NewPeer(o)
		if err != nil {
			log.Println("conn fail", err)
			return nil
		}
		peers = append(peers, cli)
	}

	var metaStor = NewMemStor(0, 0)
	r.raft = NewRaft(peers, metaStor)
	return r
}

func (r *RaftNode) Put(k string, v string) error {
	return nil
}

func (r *RaftNode) Get(k string) (string, error) {
	return "", nil
}
