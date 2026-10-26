# Report

## Discuss, whether you are going to use server-side streaming, client-side streaming, or bidirectional streaming?

We use server-side streaming because when a client joins, it makes a Join call once, and the server then keeps that stream open for as long as the client is connected and all broadcasts are sent to the client through that stream. As an example client A calls Publish with "Hello", and the server sends that message through the stream of every active client, including A's own so everyone recieves it.

## Describe your system architecture - do you have a server-client architecture, peer-to-peer, or something else?

We use a client-server architecture where the server is one process which listens on port 5050, then it keeps the list of active clients and the lamport clock on the server. All new chat joins, messages and leaves go through the server and it broadcasts it to everyone else in the chat.
Each client doesn't know the other clients, it only talks with the server and through it to the other clients.
Its not peer-to-peer because each client don't know each other and they can't reach each other without the server being in the middle.

## Describe what RPC methods are implemented, of what type, and what  messages types are used for communication

We have three different RPC methods implemented:
Join which is a server-side streaming type, and the message type is Message{sender, lamport}, this adds the client to the list of clients on the server, broadcasts that "Participant X joined Chit Chat at logical time L" to everyone including itself, and then keeps the stream open and sends every last broadcast through it
Publish which is a Unary RPC call where the message type looks like Message{sender, text, lamport}, this broadcasts the message with a new lamport time to all active clients including itself
Leave which is also a Unary RPC call where the message type looks like Message{sender, lamport}, this removes the client which ends its Join stream and then broadcasts to all other clients that "Participant X left Chit Chat at logical time L"

## Describe how you have implemented the calculation of the timestamps

We use Lamport clocks on both the server and every client, each clock starts at 0 and follows two rules:

Before sending any request, we add +1 to the current clock and adds it to the messages lamport field
When recieving a new message, we set our clock to max(current_own_clock, receieved_timestamp) + 1

This ensures that all clients and server always keeps up with the timestamp even if they have fallen behind, and makes sure they can't get out of sync since the highest current timestamp either our own or the one recieved becomes the truth and all other parts do the same, so no one can ever fall behind when recieving or sending anything.

## Diagram showing a sequence of RPC calls with their Lamport timestamps

```
Client X                      Server                      Client Y
   clock=0                      clock=2                      clock=3
      |                            |                            |
    1 |--- Join, L=1 ------------->|                            |
      |                            | max(2,1)+1 = 3             |
      |<-- "X joined", L=3 --------|--- "X joined", L=3 ------->|
    4 |                            |                            | 4
      |                            |                            |
    5 |--- Publish "Hello", L=5 -->|                            |
      |<-- Empty ------------------|                            |
      |                            | max(3,5)+1 = 6             |
      |<-- "Hello", L=6 -----------|--- "Hello", L=6 ---------->|
    7 |                            |                            | 7
      |                            |                            |
      |                            |<-- Publish "Hi X", L=8 ----| 8
      |                            |--- Empty ----------------->|
      |                            | max(6,8)+1 = 9             |
      |<-- "Hi X", L=9 ------------|--- "Hi X", L=9 ----------->|
   10 |                            |                            | 10
      |                            |                            |
   11 |--- Leave, L=11 ----------->|                            |
      |<-- Empty ------------------|                            |
      |<-- stream closed ----------|                            |
                                   | max(9,11)+1 = 12           |
                                   |--- "X left", L=12 -------->|
                                   |                            | 13
```

# Logs

## Server

