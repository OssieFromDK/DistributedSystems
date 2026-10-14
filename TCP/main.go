package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Packet2 struct {
    Type string
    Seq  uint32
    Ack  uint32
}

func client3(clientToServer, serverToClient chan Packet2, id int) {
    const timeout = 2 * time.Second
    const maxRetries = 5

    isn := rand.Uint32()

    for attempt := 1; attempt <= maxRetries; attempt++ {
        fmt.Printf("Client %d: Sending SYN seq=%d (attempt %d)\n", id, isn, attempt)
        clientToServer <- Packet2{Type: "SYN", Seq: isn, Ack: 0}

        select {
        case response := <-serverToClient:
            if response.Type == "SYN-ACK" && response.Ack == isn+1 {
                fmt.Printf("Client %d: Received SYN-ACK seq=%d ack=%d\n", id, response.Seq, response.Ack)
                clientToServer <- Packet2{Type: "ACK", Seq: isn + 1, Ack: response.Seq + 1}
                return
            }
        case <-time.After(timeout):
            fmt.Printf("Client %d: Timeout on SYN seq=%d, retrying\n", id, isn)
        }
    }
    fmt.Printf("Client %d: Gave up (ISN %d)\n", id, isn)
}

func server3(clientToServer, serverToClient chan Packet2) {
    const timeout = 2 * time.Second
    const maxRetries = 5

    packet := <-clientToServer
    if packet.Type != "SYN" { return }
    fmt.Printf("Server: Received SYN seq=%d\n", packet.Seq)

    isn := rand.Uint32()
    clientSeq := packet.Seq
    expectedAck := isn + 1

    for attempt := 1; attempt <= maxRetries; attempt++ {
        fmt.Printf("Server: Sending SYN-ACK seq=%d ack=%d (attempt %d)\n", isn, clientSeq+1, attempt)
        serverToClient <- Packet2{Type: "SYN-ACK", Seq: isn, Ack: clientSeq + 1}
        select {
        case reply := <-clientToServer:
            if reply.Type == "SYN" && reply.Seq == clientSeq {
                continue
            }
            if reply.Type == "ACK" && reply.Ack == expectedAck {
                fmt.Printf("Connection Established (server seq=%d, client seq=%d)\n", isn, reply.Seq)
                return
            }
        case <-time.After(timeout):
            fmt.Printf("Server: Timeout on SYN-ACK seq=%d, resending\n", isn)
        }
    }
}

func forwarder(clientIncoming chan Packet2, serverIncoming chan Packet2, serverOutgoing chan Packet2, clientOutgoing chan Packet2) {
    for {
        select {
        case packet := <-clientIncoming:
            if rand.Intn(100) < 20 {
                fmt.Printf("Forwarder (Client -> Server): LOST %s seq=%d ack=%d\n", packet.Type, packet.Seq, packet.Ack)
                continue
            }
            delay := time.Duration(rand.Intn(1000))
            time.Sleep(delay * time.Millisecond)
            fmt.Printf("Forwarder (Client -> Server): Forwarded %s seq=%d ack=%d after %v\n", packet.Type, packet.Seq, packet.Ack, delay*time.Millisecond)
            serverOutgoing <- packet

        case packet := <-serverIncoming:
            if rand.Intn(100) < 20 {
                fmt.Printf("Forwarder (Server -> Client): LOST %s seq=%d ack=%d\n", packet.Type, packet.Seq, packet.Ack)
                continue
            }
            delay := time.Duration(rand.Intn(1000))
            time.Sleep(delay * time.Millisecond)
            fmt.Printf("Forwarder (Server -> Client): Forwarded %s seq=%d ack=%d after %v\n", packet.Type, packet.Seq, packet.Ack, delay*time.Millisecond)
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
