package graft

import (
	pb "graft/pb"
)

func (r *Raft) handleAppendEntries(req *pb.AppendEntriesReq) (pb.AppendEntriesRsp, error) {

	return pb.AppendEntriesRsp{}, nil
}

func (r *Raft) handleVote(req *pb.VoteReq) (pb.VoteRsp, error) {
	return pb.VoteRsp{}, nil
}