```
[Server] 2026/10/25 23:18:16 event=STARTUP address=:5050
[Server] 2026/10/25 23:18:45 event=CONNECT client=Test lamport=1
[Server] 2026/10/25 23:18:45 event=BROADCAST sender=Server lamport=2 text="Participant Test joined Chit Chat at logical time 2"
[Server] 2026/10/25 23:18:45 event=DELIVER client=Test lamport=2
[Server] 2026/10/25 23:18:57 event=CONNECT client=Test2 lamport=1
[Server] 2026/10/25 23:18:57 event=BROADCAST sender=Server lamport=3 text="Participant Test2 joined Chit Chat at logical time 3"
[Server] 2026/10/25 23:18:57 event=DELIVER client=Test2 lamport=3
[Server] 2026/10/25 23:18:57 event=DELIVER client=Test lamport=3
[Server] 2026/10/25 23:19:09 event=CONNECT client=Test3 lamport=1
[Server] 2026/10/25 23:19:09 event=BROADCAST sender=Server lamport=4 text="Participant Test3 joined Chit Chat at logical time 4"
[Server] 2026/10/25 23:19:09 event=DELIVER client=Test3 lamport=4
[Server] 2026/10/25 23:19:09 event=DELIVER client=Test lamport=4
[Server] 2026/10/25 23:19:09 event=DELIVER client=Test2 lamport=4
[Server] 2026/10/25 23:19:30 event=BROADCAST sender=Test lamport=7 text="hej"
[Server] 2026/10/25 23:19:30 event=DELIVER client=Test3 lamport=7
[Server] 2026/10/25 23:19:30 event=DELIVER client=Test lamport=7
[Server] 2026/10/25 23:19:30 event=DELIVER client=Test2 lamport=7
[Server] 2026/10/25 23:19:32 event=BROADCAST sender=Test2 lamport=10 text="hej"
[Server] 2026/10/25 23:19:32 event=DELIVER client=Test lamport=10
[Server] 2026/10/25 23:19:32 event=DELIVER client=Test3 lamport=10
[Server] 2026/10/25 23:19:32 event=DELIVER client=Test2 lamport=10
[Server] 2026/10/25 23:19:35 event=BROADCAST sender=Test3 lamport=13 text="hej"
[Server] 2026/10/25 23:19:35 event=DELIVER client=Test lamport=13
[Server] 2026/10/25 23:19:35 event=DELIVER client=Test3 lamport=13
[Server] 2026/10/25 23:19:35 event=DELIVER client=Test2 lamport=13
[Server] 2026/10/25 23:21:22 event=BROADCAST sender=Test2 lamport=16 text="hej"
[Server] 2026/10/25 23:21:22 event=DELIVER client=Test lamport=16
[Server] 2026/10/25 23:21:22 event=DELIVER client=Test2 lamport=16
[Server] 2026/10/25 23:21:22 event=DELIVER client=Test3 lamport=16
[Server] 2026/10/25 23:21:22 event=BROADCAST sender=Test2 lamport=19 text="hej"
[Server] 2026/10/25 23:21:22 event=DELIVER client=Test lamport=19
[Server] 2026/10/25 23:21:22 event=DELIVER client=Test2 lamport=19
[Server] 2026/10/25 23:21:22 event=DELIVER client=Test3 lamport=19
[Server] 2026/10/25 23:21:23 event=BROADCAST sender=Test2 lamport=22 text="hej"
[Server] 2026/10/25 23:21:23 event=DELIVER client=Test lamport=22
[Server] 2026/10/25 23:21:23 event=DELIVER client=Test3 lamport=22
[Server] 2026/10/25 23:21:23 event=DELIVER client=Test2 lamport=22
[Server] 2026/10/25 23:21:23 event=BROADCAST sender=Test2 lamport=25 text="hej"
[Server] 2026/10/25 23:21:23 event=DELIVER client=Test2 lamport=25
[Server] 2026/10/25 23:21:23 event=DELIVER client=Test3 lamport=25
[Server] 2026/10/25 23:21:23 event=DELIVER client=Test lamport=25
[Server] 2026/10/25 23:21:24 event=BROADCAST sender=Test2 lamport=28 text="hej"
[Server] 2026/10/25 23:21:24 event=DELIVER client=Test3 lamport=28
[Server] 2026/10/25 23:21:24 event=DELIVER client=Test2 lamport=28
[Server] 2026/10/25 23:21:24 event=DELIVER client=Test lamport=28
[Server] 2026/10/25 23:21:30 event=BROADCAST sender=Test lamport=31 text="yoooooooooooooooo"
[Server] 2026/10/25 23:21:30 event=DELIVER client=Test lamport=31
[Server] 2026/10/25 23:21:30 event=DELIVER client=Test3 lamport=31
[Server] 2026/10/25 23:21:30 event=DELIVER client=Test2 lamport=31
[Server] 2026/10/25 23:21:32 event=BROADCAST sender=Test3 lamport=34 text="123123123123"
[Server] 2026/10/25 23:21:32 event=DELIVER client=Test lamport=34
[Server] 2026/10/25 23:21:32 event=DELIVER client=Test3 lamport=34
[Server] 2026/10/25 23:21:32 event=DELIVER client=Test2 lamport=34
[Server] 2026/10/25 23:21:35 event=BROADCAST sender=Test3 lamport=37 text="32132132132131"
[Server] 2026/10/25 23:21:35 event=DELIVER client=Test lamport=37
[Server] 2026/10/25 23:21:35 event=DELIVER client=Test3 lamport=37
[Server] 2026/10/25 23:21:35 event=DELIVER client=Test2 lamport=37
[Server] 2026/10/25 23:21:38 event=BROADCAST sender=Test lamport=40 text="555555555555"
[Server] 2026/10/25 23:21:38 event=DELIVER client=Test lamport=40
[Server] 2026/10/25 23:21:38 event=DELIVER client=Test3 lamport=40
[Server] 2026/10/25 23:21:38 event=DELIVER client=Test2 lamport=40
[Server] 2026/10/25 23:21:50 event=DISCONNECT client=Test3 lamport=42
[Server] 2026/10/25 23:21:50 event=BROADCAST sender=Server lamport=43 text="Participant Test3 left Chit Chat at logical time 43"
[Server] 2026/10/25 23:21:50 event=DELIVER client=Test2 lamport=43
[Server] 2026/10/25 23:21:50 event=DELIVER client=Test lamport=43
[Server] 2026/10/25 23:21:52 event=DISCONNECT client=Test lamport=45
[Server] 2026/10/25 23:21:52 event=BROADCAST sender=Server lamport=46 text="Participant Test left Chit Chat at logical time 46"
[Server] 2026/10/25 23:21:52 event=DELIVER client=Test2 lamport=46
[Server] 2026/10/25 23:21:54 event=BROADCAST sender=Test2 lamport=49 text="rip"
[Server] 2026/10/25 23:21:54 event=DELIVER client=Test2 lamport=49
[Server] 2026/10/25 23:21:56 event=BROADCAST sender=Test2 lamport=52 text="all alone"
[Server] 2026/10/25 23:21:56 event=DELIVER client=Test2 lamport=52
[Server] 2026/10/25 23:21:58 event=DISCONNECT client=Test2 lamport=54
[Server] 2026/10/25 23:21:58 event=BROADCAST sender=Server lamport=55 text="Participant Test2 left Chit Chat at logical time 55"
[Server] 2026/10/25 23:24:16 event=SHUTDOWN
```

