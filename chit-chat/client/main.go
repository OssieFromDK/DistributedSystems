package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	proto "chitchat/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	name := os.Args[1]
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := proto.NewChitChatClient(conn)

	lamport := make(chan int64, 1)
	lamport <- 0
	tick := func(t int64) int64 {
		next := max(<-lamport, t) + 1
		lamport <- next
		return next
	}

	stream, err := client.Join(context.Background(), &proto.Message{Sender: name, Lamport: tick(0)})
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll("logs", 0755)
	file, err := os.Create("logs/" + name + ".log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetPrefix("[Client] ")
	log.Printf("event=CONNECT client=%s", name)

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}
			local := tick(msg.Lamport)
			fmt.Printf("[%d] %s: %s\n", msg.Lamport, msg.Sender, msg.Text)
			log.Printf("event=RECEIVE client=%s sender=%s lamport=%d local=%d text=%q", name, msg.Sender, msg.Lamport, local, msg.Text)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && scanner.Text() != "/quit" {
		if utf8.RuneCountInString(scanner.Text()) > 128 {
			fmt.Println("message is longer than 128 characters")
			continue
		}

		t := tick(0)
		log.Printf("event=PUBLISH client=%s lamport=%d text=%q", name, t, scanner.Text())
		client.Publish(context.Background(), &proto.Message{Sender: name, Text: scanner.Text(), Lamport: t})
	}

	log.Printf("event=DISCONNECT client=%s", name)
	client.Leave(context.Background(), &proto.Message{Sender: name, Lamport: tick(0)})
	conn.Close()
}
