package graft

import (
	pb "graft/pb"
)

type Stor interface {
	// PerstCurTerm 持久化CurrentTerm
	PerstCurTerm(term uint64) error

	// PerstVoteFor 持久化VoteFor
	PerstVoteFor(id uint32) error

	LoadCommitIdx() (uint64, error)

	LoadCurTerm() (uint64, error)

	LoadVoteFor() (uint32, error)

	// AppendEntries 追加Entry
	AppendEntries(es []pb.Entries) error

	// DelInvalidEntries
	DelInvalidEntries(idx uint64) error

	// LoadData 加载数据
	LoadData() ([]pb.Entries, error)

	GetEntryByIdx(idx uint64) (pb.Entries, error)
}
