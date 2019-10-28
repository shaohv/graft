package graft

import (
	pb "graft/pb"
	"log"
)

type MemStor struct {
	persTerm uint64

	persVoteFor uint32

	log []pb.Entries
}

func NewMemStor(t uint64, v uint32) *MemStor {
	return &MemStor{
		persTerm:    t,
		persVoteFor: v,
	}
}
func (s *MemStor) PerstCurTerm(term uint64) error {
	log.Printf("OldTerm:%v, CurTerm:%v", s.persTerm, term)
	s.persTerm = term
	return nil
}

// PerstVoteFor 持久化VoteFor
func (s *MemStor) PerstVoteFor(id uint32) error {
	log.Printf("OldVoteFor:%v, CurVoteFor:%v", s.persVoteFor, id)
	s.persVoteFor = id
	return nil
}

func (s *MemStor) LoadCurTerm() (uint64, error) {
	return s.persTerm, nil
}

func (s *MemStor) LoadVoteFor() (uint32, error) {
	return s.persVoteFor, nil
}

func (s *MemStor) LoadCommitIdx() (uint64, error) {
	return 0, nil
}

// AppendEntries 追加Entry
func (s *MemStor) AppendEntries(es []pb.Entries) error {
	s.log = append(s.log, es...)
	return nil
}

// DelInvalidEntries
func (s *MemStor) DelInvalidEntries(es []pb.Entries) error {
	return nil
}

// LoadData 加载数据
func (s *MemStor) LoadData() ([]pb.Entries, error) {
	return s.log, nil
}