## Client (Test)

```
[Client] 2026/10/25 23:18:45 event=CONNECT client=Test
[Client] 2026/10/25 23:18:45 event=RECEIVE client=Test sender=Server lamport=2 local=3 text="Participant Test joined Chit Chat at logical time 2"
[Client] 2026/10/25 23:18:57 event=RECEIVE client=Test sender=Server lamport=3 local=4 text="Participant Test2 joined Chit Chat at logical time 3"
[Client] 2026/10/25 23:19:09 event=RECEIVE client=Test sender=Server lamport=4 local=5 text="Participant Test3 joined Chit Chat at logical time 4"
[Client] 2026/10/25 23:19:30 event=PUBLISH client=Test lamport=6 text="hej"
[Client] 2026/10/25 23:19:30 event=RECEIVE client=Test sender=Test lamport=7 local=8 text="hej"
[Client] 2026/10/25 23:19:32 event=RECEIVE client=Test sender=Test2 lamport=10 local=11 text="hej"
[Client] 2026/10/25 23:19:35 event=RECEIVE client=Test sender=Test3 lamport=13 local=14 text="hej"
[Client] 2026/10/25 23:21:22 event=RECEIVE client=Test sender=Test2 lamport=16 local=17 text="hej"
[Client] 2026/10/25 23:21:22 event=RECEIVE client=Test sender=Test2 lamport=19 local=20 text="hej"
[Client] 2026/10/25 23:21:23 event=RECEIVE client=Test sender=Test2 lamport=22 local=23 text="hej"
[Client] 2026/10/25 23:21:23 event=RECEIVE client=Test sender=Test2 lamport=25 local=26 text="hej"
[Client] 2026/10/25 23:21:24 event=RECEIVE client=Test sender=Test2 lamport=28 local=29 text="hej"
[Client] 2026/10/25 23:21:30 event=PUBLISH client=Test lamport=30 text="yoooooooooooooooo"
[Client] 2026/10/25 23:21:30 event=RECEIVE client=Test sender=Test lamport=31 local=32 text="yoooooooooooooooo"
[Client] 2026/10/25 23:21:32 event=RECEIVE client=Test sender=Test3 lamport=34 local=35 text="123123123123"
[Client] 2026/10/25 23:21:35 event=RECEIVE client=Test sender=Test3 lamport=37 local=38 text="32132132132131"
[Client] 2026/10/25 23:21:38 event=PUBLISH client=Test lamport=39 text="555555555555"
[Client] 2026/10/25 23:21:38 event=RECEIVE client=Test sender=Test lamport=40 local=41 text="555555555555"
[Client] 2026/10/25 23:21:50 event=RECEIVE client=Test sender=Server lamport=43 local=44 text="Participant Test3 left Chit Chat at logical time 43"
[Client] 2026/10/25 23:21:52 event=DISCONNECT client=Test
```

