package main

import (
	"fmt"
	"log"
	"time"
	"math/rand"

	"context"

	pb "github.com/chukmunnlee/grpc/calculator/messages"

	grpc "google.golang.org/grpc"
	ptypes "github.com/golang/protobuf/ptypes"
)

const SERVICE string = "localhost:50051"

func main() {

	conn, err := grpc.Dial(SERVICE, grpc.WithInsecure())
	defer conn.Close()

	if nil != err {
		log.Fatalf("Cannot connect to service: %v\n", err)
	}

	client := pb.NewCalculatorServiceClient(conn);

	fmt.Println("Connected to service: ", client)

	rand.Seed(time.Now().UnixNano())

	id := rand.Uint32()
	op0 := int32(rand.Int())
	op1 := int32(rand.Int())

	cal := pb.Calculation {
		Operand0: op0,
		Operand1: op1,
	}

	request := pb.CalculationRequest {
		Id: id,
		Time: ptypes.TimestampNow(),
		Calculation: &cal,
	}

	fmt.Printf(">> Add - Id: %d, op0: %d, op1: %d\n", id, op0, op1)

	resp, err := client.Add(context.Background(), &request);
	if nil != err {
		log.Fatalf("Invocation error: %v\n", err)
	}

	result := resp.GetResult();
	fmt.Printf("<< Add - Id: %d, result: %d\n", id, result)

}
