package graft

import (
	"golang.org/x/net/context"

	pb "graft/pb"
)

//PeerServerImpl ...
type PeerServerImpl struct{}

func (p *PeerServerImpl) AppendEntries(ctx context.Context, req *pb.AppendEntriesReq) (*pb.AppendEntriesRsp, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
	rsp := &pb.AppendEntriesRsp{
		Term:    1,
		RetCode: 0,
	}

	return rsp, nil
}
func (p *PeerServerImpl) Vote(ctx context.Context, req *pb.VoteReq) (*pb.VoteRsp, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
	rsp := &pb.VoteRsp{
		Term:        2,
		VoteGranted: 1,
	}

	return rsp, nil
}
