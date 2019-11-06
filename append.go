package graft

import (
	"errors"
	pb "graft/pb"
	"log"
)

func (r *Raft) handleAppendEntries(req *pb.AppendEntriesReq) (pb.AppendEntriesRsp, error) {
	//var ret = RetOK
	if req.GetTerm() < r.CurrentTerm {
		log.Println("term too small")
		return pb.AppendEntriesRsp{Term: r.CurrentTerm, RetCode: RetError},
			errors.New("term too small")
	}

	entry, err := r.metastor.GetEntryByIdx(req.GetPrevLogIdx())
	if err != nil {
		log.Println("get log failed")
		return pb.AppendEntriesRsp{Term: r.CurrentTerm, RetCode: RetError}, err
	}

	if entry.GetTerm() != req.GetTerm() {
		r.metastor.DelInvalidEntries(req.GetPrevLogIdx()) //删除该条之后的log
		log.Println("term not matching")
		return pb.AppendEntriesRsp{Term: r.CurrentTerm, RetCode: RetError}, err
	}

	//r.metastor.AppendEntries(req.GetEntries())
	commitIdx := req.GetLeaderCommitLogIdx()
	for _, e := range req.GetEntries() {
		//r.perstQ <- *e
		if commitIdx > e.GetLogIdx() {
			commitIdx = e.GetLogIdx()
		}
		r.metastor.AppendEntries([]pb.Entries{*e})
	}

	r.CommitIdx = commitIdx
	//r.metastor.PerstCurTerm()
	//r.metastor.PerstVoteFor()

	return pb.AppendEntriesRsp{Term: r.CurrentTerm, RetCode: RetError}, nil
}

func (r *Raft) handleVote(req *pb.VoteReq) (pb.VoteRsp, error) {

	if req.GetTerm() < r.CurrentTerm {
		log.Println("term too small")
		return pb.VoteRsp{Term: r.CurrentTerm, VoteGranted: RetError}, errors.New("term too small")
	}

	return pb.VoteRsp{}, nil
}
