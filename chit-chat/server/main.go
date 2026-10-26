package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	proto "chitchat/grpc"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedChitChatServer
	actions chan func()
	clients map[string]chan *proto.Message
	lamport int64
}

func main() {
	os.MkdirAll("logs", 0755)
	file, err := os.Create("logs/server.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetPrefix("[Server] ")

	s := &server{actions: make(chan func()), clients: make(map[string]chan *proto.Message)}

	go func() {
		for action := range s.actions {
			action()
		}
	}()

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterChitChatServer(grpcServer, s)

	log.Println("event=STARTUP address=:5050")
	go grpcServer.Serve(listener)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("event=SHUTDOWN")
}

func (s *server) Join(req *proto.Message, stream proto.ChitChat_JoinServer) error {
	log.Printf("event=CONNECT client=%s lamport=%d", req.Sender, req.Lamport)
	out := make(chan *proto.Message, 100)
	s.actions <- func() {
		s.clients[req.Sender] = out
		s.lamport = max(s.lamport, req.Lamport) + 1
		s.broadcast("Server", fmt.Sprintf("Participant %s joined Chit Chat at logical time %d", req.Sender, s.lamport))
	}

	for msg := range out {
		stream.Send(msg)
		log.Printf("event=DELIVER client=%s lamport=%d", req.Sender, msg.Lamport)
	}
	return nil
}

func (s *server) Publish(ctx context.Context, msg *proto.Message) (*proto.Empty, error) {
	s.actions <- func() {
		s.lamport = max(s.lamport, msg.Lamport) + 1
		s.broadcast(msg.Sender, msg.Text)
	}
	return &proto.Empty{}, nil
}

func (s *server) Leave(ctx context.Context, msg *proto.Message) (*proto.Empty, error) {
	log.Printf("event=DISCONNECT client=%s lamport=%d", msg.Sender, msg.Lamport)
	s.actions <- func() {
		close(s.clients[msg.Sender])
		delete(s.clients, msg.Sender)
		s.lamport = max(s.lamport, msg.Lamport) + 1
		s.broadcast("Server", fmt.Sprintf("Participant %s left Chit Chat at logical time %d", msg.Sender, s.lamport))
	}
	return &proto.Empty{}, nil
}

func (s *server) broadcast(sender, text string) {
	log.Printf("event=BROADCAST sender=%s lamport=%d text=%q", sender, s.lamport, text)
	for _, out := range s.clients {
		out <- &proto.Message{Sender: sender, Text: text, Lamport: s.lamport}
	}
}
