package main

import (
	"fmt"
	"log"
	"time"
	"net"
	"strings"
	"context"
	pb "github.com/chukmunnlee/echo/messages"
	"google.golang.org/grpc"
	ptypes "github.com/golang/protobuf/ptypes"
)

type server struct{ }

func (*server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {

	echoMessage := req.GetData()

	log.Printf("Echo request: Id: %d\n", echoMessage.GetId())

	original := echoMessage.GetData()
	toUpper := strings.ToUpper(original)

	response := pb.EchoResponse {
		Time: ptypes.TimestampNow(),
		Code: 0,
		Status: "OK",
		Result: original + ": " + toUpper,
	}

	return &response, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051");
	if nil != err {
		log.Fatalf("Cannot listen on 50051: %s\n", err)
	}

	s := grpc.NewServer()

	pb.RegisterEchoServiceServer(s, &server{})


	fmt.Printf("Go gRPC ECHO Server: %s\n", time.Now())

	if err := s.Serve(lis); nil != err {
		log.Fatalf("Cannot start service: %v\n", err);
	}
}
