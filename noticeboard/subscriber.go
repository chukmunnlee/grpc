package main

import (
	"fmt"
	"log"
	"io"
	"context"

	pb "github.com/chukmunnlee/grpc/noticeboard/messages"

	grpc "google.golang.org/grpc"
)

const SERVICE string = "localhost:50051"

func main() {

	conn, err := grpc.Dial(SERVICE, grpc.WithInsecure())
	if nil != err {
		log.Fatalf("Cannot connect to NoticeBoardService %v\n", err)
	}
	defer conn.Close()

	c := pb.NewNoticeBoardServiceClient(conn)

	stream, s_err := c.Subscribe(context.Background(), &pb.SubscribeRequest{})
	if nil != s_err {
		log.Fatalf("Cannot subscribe to NoticeBoardService: %v\n", s_err)
	}

	for {
		resp, err := stream.Recv()
		if io.EOF == err {
			break
		}
		if nil != err {
			log.Fatalf("Error when receiving message: %\v\n", err)
		}
		notice := resp.GetNotice()
		fmt.Printf("[notice] id: %s, from: %s, note: %s\n", notice.GetId(), notice.GetFrom(), notice.GetNote())
	}

}
