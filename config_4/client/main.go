package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "go_config/proto"
	"log"
	"os"

	"google.golang.org/grpc"
)

var client pb.MyServiceClient

func main() {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewMyServiceClient(conn)

	// InsertData()

	GetData()
}

func InsertData() {
	key := "my.client.number"
	value := 9894364

	valueJSON, err := json.Marshal(value)
	fmt.Println("bytes of value: ", valueJSON)
	if err != nil {
		log.Fatalf("Failed to marshal value to JSON: %v", err)
	}

	req := pb.Request{
		Key:   key,
		Value: string(valueJSON),
	}

	_, err = client.InsertData(context.Background(), &req)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}
}

func GetData() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <prefix>")
		return
	}

	prefix := os.Args[1]
	req := pb.GetDataRequest{
		Prefix: prefix,
	}

	_, err := client.GetData(context.Background(), &req)
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
}