## Client (Test2)

```
[Client] 2026/10/25 23:18:57 event=CONNECT client=Test2
[Client] 2026/10/25 23:18:57 event=RECEIVE client=Test2 sender=Server lamport=3 local=4 text="Participant Test2 joined Chit Chat at logical time 3"
[Client] 2026/10/25 23:19:09 event=RECEIVE client=Test2 sender=Server lamport=4 local=5 text="Participant Test3 joined Chit Chat at logical time 4"
[Client] 2026/10/25 23:19:30 event=RECEIVE client=Test2 sender=Test lamport=7 local=8 text="hej"
[Client] 2026/10/25 23:19:32 event=PUBLISH client=Test2 lamport=9 text="hej"
[Client] 2026/10/25 23:19:32 event=RECEIVE client=Test2 sender=Test2 lamport=10 local=11 text="hej"
[Client] 2026/10/25 23:19:35 event=RECEIVE client=Test2 sender=Test3 lamport=13 local=14 text="hej"
[Client] 2026/10/25 23:21:22 event=PUBLISH client=Test2 lamport=15 text="hej"
[Client] 2026/10/25 23:21:22 event=RECEIVE client=Test2 sender=Test2 lamport=16 local=17 text="hej"
[Client] 2026/10/25 23:21:22 event=PUBLISH client=Test2 lamport=18 text="hej"
[Client] 2026/10/25 23:21:22 event=RECEIVE client=Test2 sender=Test2 lamport=19 local=20 text="hej"
[Client] 2026/10/25 23:21:23 event=PUBLISH client=Test2 lamport=21 text="hej"
[Client] 2026/10/25 23:21:23 event=RECEIVE client=Test2 sender=Test2 lamport=22 local=23 text="hej"
[Client] 2026/10/25 23:21:23 event=PUBLISH client=Test2 lamport=24 text="hej"
[Client] 2026/10/25 23:21:23 event=RECEIVE client=Test2 sender=Test2 lamport=25 local=26 text="hej"
[Client] 2026/10/25 23:21:24 event=PUBLISH client=Test2 lamport=27 text="hej"
[Client] 2026/10/25 23:21:24 event=RECEIVE client=Test2 sender=Test2 lamport=28 local=29 text="hej"
[Client] 2026/10/25 23:21:30 event=RECEIVE client=Test2 sender=Test lamport=31 local=32 text="yoooooooooooooooo"
[Client] 2026/10/25 23:21:32 event=RECEIVE client=Test2 sender=Test3 lamport=34 local=35 text="123123123123"
[Client] 2026/10/25 23:21:35 event=RECEIVE client=Test2 sender=Test3 lamport=37 local=38 text="32132132132131"
[Client] 2026/10/25 23:21:38 event=RECEIVE client=Test2 sender=Test lamport=40 local=41 text="555555555555"
[Client] 2026/10/25 23:21:50 event=RECEIVE client=Test2 sender=Server lamport=43 local=44 text="Participant Test3 left Chit Chat at logical time 43"
[Client] 2026/10/25 23:21:52 event=RECEIVE client=Test2 sender=Server lamport=46 local=47 text="Participant Test left Chit Chat at logical time 46"
[Client] 2026/10/25 23:21:54 event=PUBLISH client=Test2 lamport=48 text="rip"
[Client] 2026/10/25 23:21:54 event=RECEIVE client=Test2 sender=Test2 lamport=49 local=50 text="rip"
[Client] 2026/10/25 23:21:56 event=PUBLISH client=Test2 lamport=51 text="all alone"
[Client] 2026/10/25 23:21:56 event=RECEIVE client=Test2 sender=Test2 lamport=52 local=53 text="all alone"
[Client] 2026/10/25 23:21:58 event=DISCONNECT client=Test2
```

