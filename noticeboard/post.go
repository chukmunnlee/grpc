package main

import (
	"fmt"
	"log"
	"os"
	"context"

	uuid "github.com/google/uuid"
	grpc "google.golang.org/grpc"

	pb "github.com/chukmunnlee/noticeboard/messages"

	ptypes "github.com/golang/protobuf/ptypes"
)

const NOTICE_SERVICE string = "localhost:50051"

func main() {

	if len(os.Args) < 3 {
		log.Println("Missing arguments: sender message");
	}

	u, u_err := uuid.NewRandom()
	if nil != u_err {
		log.Fatalf("Cannot generate UUID: %v\n", u_err)
	}
	notice := pb.Notice {
		Id: u.String()[:8],
		From: os.Args[1],
		Note: os.Args[2],
	}
	req := pb.PostNoticeRequest {
		Time: ptypes.TimestampNow(),
		Note: &notice,
	}

	conn, err := grpc.Dial(NOTICE_SERVICE, grpc.WithInsecure())
	defer conn.Close()

	c := pb.NewNoticeBoardServiceClient(conn)
	resp, err := c.Post(context.Background(), &req)
	if nil != err {
		log.Fatalf("Cannot post message: %v\n", err)
	}
	fmt.Printf("Message posted: %v\n", resp)

}
