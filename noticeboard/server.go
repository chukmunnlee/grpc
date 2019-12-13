package main

import (
	"fmt"
	"log"
	"net"
	"context"

	pb "github.com/chukmunnlee/grpc/noticeboard/messages"

	ptypes "github.com/golang/protobuf/ptypes"
	grpc "google.golang.org/grpc"
)

const TCP string = "tcp";
const PORT string = "0.0.0.0:50051";

type server struct {
	pb.UnimplementedNoticeBoardServiceServer
	channel chan pb.Notice
}

func (nbserv *server) Post(ctx context.Context, req *pb.PostNoticeRequest) (*pb.PostNoticeResponse, error) {
	note := req.GetNote()
	log.Printf("[Post] id: %s, from: %s, note: %s\n", note.GetId(), note.GetFrom(), note.GetNote()) 
	resp := pb.PostNoticeResponse {
		Id: note.GetId(),
		Code: 201,
		Status: "Accepted",
	}

	go func(c chan pb.Notice, r pb.Notice) {
		c <- r
	}(nbserv.channel, *note)

	return &resp, nil
}

func (nbserv *server) Subscribe(req *pb.SubscribeRequest, stream pb.NoticeBoardService_SubscribeServer) error {
	log.Println("New subscription")
	for notice := range nbserv.channel {
		log.Printf("\tSending %v\n", notice)
		resp :=  pb.SubscribeResponse{
			Time: ptypes.TimestampNow(),
			Notice: &notice,
		}
		if err := stream.Send(&resp); nil != err {
			log.Printf("Encounted error in Send. Closing stream: %v\n", err)
			return err
		}
	}
	return nil
}

func main() {

	s := grpc.NewServer()
	serv := server {
		channel: make(chan pb.Notice),
	}

	pb.RegisterNoticeBoardServiceServer(s, &serv)

	fmt.Println("Start NoticeBoardService")

	lis, err := net.Listen(TCP, PORT)
	if nil != err {
		log.Fatalf("Cannot listen on %s: &v\n", PORT, err)
	}
	if err := s.Serve(lis); nil != err {
		log.Fatalf("Cannot start NoticeBoardService: %v\n", err)
	}
}
