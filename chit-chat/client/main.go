package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "chitchat/grpc"
)

func main() {
	var clientID string
	clientID = "Client1"

	if len(os.Args) > 1 {
		clientID = os.Args[1]
	}

	var connection *grpc.ClientConn
	var err error

	connection, err = grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Could not connect to server: ", err)
	}
	defer connection.Close()

	log.Println("Connected to chat server")

	var client pb.ChatServiceClient
	client = pb.NewChatServiceClient(connection)

	var stream pb.ChatService_ConnectClient
	stream, err = client.Connect(context.Background())
	if err != nil {
		log.Fatal("Could not open chat stream: ", err)
	}

	log.Println("Chat stream opened")

	var joinRequest *pb.ClientRequest
	joinRequest = &pb.ClientRequest{
		ClientId: clientID,
		Action:   pb.Action_JOIN,
		Content:  "",
	}

	err = stream.Send(joinRequest)
	if err != nil {
		log.Fatal("Could not send join request: ", err)
	}

	log.Println("Joined chat as:", clientID)

	go func() {
		for {
			var response *pb.ServerMessage
			var receiveErr error

			response, receiveErr = stream.Recv()
			if receiveErr != nil {
				log.Println("Connection to server closed:", receiveErr)
				return
			}

			log.Printf(
				"[%d] %s",
				response.GetLogicalTime(),
				response.GetContent(),
			)
		}
	}()

	log.Println("Type a message and press Enter. Type /quit to leave.")

	var scanner *bufio.Scanner
	scanner = bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var content string
		content = scanner.Text()

		if content == "/quit" {
			var leaveRequest *pb.ClientRequest
			leaveRequest = &pb.ClientRequest{
				ClientId: clientID,
				Action:   pb.Action_LEAVE,
				Content:  "",
			}

			err = stream.Send(leaveRequest)
			if err != nil {
				log.Println("Could not send leave request:", err)
			}

			var closeErr error
			closeErr = stream.CloseSend()
			if closeErr != nil {
				log.Println("Could not close chat stream:", closeErr)
			}

			log.Println("Leaving chat")
			return
		}

		var messageRequest *pb.ClientRequest
		messageRequest = &pb.ClientRequest{
			ClientId: clientID,
			Action:   pb.Action_PUBLISH,
			Content:  content,
		}

		err = stream.Send(messageRequest)
		if err != nil {
			log.Println("Could not send message:", err)
			return
		}
	}

	if scanner.Err() != nil {
		log.Println("Error reading input:", scanner.Err())
	}

	var leaveRequest *pb.ClientRequest
	leaveRequest = &pb.ClientRequest{
		ClientId: clientID,
		Action:   pb.Action_LEAVE,
		Content:  "",
	}

	err = stream.Send(leaveRequest)
	if err != nil {
		log.Println("Could not send leave request:", err)
	}

	var closeErr error
	closeErr = stream.CloseSend()
	if closeErr != nil {
		log.Println("Could not close chat stream:", closeErr)
	}
}
