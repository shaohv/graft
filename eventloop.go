package graft

import (
	pb "graft/pb"
	"log"
	"math/rand"
	"time"
)

func randomTmout() <-chan time.Time {
	electTmout := rand.Intn(ElectTime) + ElectTime
	return time.After(time.Duration(electTmout) * time.Millisecond)
}

/*
	1. 设置选举超时时间
	2. 将blockQ上的数据保存到metastor，并应用到stat machine
	3.
*/
func (r *Raft) followerLoop() {
	var electTmout = randomTmout()
	var err error

	for r.state == FOLLOWER {
		update := false
		select {
		case <-r.stopC:
			log.Println("raft stoped")
			return
		case <-electTmout:
			// change state
			r.changeState(CANDIDATE)
		case task := <-r.blockQ:
			update = true
			switch req := task.args.(type) {
			case *pb.AppendEntriesReq:
				rsp, _ := (task.reply).(*pb.AppendEntriesRsp)
				*rsp, _ = r.handleAppendEntries(req)
			case *pb.VoteReq:
				rsp, _ := (task.reply).(*pb.VoteRsp)
				*rsp, _ = r.handleVote(req)
			}

			task.err <- err
		}

		if update {
			electTmout = randomTmout()
		}
	}
}

func (r *Raft) candidateLoop() {

}

func (r *Raft) leaderLoop() {

}
