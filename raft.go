package graft

import (
	"errors"
	"graft/pb"
	"log"
	"time"
)

const (
	// ElectTime ...
	ElectTime = 150
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

var (
	ErrNotLeader = errors.New("Raft server not leader!")
	ErrStop      = errors.New("Raft server stopped!")
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

type message struct {
	args  interface{}
	reply interface{}
	err   chan error
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

	blockQ chan *message //leader rcv KvsReq from KV service layer or follower rcv Req from leader
	applyR chan KvsRsp   // response KvsRsp to KV service layer, leader

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

		blockQ: make(chan *message, 256),
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

func (rf *Raft) deliver(value interface{}, r interface{}) error {
	if rf.state != RUNNING {
		return ErrStop
	}

	task := &message{args: value, reply: r, err: make(chan error, 1)}

	//deliver the task to event loop
	select {
	case rf.blockQ <- task:
	case <-rf.stopC:
		return ErrStop
	}
	// wait for the task been handle over, and return
	select {
	case <-rf.stopC:
		return ErrStop
	case suc := <-task.err:
		return suc
	}
}

func (rf *Raft) changeState(state int) {
	rf.state = state
	// 持久化
}
