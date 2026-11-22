# Auction System

## How to run

Run nodes first:

```
go run ./node 1
go run ./node 2
go run ./node 3
```

Then clients:

```
go run ./client alice
go run ./client bob
```

Use `bid <amount>` from a client to place a bid and `result` to see the highest bid that is currently there, or the winner once the auction is over
The auction ends 100 seconds after node 1 starts, and you can crash any one node with Ctrl+C inside its terminal