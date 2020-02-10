package grpc

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
)

func CreateConsumer(client pb.ConsumerClient, consumer *pb.ConsumerRequest) {
	resp, err := client.CreateConsumer(context.Background(), consumer)
	if err != nil {
		log.Fatalf("Could not create consumer: %v", err)
	}
	if resp.Success {
		log.Printf("A new consumer has been added with id: %d", resp.Id)
	}
}

func GetConsumers(client pb.ConsumerClient, filter *pb.ConsumerFilter) {
	// calling the streaming API
	stream, err := client.GetConsumers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		consumer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", consumer)
	}
}

func NewClient(address string) *pb.ConsumerClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewConsumerClient(conn)
	return &client
}
