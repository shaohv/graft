package graft

import (
	"errors"
	"graft/pb"
	"log"
	"time"
)

const (
	// ElectTime ...
	ElectTime = time.Millisecond * 150
	HbTime    = time.Second * 5
)

const (
	INIT int = iota
	RUNNING
	STOP
)

const (
	FOLLOWER int = iota
	CANDIDATE
	LEADER
)

// KVd 对外的
type KvsReq struct {
	k string
	v string
}

type KvsRsp struct {
	RetCode int32
	LogIdx  uint64
}

type Raft struct {
	CurrentTerm uint64
	VotedFor    uint32
	Log         []pb.Entries

	CommitIdx   uint64
	LastApplied uint64

	NextIdx  uint64
	MatchIdx uint64

	Peers []*PeerCli

	blockQ chan pb.Entries // receive KvsReq from KV service layer, leader
	applyR chan KvsRsp     // response KvsRsp to KV service layer, leader

	perstQ chan pb.Entries // send entries to persist layer, every node

	stopC chan int

	state int
	role  int

	metastor Stor // load meta
}

func NewRaft(ps []*PeerCli, stor Stor) *Raft {
	return &Raft{
		CurrentTerm: 0,
		VotedFor:    0,
		Log:         make([]pb.Entries, 1024),

		CommitIdx:   0,
		LastApplied: 0,

		NextIdx:  0,
		MatchIdx: 0,
		Peers:    ps,

		blockQ: make(chan pb.Entries, 256),
		applyR: make(chan KvsRsp, 256),
		stopC:  make(chan int),

		state:    INIT,
		metastor: stor,
	}
}

/*
	设置raft状态
	加载metadata: commit Idx, CurentTime, votedfor ?
	加载log,并apply到stat machine
	启动event loop，
		event loop中 启动election 逻辑
*/
func (r *Raft) Init() error {
	if r.state != INIT {
		log.Println("raft node already initiated!")
		return errors.New("raft already initiated!")
	}

	r.CurrentTerm, _ = r.metastor.LoadCurTerm()
	r.CommitIdx, _ = r.metastor.LoadCommitIdx()

	r.Log, _ = r.metastor.LoadData()
	// apply to stat machine

	r.state = RUNNING
	r.role = FOLLOWER

	go r.eventLoop()

	return nil
}

func (r *Raft) eventLoop() {
	if r.state != RUNNING {
		return
	}

	for {
		switch r.role {
		case FOLLOWER:
			r.followerLoop()
		case CANDIDATE:
			r.candidateLoop()
		case LEADER:
			r.leaderLoop()
		}
	}
}
