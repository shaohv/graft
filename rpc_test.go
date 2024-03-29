package graft

import (
	"fmt"
	"graft/pb"
	"log"
	"sync"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func TestRpcVote(t *testing.T) {
	const server = "127.0.0.1:8788"
	var servC = make(chan *grpc.Server, 1)
	go StartGrpcServer(server, servC)
	s := <-servC

	var wg sync.WaitGroup
	wg.Add(1)

	go func(addr string) {
		log.Println("client conn ...")
		defer wg.Done()
		client, err := NewPeer(addr)
		if err != nil {
			log.Println("conn fail", err)
			return
		}
		defer client.DelPeer()

		req := &pb.VoteReq{}
		rsp, err := client.Cli.Vote(context.Background(), req)
		if err != nil {
			fmt.Println("rcv rsp fail", err)
			return
		}
		fmt.Printf("Term=%v, VoteGranted=%v", rsp.Term, rsp.VoteGranted)
	}(server)

	wg.Wait()
	s.Stop()
}

func TestRpcAppendEnties(t *testing.T) {
	const server = "127.0.0.1:8788"
	var servC = make(chan *grpc.Server, 1)
	go StartGrpcServer(server, servC)
	s := <-servC

	var wg sync.WaitGroup
	wg.Add(1)

	go func(addr string) {
		log.Println("client conn ...")
		defer wg.Done()
		client, err := NewPeer(addr)
		if err != nil {
			log.Println("conn fail", err)
			return
		}
		defer client.DelPeer()

		req := &pb.AppendEntriesReq{}
		rsp, err := client.Cli.AppendEntries(context.Background(), req)
		if err != nil {
			fmt.Println("rcv rsp fail", err)
			return
		}
		fmt.Printf("Term=%v, RetCode=%v", rsp.Term, rsp.RetCode)
	}(server)

	wg.Wait()
	s.Stop()
}
