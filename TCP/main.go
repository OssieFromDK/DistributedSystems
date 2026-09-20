package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Packet2 struct {
    Type   string
    Number uint32
}

func client3(clientToServer, serverToClient chan Packet2, id int) {
    const timeout = 2 * time.Second
    const maxRetries = 5

    isn := rand.Uint32()

    for attempt := 1; attempt <= maxRetries; attempt++ {
        fmt.Printf("Client %d: Sending SYN %d (attempt %d)\n", id, isn, attempt)
        clientToServer <- Packet2{Type: "SYN", Number: isn}

        select {
        case response := <-serverToClient:
            if response.Type == "SYN-ACK" && response.Number == isn+1 {
                fmt.Printf("Client %d: Received SYN-ACK %d\n", id, response.Number)
                clientToServer <- Packet2{Type: "ACK", Number: isn + 1}
                return
            }
        case <-time.After(timeout):
            fmt.Printf("Client %d: Timeout on SYN %d, retrying\n", id, isn)
        }
    }
    fmt.Printf("Client %d: Gave up (ISN %d)\n", id, isn)
}

func server3(clientToServer, serverToClient chan Packet2) {
    const timeout = 2 * time.Second
    const maxRetries = 5

    packet := <-clientToServer
    if packet.Type != "SYN" { return }
    fmt.Printf("Server: Received SYN %d\n", packet.Number)
    expected := packet.Number + 1

    for attempt := 1; attempt <= maxRetries; attempt++ {
        serverToClient <- Packet2{Type: "SYN-ACK", Number: expected}
        select {
        case ack := <-clientToServer:
            if ack.Type == "SYN" && ack.Number == packet.Number {
                continue
            }
            if ack.Type == "ACK" && ack.Number == expected {
                fmt.Printf("Connection %d Established\n", ack.Number)
                return
            }
        case <-time.After(timeout):
            fmt.Printf("Server: Timeout on SYN-ACK %d, resending\n", expected)
        }
    }
}

func forwarder(clientIncoming chan Packet2, serverIncoming chan Packet2, serverOutgoing chan Packet2, clientOutgoing chan Packet2) {
    for {
        select {
        case packet := <-clientIncoming:
            if rand.Intn(100) < 20 {
                fmt.Printf("Forwarder (Client -> Server): LOST %s %d\n", packet.Type, packet.Number)
                continue
            }
            delay := time.Duration(rand.Intn(1000))
            time.Sleep(delay * time.Millisecond)
            fmt.Printf("Forwarder (Client -> Server): Forwarded %s %d after %v\n", packet.Type, packet.Number, delay*time.Millisecond)
            serverOutgoing <- packet

        case packet := <-serverIncoming:
            if rand.Intn(100) < 20 {
                fmt.Printf("Forwarder (Server -> Client): LOST %s %d\n", packet.Type, packet.Number)
                continue
            }
            delay := time.Duration(rand.Intn(1000))
            time.Sleep(delay * time.Millisecond)
            fmt.Printf("Forwarder (Server -> Client): Forwarded %s %d after %v\n", packet.Type, packet.Number, delay*time.Millisecond)
            clientOutgoing <- packet
        }
    }
}

func main() {

    for i := 1; i <= 5; i++ {
        clientToForwarder := make(chan Packet2)
        forwarderToServer := make(chan Packet2)

        serverToForwarder := make(chan Packet2)
        forwarderToClient := make(chan Packet2)

        go forwarder(clientToForwarder, serverToForwarder, forwarderToServer, forwarderToClient)

        go client3(clientToForwarder, forwarderToClient, i)
        go server3(forwarderToServer, serverToForwarder)
    }

    time.Sleep(5 * time.Second)
}