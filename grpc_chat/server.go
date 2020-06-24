package main

import (
	"context"
	"log"
	"net"

	"github.com/rudiarta/go_grpc/grpc_chat/chat"
	"google.golang.org/grpc"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *chat.HelloRequest) (*chat.HelloReply, error) {
	log.Printf("Received " + in.GetName())
	return &chat.HelloReply{Message: "Hello: " + in.GetName()}, nil
}

func main() {
	var c chat.GreeterServer = &Server{}
	listen, err := net.Listen("tcp", ":1006")
	if err != nil {
		log.Printf("Failed to Lister: ", err)
	}

	server := grpc.NewServer()
	chat.RegisterGreeterServer(server, c)
	if err := server.Serve(listen); err != nil {
		log.Printf("Fail to serve: ", err)
	}
}
