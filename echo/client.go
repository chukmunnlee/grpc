package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"context"
	"math/rand"

	"google.golang.org/grpc"

	pb "github.com/chukmunnlee/echo/messages"
	ptypes "github.com/golang/protobuf/ptypes"
)

func main() {

	message := "Standard greeter"
	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	rand.Seed(time.Now().UnixNano())

	echoMessage := pb.EchoMessage {
		Time: ptypes.TimestampNow(),
		Id: rand.Uint32(),
		Data: message,
	}

	request := pb.EchoRequest {
		Time: echoMessage.Time,
		Data: &echoMessage,
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()

	if nil != err {
		log.Fatalf("Cannot connect to EchoService on port 50051: %s\n", err)
	}

	c := pb.NewEchoServiceClient(conn)
	fmt.Printf("Connected to EchoService:")

	response, err := c.Echo(context.Background(), &request)
	if nil != err {
		log.Fatalf("Echo invocation error: %v\n", err);
	}

	fmt.Printf("Id: %d, Response: code: %d, status: %s, result: %s\n", 
			echoMessage.GetId(),
			response.GetCode(), response.GetStatus(), response.GetResult())
}

