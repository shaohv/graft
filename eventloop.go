package graft

import (
	"log"
)

/*
	1. 设置选举超时时间
	2. 将blockQ上的数据保存到metastor，并应用到stat machine
	3.
*/
func (r *Raft) followerLoop() {

	for {
		select {
		case <-r.stopC:
			log.Println("raft stoped")
			return
			//case <-time.After():

		}

	}
}

func (r *Raft) candidateLoop() {

}

func (r *Raft) leaderLoop() {

}
