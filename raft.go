package graft

import (
	"graft/pb"
)

type Raft struct {
	CurrentTerm uint64
	VotedFor    uint32
	Log         []pb.Entries

	CommitIdx   uint64
	LastApplied uint64

	NextIdx  uint64
	MatchIdx uint64

	Peers []PeerCli
}

func NewRaft() *Raft {
	return &Raft{
		CurrentTerm: 0,
		VotedFor:    0,
		Log:         make([]pb.Entries, 1024),

		CommitIdx:   0,
		LastApplied: 0,

		NextIdx:  0,
		MatchIdx: 0,
	}
}
