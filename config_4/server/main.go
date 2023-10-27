package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "go_config/proto"
	"log"
	"net"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedMyServiceServer
}
type KeyValue struct {
	ID    string        `json:"_id"`
	Key   string        `json:"key"`
	Value bson.RawValue `json:"value"`
}

func (s *server) InsertData(ctx context.Context, req *pb.Request) (*emptypb.Empty, error) {

	var value interface{}
	err := json.Unmarshal([]byte(req.Value), &value)
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("kishore").Collection("nithish")

	document := bson.M{
		"Key":   req.Key,
		"Value": value,
	}

	result, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		fmt.Println("error in inserting in db")
		return nil, err
	}

	fmt.Println(result)

	return &emptypb.Empty{}, nil
}

func (s *server) GetData(ctx context.Context, req *pb.GetDataRequest) (*emptypb.Empty, error) {
	
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	
	collection := client.Database("kishore").Collection("nithish")

	
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var data []KeyValue
	for cursor.Next(context.TODO()) {
		var kv KeyValue
		if err := cursor.Decode(&kv); err != nil {
			log.Fatal(err)
		}
		data = append(data, kv)
	}

	
	pattern := regexp.QuoteMeta(req.Prefix)
	result := make(map[string]map[string]interface{})
	for _, item := range data {
		if regex := regexp.MustCompile(pattern); regex.MatchString(item.Key) {
			subkey := item.Key[len(req.Prefix):]

			if result[subkey] == nil {
				result[subkey] = make(map[string]interface{})
			}

			
			var value interface{}
			if err := item.Value.Unmarshal(&value); err != nil {
				fmt.Println("Error parsing value:", err)
				return nil, nil
			}

			result[subkey]["value"] = value
		}
	}

	
	fmt.Println("Output:")
	for subkey, values := range result {
		fmt.Printf("%s{\n", subkey)
		for key, v := range values {
			fmt.Printf("  %s: %v\n", key, v)
		}
		fmt.Println("}")
	}
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Listening")
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterMyServiceServer(s, &server{})
	if err2 := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err2)
	}
}
