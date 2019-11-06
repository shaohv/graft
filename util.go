package graft

import (
	"math/rand"
	"time"
)

const (
	RetOK    = 0
	RetError = -1
)

func randomTmout(t int) <-chan time.Time {
	electTmout := rand.Intn(t) + t
	return time.After(time.Duration(electTmout) * time.Millisecond)
}
