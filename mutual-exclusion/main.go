package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	proto "mutualexclusion/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	released = "RELEASED"
	wanted   = "WANTED"
	held     = "HELD"
)

type node struct {
	proto.UnimplementedRicartAgrawalaServer
	id             int32
	peers          map[int32]proto.RicartAgrawalaClient
	mu             sync.Mutex
	state          string
	clock          int64
	requestLamport int64
	deferred       []chan struct{}
}

func main() {
	addresses := []string{"localhost:5001", "localhost:5002", "localhost:5003"}
	id, _ := strconv.Atoi(os.Args[1])

	listener, err := net.Listen("tcp", addresses[id-1])
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll("logs", 0755)
	file, err := os.Create(fmt.Sprintf("logs/node%d.log", id))
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(fmt.Sprintf("[Node %d] ", id))

	peers := make(map[int32]proto.RicartAgrawalaClient)
	for i, address := range addresses {
		if i == id-1 {
			continue
		}
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.Config{BaseDelay: time.Second, Multiplier: 1, MaxDelay: time.Second}, MinConnectTimeout: 5 * time.Second}))
		if err != nil {
			log.Fatal(err)
		}
		peers[int32(i+1)] = proto.NewRicartAgrawalaClient(conn)
	}

	n := &node{id: int32(id), peers: peers, state: released}
	server := grpc.NewServer()
	proto.RegisterRicartAgrawalaServer(server, n)
	go server.Serve(listener)
	log.Printf("event=STARTUP address=%s peers=%d", addresses[id-1], len(peers))

	for {
		time.Sleep(time.Second + rand.N(3*time.Second))
		n.enter()
		time.Sleep(time.Second)
		n.exit()
	}
}

func (n *node) tick(t int64) int64 {
	n.clock = max(n.clock, t) + 1
	return n.clock
}

func (n *node) enter() {
	n.mu.Lock()
	n.state = wanted
	lamport := n.tick(0)
	n.requestLamport = lamport
	log.Printf("event=REQUEST_CS lamport=%d", lamport)
	for peer := range n.peers {
		log.Printf("event=SEND_REQUEST to=%d lamport=%d", peer, lamport)
	}
	n.mu.Unlock()

	var wg sync.WaitGroup
	replies := 0
	for peer, client := range n.peers {
		wg.Go(func() {
			for {
				reply, err := client.Request(context.Background(), &proto.Request{Id: n.id, Lamport: lamport}, grpc.WaitForReady(true))
				if err != nil {
					log.Printf("event=RETRY to=%d error=%q", peer, err)
					time.Sleep(time.Second)
					continue
				}
				n.mu.Lock()
				local := n.tick(reply.Lamport)
				replies++
				log.Printf("event=RECEIVE_REPLY from=%d lamport=%d local=%d replies=%d/%d", reply.Id, reply.Lamport, local, replies, len(n.peers))
				n.mu.Unlock()
				return
			}
		})
	}
	wg.Wait()

	n.mu.Lock()
	n.state = held
	log.Printf("event=ENTER_CS lamport=%d", n.clock)
	n.mu.Unlock()
}

func (n *node) exit() {
	n.mu.Lock()
	n.state = released
	log.Printf("event=EXIT_CS lamport=%d deferred=%d", n.clock, len(n.deferred))
	for _, release := range n.deferred {
		close(release)
	}
	n.deferred = nil
	n.mu.Unlock()
}

func (n *node) Request(ctx context.Context, req *proto.Request) (*proto.Reply, error) {
	n.mu.Lock()
	local := n.tick(req.Lamport)
	log.Printf("event=RECEIVE_REQUEST from=%d lamport=%d local=%d state=%s", req.Id, req.Lamport, local, n.state)
	if n.state == held || (n.state == wanted && (n.requestLamport < req.Lamport || (n.requestLamport == req.Lamport && n.id < req.Id))) {
		release := make(chan struct{})
		n.deferred = append(n.deferred, release)
		log.Printf("event=DEFER_REPLY to=%d", req.Id)
		n.mu.Unlock()
		select {
		case <-release:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		n.mu.Lock()
	}

	lamport := n.tick(0)
	log.Printf("event=SEND_REPLY to=%d lamport=%d", req.Id, lamport)
	n.mu.Unlock()
	return &proto.Reply{Id: n.id, Lamport: lamport}, nil
}
