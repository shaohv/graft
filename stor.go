package graft

import (
	pb "graft/pb"
)

type Stor interface {
	AppendEntries(es []pb.Entries) error
	DelInvalidEntries()
}
