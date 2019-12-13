package main

import (
	"fmt"
	"log"
	"net"
	"context"

	pb "github.com/chukmunnlee/calculator/messages"

	"google.golang.org/grpc"
	ptypes "github.com/golang/protobuf/ptypes"
)

const TCP string = "tcp"
const HOST string = "0.0.0.0:50051"

type server struct { } 

func (*server) Add(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {

	reqId := req.GetId()
	calc := req.GetCalculation()
	op0 := calc.GetOperand0()
	op1 := calc.GetOperand1()
	result := op0 + op1

	log.Printf("<< Add - id: %d, op0: %d, op1: %d result: %d\n", reqId, op0, op1, result);

	response := pb.CalculationResponse{
		Id: reqId,
		Time: ptypes.TimestampNow(),
		Result: result,
		Status: "OK",
		Code: 201,
	}

	return &response, nil
}


func main() {

	lis, err := net.Listen(TCP, HOST)
	if nil != err {
		log.Fatalf("Cannot listen to port %d: %v\n", HOST, err)
	}

	s := grpc.NewServer()

	pb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Starting CalculatorService")
	if err := s.Serve(lis); nil != err {
		log.Fatalf("Cannot start CalculatorService: %v\n", err)
	}
}
