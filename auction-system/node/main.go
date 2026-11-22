package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"maps"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	proto "auctionsystem/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type node struct {
	proto.UnimplementedAuctionServer
	peers  map[int32]proto.AuctionClient
	mu     sync.Mutex
	bidder string
	amount int64
	end    time.Time
	acks   map[string]*proto.Ack
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

	peers := make(map[int32]proto.AuctionClient)
	for i, address := range addresses {
		if i == id-1 {
			continue
		}
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		peers[int32(i+1)] = proto.NewAuctionClient(conn)
	}

	n := &node{peers: peers, end: time.Now().Add(100 * time.Second), acks: make(map[string]*proto.Ack)}
	server := grpc.NewServer()
	proto.RegisterAuctionServer(server, n)
	log.Printf("event=STARTUP address=%s end=%s", addresses[id-1], n.end.Format(time.TimeOnly))
	server.Serve(listener)
}

func (n *node) Bid(ctx context.Context, req *proto.BidRequest) (*proto.Ack, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	last, registered := n.acks[req.Bidder]
	if registered && last.Id == req.Id {
		log.Printf("event=DUPLICATE bidder=%s id=%d outcome=%s", req.Bidder, req.Id, last.Outcome)
		return last, nil
	}
	if !registered {
		log.Printf("event=REGISTER bidder=%s", req.Bidder)
	}

	ack := &proto.Ack{Id: req.Id, Outcome: "FAIL"}
	if time.Now().After(n.end) {
		ack.Outcome = "EXCEPTION"
	} else if req.Amount > n.amount {
		ack.Outcome = "SUCCESS"
		n.bidder, n.amount = req.Bidder, req.Amount
	}
	n.acks[req.Bidder] = ack
	log.Printf("event=BID bidder=%s id=%d amount=%d outcome=%s", req.Bidder, req.Id, req.Amount, ack.Outcome)
	n.replicate()
	return ack, nil
}

func (n *node) Result(ctx context.Context, req *proto.Empty) (*proto.Outcome, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	outcome := &proto.Outcome{Over: time.Now().After(n.end), Bidder: n.bidder, Amount: n.amount}
	log.Printf("event=RESULT over=%t bidder=%s amount=%d", outcome.Over, outcome.Bidder, outcome.Amount)
	n.replicate()
	return outcome, nil
}

func (n *node) Replicate(ctx context.Context, state *proto.State) (*proto.Empty, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.bidder, n.amount, n.end = state.Bidder, state.Amount, time.Unix(0, state.End)
	clear(n.acks)
	maps.Copy(n.acks, state.Acks)
	log.Printf("event=UPDATE bidder=%s amount=%d bidders=%d", n.bidder, n.amount, len(n.acks))
	return &proto.Empty{}, nil
}

func (n *node) replicate() {
	state := &proto.State{Bidder: n.bidder, Amount: n.amount, End: n.end.UnixNano(), Acks: n.acks}
	replies := 0
	for id, peer := range n.peers {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := peer.Replicate(ctx, state)
		cancel()
		if err != nil {
			log.Printf("event=REPLICA_DOWN node=%d error=%q", id, err)
			continue
		}
		replies++
	}
	log.Printf("event=REPLICATED replies=%d/%d", replies, len(n.peers))
}
