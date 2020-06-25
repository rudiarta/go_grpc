package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/rudiarta/go_grpc/grpc_chat_stream/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type Connection struct {
	stream chat.BroadCast_CreateStreamServer
	id     string
	active bool
	err    chan error
}

type Server struct {
	Connection []*Connection
}

func (s *Server) CreateStream(pconn *chat.Connect, stream chat.BroadCast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		err:    make(chan error),
	}
	s.Connection = append(s.Connection, conn)

	return <-conn.err
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *chat.Message) (*chat.Close, error) {
	wait := sync.WaitGroup{}

	for _, conn := range s.Connection {
		// log.Println(conn.id)
		wait.Add(1)

		go func(msg *chat.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				log.Println("Sending message " + msg.Id + " to " + conn.id)
				err := conn.stream.Send(msg)

				if err != nil {
					grpclog.Errorf("Error with stream %v to %v", conn.stream, err)
					conn.active = false
					conn.err <- err
				}
			}

		}(msg, conn)
	}

	wait.Wait()

	return &chat.Close{}, nil
}

func main() {
	var connection []*Connection
	var server chat.BroadCastServer

	server = &Server{connection}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":1006")
	if err != nil {
		log.Printf("Error createing the server %v", err)
	}

	grpclog.Info("Starting server at port :1006")
	chat.RegisterBroadCastServer(grpcServer, server)
	grpcServer.Serve(listener)
}
