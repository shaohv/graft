package graft

import (
	pb "graft/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

//StartGrpcServer start grpc server, wait gprc client connect
func StartGrpcServer(addr string, servC chan *grpc.Server) {
	log.Println("Start grpc server now")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterPeerServer(s, &PeerServerImpl{})
	servC <- s
	s.Serve(lis)
}

type PeerCli struct {
	Conn *grpc.ClientConn
	Cli  pb.PeerClient
}

// GetPeer: setup connections with grpc server and return grpc client
func NewPeer(addr string) (*PeerCli, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("conn addr:%v fail:%v", addr, err)
		return nil, err
	}

	client := pb.NewPeerClient(conn)
	cli := &PeerCli{
		Conn: conn,
		Cli:  client,
	}
	return cli, nil
}

func (p *PeerCli) DelPeer() {
	p.Conn.Close()
}
