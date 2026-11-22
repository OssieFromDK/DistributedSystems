package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	proto "auctionsystem/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	name := os.Args[1]
	addresses := []string{"localhost:5001", "localhost:5002", "localhost:5003"}
	var nodes []proto.AuctionClient
	for _, address := range addresses {
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		nodes = append(nodes, proto.NewAuctionClient(conn))
	}

	os.MkdirAll("logs", 0755)
	file, err := os.Create("logs/" + name + ".log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(fmt.Sprintf("[Client %s] ", name))

	current := 0
	call := func(rpc func(ctx context.Context, node proto.AuctionClient) error) {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := rpc(ctx, nodes[current])
			cancel()
			if err == nil {
				return
			}
			log.Printf("event=NODE_DOWN node=%d error=%q", current+1, err)
			current = (current + 1) % len(nodes)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[0] == "bid" {
			amount, _ := strconv.ParseInt(fields[1], 10, 64)
			id := time.Now().UnixNano()
			var ack *proto.Ack
			call(func(ctx context.Context, node proto.AuctionClient) (err error) {
				ack, err = node.Bid(ctx, &proto.BidRequest{Bidder: name, Id: id, Amount: amount})
				return err
			})
			fmt.Println(ack.Outcome)
			log.Printf("event=BID node=%d id=%d amount=%d outcome=%s", current+1, id, amount, ack.Outcome)
		} else if len(fields) == 1 && fields[0] == "result" {
			var outcome *proto.Outcome
			call(func(ctx context.Context, node proto.AuctionClient) (err error) {
				outcome, err = node.Result(ctx, &proto.Empty{})
				return err
			})
			if outcome.Over {
				fmt.Printf("auction is over, winner is %s with %d\n", outcome.Bidder, outcome.Amount)
			} else {
				fmt.Printf("highest bid is %d by %s\n", outcome.Amount, outcome.Bidder)
			}
			log.Printf("event=RESULT node=%d over=%t bidder=%s amount=%d", current+1, outcome.Over, outcome.Bidder, outcome.Amount)
		} else {
			fmt.Println("commands: bid <amount>, result")
		}
	}
}
