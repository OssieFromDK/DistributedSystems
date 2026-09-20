package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Packet2 struct {
    Type   string
    Number int
}

func client3(clientToServer chan Packet2, serverToClient chan Packet2, counter int) {
    fmt.Printf("Client: Sending SYN %d\n", counter)
    time.Sleep(500 * time.Millisecond)

    clientToServer <- Packet2{
        Type:   "SYN",
        Number: counter,
    }

    response := <-serverToClient

    if response.Type == "SYN-ACK" && response.Number == counter+1 {
        fmt.Printf("Client: Received SYN-ACK %d\n", response.Number)

        clientToServer <- Packet2{
            Type:   "ACK",
            Number: counter + 1,
        }
    }
}

func server3(clientToServer chan Packet2, serverToClient chan Packet2) {
    packet := <-clientToServer

    if packet.Type == "SYN" {
        fmt.Printf("Server: Received SYN %d\n", packet.Number)
        time.Sleep(500 * time.Millisecond)

        serverToClient <- Packet2{
            Type:   "SYN-ACK",
            Number: packet.Number + 1,
        }
    }

    packet = <-clientToServer

    if packet.Type == "ACK" {
        fmt.Printf("Server: Received ACK %d\n", packet.Number)
        fmt.Printf("Connection %d Established With Success!\n", packet.Number)
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