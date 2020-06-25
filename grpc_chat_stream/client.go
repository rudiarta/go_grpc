package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rudiarta/go_grpc/grpc_chat_stream/chat"
	"google.golang.org/grpc"
)

var client chat.BroadCastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func Connect(user *chat.User) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &chat.Connect{
		User:   user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	wait.Add(1)
	go func(str chat.BroadCast_CreateStreamClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}

			fmt.Printf("%v : %s \n", msg.Id, msg.Message)
		}
	}(stream)

	return streamerror
}

func main() {
	timestamp := time.Now()
	done := make(chan int)

	name := flag.String("N", "Anon", "The Name of the user")
	flag.Parse()

	id := sha256.Sum256([]byte(timestamp.String() + *name))
	conn, err := grpc.Dial("127.0.0.1:1006", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to service: %v", err)
	}

	client = chat.NewBroadCastClient(conn)
	user := &chat.User{
		Id:          hex.EncodeToString(id[:]),
		DisplayName: *name,
	}

	Connect(user)

	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &chat.Message{
				Id:        user.Id,
				Message:   scanner.Text(),
				Timestamp: timestamp.String(),
			}

			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error Sending Message: %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
