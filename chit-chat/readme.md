# Chit Chat

This is our chat program made in Go with grpc.

To run it, you first have to start the server with `go run ./server`. Then open more terminals and run the clients with `go run ./client navn`, `go run ./client navn1` etc.

You can send messages and they get sent to the other clients too. You can use `/quit` to leave the chat.

We use gRPC and protobuf for communication between the server and clients. The server also uses logical timestamps for the messages.

Remember to start the server before starting the clients.