## Client (Test3)

```
[Client] 2026/10/25 23:19:09 event=CONNECT client=Test3
[Client] 2026/10/25 23:19:09 event=RECEIVE client=Test3 sender=Server lamport=4 local=5 text="Participant Test3 joined Chit Chat at logical time 4"
[Client] 2026/10/25 23:19:30 event=RECEIVE client=Test3 sender=Test lamport=7 local=8 text="hej"
[Client] 2026/10/25 23:19:32 event=RECEIVE client=Test3 sender=Test2 lamport=10 local=11 text="hej"
[Client] 2026/10/25 23:19:35 event=PUBLISH client=Test3 lamport=12 text="hej"
[Client] 2026/10/25 23:19:35 event=RECEIVE client=Test3 sender=Test3 lamport=13 local=14 text="hej"
[Client] 2026/10/25 23:21:22 event=RECEIVE client=Test3 sender=Test2 lamport=16 local=17 text="hej"
[Client] 2026/10/25 23:21:22 event=RECEIVE client=Test3 sender=Test2 lamport=19 local=20 text="hej"
[Client] 2026/10/25 23:21:23 event=RECEIVE client=Test3 sender=Test2 lamport=22 local=23 text="hej"
[Client] 2026/10/25 23:21:23 event=RECEIVE client=Test3 sender=Test2 lamport=25 local=26 text="hej"
[Client] 2026/10/25 23:21:24 event=RECEIVE client=Test3 sender=Test2 lamport=28 local=29 text="hej"
[Client] 2026/10/25 23:21:30 event=RECEIVE client=Test3 sender=Test lamport=31 local=32 text="yoooooooooooooooo"
[Client] 2026/10/25 23:21:32 event=PUBLISH client=Test3 lamport=33 text="123123123123"
[Client] 2026/10/25 23:21:32 event=RECEIVE client=Test3 sender=Test3 lamport=34 local=35 text="123123123123"
[Client] 2026/10/25 23:21:35 event=PUBLISH client=Test3 lamport=36 text="32132132132131"
[Client] 2026/10/25 23:21:35 event=RECEIVE client=Test3 sender=Test3 lamport=37 local=38 text="32132132132131"
[Client] 2026/10/25 23:21:38 event=RECEIVE client=Test3 sender=Test lamport=40 local=41 text="555555555555"
[Client] 2026/10/25 23:21:50 event=DISCONNECT client=Test3
```