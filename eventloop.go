package graft

import (
	pb "graft/pb"
	"log"
	"math/rand"
	"time"
)

func randomTmout(t int) <-chan time.Time {
	electTmout := rand.Intn(t) + t
	return time.After(time.Duration(electTmout) * time.Millisecond)
}

func (r *Raft) handleTask(task *message) {
	var err error
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

/*
	1. 设置选举超时时间
	2. 将blockQ上的数据保存到metastor，并应用到stat machine
	3.
*/
func (r *Raft) followerLoop() {
	var electTmout = randomTmout(ElectTime)
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
			r.handleTask(task)
			update = true
		}

		if update {
			electTmout = randomTmout(HbTime)
		}
	}
}

func (r *Raft) requestVotes() {

}

/*
	1. 发送选举请求
	2. 处理blockQ消息
*/
func (r *Raft) candidateLoop() {
	//var electTmout = randomTmout(ElectTime)
	for r.state == CANDIDATE {

		// 发送选举请求
		r.requestVotes()

		select {
		case <-r.stopC:
			log.Println("raft stopped!")
			return
		//case <-electTmout:
		case task := <-r.blockQ:
			r.handleTask(task)
		}

	}
}

// appendLog 心跳消息中将数据带给follower
func (r *Raft) appendLog() {

}

/*
	0. 从candidates -> leader后，第一次同步特殊处理
	1. 处理blockQ
	2. hb同步
	3. 处理blockQ
*/
func (r *Raft) leaderLoop() {
	var hbTmout = randomTmout(HbTime)

	// first sync

	for r.state == LEADER {
		select {
		case <-r.stopC:
			log.Println("raft stopped!")
			return
		case <-hbTmout:
			r.appendLog()
			hbTmout = randomTmout(HbTime)
		case task := <-r.blockQ:
			r.handleTask(task)
		}
	}
}
