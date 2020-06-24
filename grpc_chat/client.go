package main

import (
	"context"
	"log"
	"time"

	"github.com/rudiarta/go_grpc/grpc_chat/chat"
	"google.golang.org/grpc"
)

func main() {
	con, err := grpc.Dial("127.0.0.1:1006", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Did not connect: ", err)
	}
	defer con.Close()

	client := chat.NewGreeterClient(con)

	name := "testUser"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.SayHello(ctx, &chat.HelloRequest{Name: name})
	if err != nil {
		log.Fatal("Couldn't Greet: ", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
