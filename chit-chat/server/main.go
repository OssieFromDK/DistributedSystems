package main

import (
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "chitchat/grpc"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer

	mutex       sync.Mutex
	clients     map[string]chan *pb.ServerMessage
	logicalTime int64
}

func main() {
	startServer()
}

func startServer() {
	const port string = ":50051"

	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Could not start server: ", err)
	}

	var grpcServer *grpc.Server
	grpcServer = grpc.NewServer()

	var chatServer *ChatServer
	chatServer = &ChatServer{
		clients: make(map[string]chan *pb.ServerMessage),
	}

	pb.RegisterChatServiceServer(grpcServer, chatServer)

	log.Println("Chat server running on port 50051")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Server stopped unexpectedly: ", err)
	}
}

func (server *ChatServer) Connect(
	stream grpc.BidiStreamingServer[pb.ClientRequest, pb.ServerMessage],
) error {
	log.Println("Client connected")

	var outgoing chan *pb.ServerMessage
	outgoing = make(chan *pb.ServerMessage, 32)

	var clientID string
	var registered bool
	registered = false

	go func() {
		for message := range outgoing {
			var sendErr error
			sendErr = stream.Send(message)
			if sendErr != nil {
				log.Println("Could not send message to client:", sendErr)
				return
			}
		}
	}()

	defer func() {
		if registered {
			server.removeClient(clientID, outgoing)
		} else {
			close(outgoing)
		}
	}()

	for {
		var request *pb.ClientRequest
		var err error

		request, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed its stream")
			} else {
				log.Println("Client disconnected:", err)
			}
			return nil
		}

		log.Println(
			"Received request from client:",
			request.GetClientId(),
		)

		switch request.GetAction() {
		case pb.Action_JOIN:
			if registered {
				continue
			}

			clientID = request.GetClientId()
			if clientID == "" {
				outgoing <- &pb.ServerMessage{
					Content:   "Client ID cannot be empty.",
					EventType: "ERROR",
				}
				return nil
			}

			server.mutex.Lock()

			_, exists := server.clients[clientID]
			if exists {
				server.mutex.Unlock()

				outgoing <- &pb.ServerMessage{
					Content:   "That client ID is already in use.",
					EventType: "ERROR",
				}
				return nil
			}

			server.clients[clientID] = outgoing
			registered = true

			server.logicalTime++
			outgoing <- &pb.ServerMessage{
				ClientId:    "Server",
				Content:     "Welcome to Chit Chat!",
				LogicalTime: server.logicalTime,
				EventType:   "WELCOME",
			}

			server.logicalTime++
			var joinMessage *pb.ServerMessage
			joinMessage = &pb.ServerMessage{
				ClientId:    clientID,
				Content:     clientID + " joined the chat.",
				LogicalTime: server.logicalTime,
				EventType:   "JOIN",
			}

			server.broadcastLocked(joinMessage)
			server.mutex.Unlock()

			log.Println("Client joined:", clientID)

		case pb.Action_PUBLISH:
			if !registered {
				outgoing <- &pb.ServerMessage{
					Content:   "Join the chat before sending messages.",
					EventType: "ERROR",
				}
				continue
			}

			server.broadcast(
				"PUBLISH",
				clientID,
				request.GetContent(),
			)

		case pb.Action_LEAVE:
			if registered {
				server.removeClient(clientID, outgoing)
				registered = false
			}
			return nil

		default:
			log.Println("Received an unknown action from:", request.GetClientId())
		}
	}
}

func (server *ChatServer) broadcast(
	eventType string,
	clientID string,
	content string,
) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	server.logicalTime++

	var message *pb.ServerMessage
	message = &pb.ServerMessage{
		ClientId:    clientID,
		Content:     content,
		LogicalTime: server.logicalTime,
		EventType:   eventType,
	}

	server.broadcastLocked(message)
}

func (server *ChatServer) broadcastLocked(message *pb.ServerMessage) {
	for clientID, outgoing := range server.clients {
		select {
		case outgoing <- message:
		default:
			log.Println("Outgoing queue full for client:", clientID)
		}
	}
}

// removeClient removes a client
func (server *ChatServer) removeClient(
	clientID string,
	outgoing chan *pb.ServerMessage,
) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	currentChannel, exists := server.clients[clientID]
	if !exists || currentChannel != outgoing {
		close(outgoing)
		return
	}

	delete(server.clients, clientID)
	close(outgoing)

	server.logicalTime++

	var leaveMessage *pb.ServerMessage
	leaveMessage = &pb.ServerMessage{
		ClientId:    clientID,
		Content:     clientID + " left the chat.",
		LogicalTime: server.logicalTime,
		EventType:   "LEAVE",
	}

	server.broadcastLocked(leaveMessage)

	log.Println("Client left:", clientID)
}
