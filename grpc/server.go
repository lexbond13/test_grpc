package grpc

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"strings"
	pb "test_repo/grpc/consumer"
)

type server struct {
	savedConsumers []*pb.ConsumerRequest
}

// CreateCustomer creates a new Customer
func (s *server) CreateConsumer(ctx context.Context, in *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	s.savedConsumers = append(s.savedConsumers, in)
	return &pb.ConsumerResponse{Id: in.Id, Success: true}, nil
}

func (s *server) GetConsumers(filter *pb.ConsumerFilter, stream pb.Consumer_GetConsumersServer) error {
	for _, customer := range s.savedConsumers {
		if filter.Keyword != "" {
			if !strings.Contains(customer.Name, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil
}

func Run(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	serv := grpc.NewServer()
	pb.RegisterConsumerServer(serv, &server{})
	serv.Serve(lis)
}
