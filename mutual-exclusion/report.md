# Report

## Code Source Repo
https://github.com/OssieFromDK/DistributedSystems/tree/main/mutual-exclusion

## Describe how your system meets the System Requirements (R1, R2, and R3)

R1 is met since we have a system where 3 nodes can start, and each node acts as both a gRPC server and a client to the two other nodes, which lets them all talk to each other as is done in Ricart-Argawala. We emulate the critical section by logging ENTER_CS, then waiting for 1 second and then loggin EXIT_CS. Each node asks for access to the critical section on their own random timer which means any node can request access at any given time.

R2 is met as a node can only enter after both the other nodes have replied to its request to enter. A node will hold back its reply while it itself is in the Critical Section or while its waiting itself to enter the critical section and its request is older than the one that just came. They compare based on the lamport clock, and if a tie by their id where the lower one wins. Both nodes will be comparing the same timestamps, so only one of them can ever get all the replies and enter, while the other will be waiting its turn. 

R3 is met since a node will only hold back replies to request that are newer than its own or while its in the Critical Section, and sends all of them when it leaves. That means the oldest request always gets it next. 

## Provide a discussion of your algorithm, using examples from your logs (as per Technical Requirement 5)

We have implemented Ricart-Agrawala with Lamport clocks, where a node that wishes to enter the critical section, sets itself to WANTED, and then sends a request with its lamport time to the other nodes in the system, and it will not enter until all the other nodes have answered. The reply is the response of the gRPC call so holding a reply just means that the call does not return until the node exits.

The log lines below is from the beginning of our run, here we can see that Node 2 asks first with lamport 1, and then node 1 replies right away. Then node 1 asks with lamport 4 and node 2 holds that reply back because 1 is less than 4. Then node 3 starts running a little bit later and replies to both requests, so now node 2 enters the critical section first, and then when node 2 exits the critical section, it sends the held reply back to node 1, which now has both replies and then enters.

```
2026/11/08 21:48:27.369724 [Node 1] event=STARTUP address=localhost:5001 peers=2
2026/11/08 21:48:28.583766 [Node 2] event=STARTUP address=localhost:5002 peers=2
2026/11/08 21:48:29.792411 [Node 2] event=REQUEST_CS lamport=1
2026/11/08 21:48:29.792411 [Node 2] event=SEND_REQUEST to=1 lamport=1
2026/11/08 21:48:29.792411 [Node 2] event=SEND_REQUEST to=3 lamport=1
2026/11/08 21:48:29.799414 [Node 1] event=RECEIVE_REQUEST from=2 lamport=1 local=2 state=RELEASED
2026/11/08 21:48:29.799414 [Node 1] event=SEND_REPLY to=2 lamport=3
2026/11/08 21:48:29.800412 [Node 2] event=RECEIVE_REPLY from=1 lamport=3 local=4 replies=1/2
2026/11/08 21:48:29.895411 [Node 1] event=REQUEST_CS lamport=4
2026/11/08 21:48:29.895411 [Node 1] event=SEND_REQUEST to=2 lamport=4
2026/11/08 21:48:29.895411 [Node 1] event=SEND_REQUEST to=3 lamport=4
2026/11/08 21:48:29.896411 [Node 2] event=RECEIVE_REQUEST from=1 lamport=4 local=5 state=WANTED
2026/11/08 21:48:29.896411 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:48:31.339760 [Node 3] event=STARTUP address=localhost:5003 peers=2
2026/11/08 21:48:31.800801 [Node 3] event=RECEIVE_REQUEST from=2 lamport=1 local=2 state=RELEASED
2026/11/08 21:48:31.800801 [Node 3] event=SEND_REPLY to=2 lamport=3
2026/11/08 21:48:31.800801 [Node 2] event=RECEIVE_REPLY from=3 lamport=3 local=6 replies=2/2
2026/11/08 21:48:31.800801 [Node 2] event=ENTER_CS lamport=6
2026/11/08 21:48:31.897149 [Node 3] event=RECEIVE_REQUEST from=1 lamport=4 local=5 state=RELEASED
2026/11/08 21:48:31.898150 [Node 3] event=SEND_REPLY to=1 lamport=6
2026/11/08 21:48:31.898150 [Node 1] event=RECEIVE_REPLY from=3 lamport=6 local=7 replies=1/2
2026/11/08 21:48:32.800921 [Node 2] event=EXIT_CS lamport=6 deferred=1
2026/11/08 21:48:32.800921 [Node 2] event=SEND_REPLY to=1 lamport=7
2026/11/08 21:48:32.800921 [Node 1] event=RECEIVE_REPLY from=2 lamport=7 local=8 replies=2/2
2026/11/08 21:48:32.801624 [Node 1] event=ENTER_CS lamport=8
```

## Logs

### Node 1

```
2026/11/08 21:48:27.369724 [Node 1] event=STARTUP address=localhost:5001 peers=2
2026/11/08 21:48:29.799414 [Node 1] event=RECEIVE_REQUEST from=2 lamport=1 local=2 state=RELEASED
2026/11/08 21:48:29.799414 [Node 1] event=SEND_REPLY to=2 lamport=3
2026/11/08 21:48:29.895411 [Node 1] event=REQUEST_CS lamport=4
2026/11/08 21:48:29.895411 [Node 1] event=SEND_REQUEST to=2 lamport=4
2026/11/08 21:48:29.895411 [Node 1] event=SEND_REQUEST to=3 lamport=4
2026/11/08 21:48:31.898150 [Node 1] event=RECEIVE_REPLY from=3 lamport=6 local=7 replies=1/2
2026/11/08 21:48:32.800921 [Node 1] event=RECEIVE_REPLY from=2 lamport=7 local=8 replies=2/2
2026/11/08 21:48:32.801624 [Node 1] event=ENTER_CS lamport=8
2026/11/08 21:48:33.485550 [Node 1] event=RECEIVE_REQUEST from=3 lamport=7 local=9 state=HELD
2026/11/08 21:48:33.486073 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:48:33.802829 [Node 1] event=EXIT_CS lamport=9 deferred=1
2026/11/08 21:48:33.802864 [Node 1] event=SEND_REPLY to=3 lamport=10
2026/11/08 21:48:33.946452 [Node 1] event=RECEIVE_REQUEST from=2 lamport=10 local=11 state=RELEASED
2026/11/08 21:48:33.946452 [Node 1] event=SEND_REPLY to=2 lamport=12
2026/11/08 21:48:35.152282 [Node 1] event=REQUEST_CS lamport=13
2026/11/08 21:48:35.152282 [Node 1] event=SEND_REQUEST to=2 lamport=13
2026/11/08 21:48:35.152282 [Node 1] event=SEND_REQUEST to=3 lamport=13
2026/11/08 21:48:35.152835 [Node 1] event=RECEIVE_REPLY from=3 lamport=15 local=16 replies=1/2
2026/11/08 21:48:35.804863 [Node 1] event=RECEIVE_REPLY from=2 lamport=16 local=17 replies=2/2
2026/11/08 21:48:35.804863 [Node 1] event=ENTER_CS lamport=17
2026/11/08 21:48:36.758632 [Node 1] event=RECEIVE_REQUEST from=3 lamport=16 local=18 state=HELD
2026/11/08 21:48:36.758632 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:48:36.805741 [Node 1] event=EXIT_CS lamport=18 deferred=1
2026/11/08 21:48:36.805787 [Node 1] event=SEND_REPLY to=3 lamport=19
2026/11/08 21:48:38.236784 [Node 1] event=RECEIVE_REQUEST from=2 lamport=19 local=20 state=RELEASED
2026/11/08 21:48:38.236784 [Node 1] event=SEND_REPLY to=2 lamport=21
2026/11/08 21:48:39.961400 [Node 1] event=REQUEST_CS lamport=22
2026/11/08 21:48:39.961400 [Node 1] event=SEND_REQUEST to=2 lamport=22
2026/11/08 21:48:39.961400 [Node 1] event=SEND_REQUEST to=3 lamport=22
2026/11/08 21:48:39.962007 [Node 1] event=RECEIVE_REPLY from=3 lamport=24 local=25 replies=1/2
2026/11/08 21:48:39.962007 [Node 1] event=RECEIVE_REPLY from=2 lamport=26 local=27 replies=2/2
2026/11/08 21:48:39.962007 [Node 1] event=ENTER_CS lamport=27
2026/11/08 21:48:40.054925 [Node 1] event=RECEIVE_REQUEST from=3 lamport=25 local=28 state=HELD
2026/11/08 21:48:40.054925 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:48:40.962924 [Node 1] event=EXIT_CS lamport=28 deferred=1
2026/11/08 21:48:40.962924 [Node 1] event=SEND_REPLY to=3 lamport=29
2026/11/08 21:48:42.037145 [Node 1] event=RECEIVE_REQUEST from=2 lamport=29 local=30 state=RELEASED
2026/11/08 21:48:42.037145 [Node 1] event=SEND_REPLY to=2 lamport=31
2026/11/08 21:48:43.447654 [Node 1] event=RECEIVE_REQUEST from=3 lamport=33 local=34 state=RELEASED
2026/11/08 21:48:43.447654 [Node 1] event=SEND_REPLY to=3 lamport=35
2026/11/08 21:48:43.970199 [Node 1] event=REQUEST_CS lamport=36
2026/11/08 21:48:43.970199 [Node 1] event=SEND_REQUEST to=2 lamport=36
2026/11/08 21:48:43.970199 [Node 1] event=SEND_REQUEST to=3 lamport=36
2026/11/08 21:48:43.970199 [Node 1] event=RECEIVE_REPLY from=2 lamport=38 local=39 replies=1/2
2026/11/08 21:48:44.449919 [Node 1] event=RECEIVE_REPLY from=3 lamport=39 local=40 replies=2/2
2026/11/08 21:48:44.449919 [Node 1] event=ENTER_CS lamport=40
2026/11/08 21:48:45.451586 [Node 1] event=EXIT_CS lamport=40 deferred=0
2026/11/08 21:48:45.740213 [Node 1] event=RECEIVE_REQUEST from=2 lamport=39 local=41 state=RELEASED
2026/11/08 21:48:45.740213 [Node 1] event=SEND_REPLY to=2 lamport=42
2026/11/08 21:48:45.859942 [Node 1] event=RECEIVE_REQUEST from=3 lamport=42 local=43 state=RELEASED
2026/11/08 21:48:45.859942 [Node 1] event=SEND_REPLY to=3 lamport=44
2026/11/08 21:48:47.095362 [Node 1] event=REQUEST_CS lamport=45
2026/11/08 21:48:47.095362 [Node 1] event=SEND_REQUEST to=2 lamport=45
2026/11/08 21:48:47.095362 [Node 1] event=SEND_REQUEST to=3 lamport=45
2026/11/08 21:48:47.095362 [Node 1] event=RECEIVE_REPLY from=2 lamport=48 local=49 replies=1/2
2026/11/08 21:48:47.742785 [Node 1] event=RECEIVE_REPLY from=3 lamport=49 local=50 replies=2/2
2026/11/08 21:48:47.742785 [Node 1] event=ENTER_CS lamport=50
2026/11/08 21:48:48.213898 [Node 1] event=RECEIVE_REQUEST from=2 lamport=49 local=51 state=HELD
2026/11/08 21:48:48.213898 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:48:48.744171 [Node 1] event=EXIT_CS lamport=51 deferred=1
2026/11/08 21:48:48.744171 [Node 1] event=SEND_REPLY to=2 lamport=52
2026/11/08 21:48:49.449148 [Node 1] event=RECEIVE_REQUEST from=3 lamport=52 local=53 state=RELEASED
2026/11/08 21:48:49.449148 [Node 1] event=SEND_REPLY to=3 lamport=54
2026/11/08 21:48:51.229880 [Node 1] event=REQUEST_CS lamport=55
2026/11/08 21:48:51.229880 [Node 1] event=SEND_REQUEST to=2 lamport=55
2026/11/08 21:48:51.229880 [Node 1] event=SEND_REQUEST to=3 lamport=55
2026/11/08 21:48:51.229880 [Node 1] event=RECEIVE_REPLY from=2 lamport=57 local=58 replies=1/2
2026/11/08 21:48:51.229880 [Node 1] event=RECEIVE_REPLY from=3 lamport=58 local=59 replies=2/2
2026/11/08 21:48:51.229880 [Node 1] event=ENTER_CS lamport=59
2026/11/08 21:48:52.230606 [Node 1] event=EXIT_CS lamport=59 deferred=0
2026/11/08 21:48:52.287606 [Node 1] event=RECEIVE_REQUEST from=3 lamport=59 local=60 state=RELEASED
2026/11/08 21:48:52.287606 [Node 1] event=SEND_REPLY to=3 lamport=61
2026/11/08 21:48:53.076607 [Node 1] event=RECEIVE_REQUEST from=2 lamport=62 local=63 state=RELEASED
2026/11/08 21:48:53.076607 [Node 1] event=SEND_REPLY to=2 lamport=64
2026/11/08 21:48:53.866702 [Node 1] event=REQUEST_CS lamport=65
2026/11/08 21:48:53.866702 [Node 1] event=SEND_REQUEST to=2 lamport=65
2026/11/08 21:48:53.866702 [Node 1] event=SEND_REQUEST to=3 lamport=65
2026/11/08 21:48:53.866702 [Node 1] event=RECEIVE_REPLY from=3 lamport=67 local=68 replies=1/2
2026/11/08 21:48:54.290144 [Node 1] event=RECEIVE_REPLY from=2 lamport=68 local=69 replies=2/2
2026/11/08 21:48:54.290144 [Node 1] event=ENTER_CS lamport=69
2026/11/08 21:48:55.290483 [Node 1] event=EXIT_CS lamport=69 deferred=0
2026/11/08 21:48:55.402492 [Node 1] event=RECEIVE_REQUEST from=3 lamport=68 local=70 state=RELEASED
2026/11/08 21:48:55.402492 [Node 1] event=SEND_REPLY to=3 lamport=71
2026/11/08 21:48:57.587893 [Node 1] event=REQUEST_CS lamport=72
2026/11/08 21:48:57.587893 [Node 1] event=SEND_REQUEST to=2 lamport=72
2026/11/08 21:48:57.587893 [Node 1] event=SEND_REQUEST to=3 lamport=72
2026/11/08 21:48:57.588557 [Node 1] event=RECEIVE_REPLY from=3 lamport=75 local=76 replies=1/2
2026/11/08 21:48:57.588557 [Node 1] event=RECEIVE_REPLY from=2 lamport=74 local=77 replies=2/2
2026/11/08 21:48:57.588557 [Node 1] event=ENTER_CS lamport=77
2026/11/08 21:48:58.099310 [Node 1] event=RECEIVE_REQUEST from=2 lamport=75 local=78 state=HELD
2026/11/08 21:48:58.099310 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:48:58.589816 [Node 1] event=EXIT_CS lamport=78 deferred=1
2026/11/08 21:48:58.589816 [Node 1] event=SEND_REPLY to=2 lamport=79
2026/11/08 21:48:59.650893 [Node 1] event=RECEIVE_REQUEST from=3 lamport=78 local=80 state=RELEASED
2026/11/08 21:48:59.650893 [Node 1] event=SEND_REPLY to=3 lamport=81
2026/11/08 21:48:59.682893 [Node 1] event=REQUEST_CS lamport=82
2026/11/08 21:48:59.682893 [Node 1] event=SEND_REQUEST to=3 lamport=82
2026/11/08 21:48:59.682893 [Node 1] event=SEND_REQUEST to=2 lamport=82
2026/11/08 21:48:59.682893 [Node 1] event=RECEIVE_REPLY from=2 lamport=84 local=85 replies=1/2
2026/11/08 21:49:00.651894 [Node 1] event=RECEIVE_REPLY from=3 lamport=86 local=87 replies=2/2
2026/11/08 21:49:00.651894 [Node 1] event=ENTER_CS lamport=87
2026/11/08 21:49:01.614867 [Node 1] event=RECEIVE_REQUEST from=2 lamport=85 local=88 state=HELD
2026/11/08 21:49:01.614867 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:01.652290 [Node 1] event=EXIT_CS lamport=88 deferred=1
2026/11/08 21:49:01.652335 [Node 1] event=SEND_REPLY to=2 lamport=89
2026/11/08 21:49:04.278450 [Node 1] event=RECEIVE_REQUEST from=3 lamport=89 local=90 state=RELEASED
2026/11/08 21:49:04.278450 [Node 1] event=SEND_REPLY to=3 lamport=91
2026/11/08 21:49:04.588451 [Node 1] event=REQUEST_CS lamport=92
2026/11/08 21:49:04.588451 [Node 1] event=SEND_REQUEST to=3 lamport=92
2026/11/08 21:49:04.588451 [Node 1] event=SEND_REQUEST to=2 lamport=92
2026/11/08 21:49:04.589452 [Node 1] event=RECEIVE_REPLY from=2 lamport=94 local=95 replies=1/2
2026/11/08 21:49:05.280415 [Node 1] event=RECEIVE_REPLY from=3 lamport=95 local=96 replies=2/2
2026/11/08 21:49:05.280415 [Node 1] event=ENTER_CS lamport=96
2026/11/08 21:49:05.787007 [Node 1] event=RECEIVE_REQUEST from=2 lamport=95 local=97 state=HELD
2026/11/08 21:49:05.787007 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:06.280857 [Node 1] event=EXIT_CS lamport=97 deferred=1
2026/11/08 21:49:06.280857 [Node 1] event=SEND_REPLY to=2 lamport=98
2026/11/08 21:49:06.726118 [Node 1] event=RECEIVE_REQUEST from=3 lamport=98 local=99 state=RELEASED
2026/11/08 21:49:06.726118 [Node 1] event=SEND_REPLY to=3 lamport=100
2026/11/08 21:49:08.570254 [Node 1] event=REQUEST_CS lamport=101
2026/11/08 21:49:08.570309 [Node 1] event=SEND_REQUEST to=2 lamport=101
2026/11/08 21:49:08.570309 [Node 1] event=SEND_REQUEST to=3 lamport=101
2026/11/08 21:49:08.570835 [Node 1] event=RECEIVE_REPLY from=3 lamport=104 local=105 replies=1/2
2026/11/08 21:49:08.570835 [Node 1] event=RECEIVE_REPLY from=2 lamport=103 local=106 replies=2/2
2026/11/08 21:49:08.570835 [Node 1] event=ENTER_CS lamport=106
2026/11/08 21:49:09.571559 [Node 1] event=EXIT_CS lamport=106 deferred=0
2026/11/08 21:49:09.828145 [Node 1] event=RECEIVE_REQUEST from=3 lamport=105 local=107 state=RELEASED
2026/11/08 21:49:09.828722 [Node 1] event=SEND_REPLY to=3 lamport=108
2026/11/08 21:49:10.941783 [Node 1] event=RECEIVE_REQUEST from=2 lamport=108 local=109 state=RELEASED
2026/11/08 21:49:10.941783 [Node 1] event=SEND_REPLY to=2 lamport=110
2026/11/08 21:49:11.700783 [Node 1] event=REQUEST_CS lamport=111
2026/11/08 21:49:11.700783 [Node 1] event=SEND_REQUEST to=3 lamport=111
2026/11/08 21:49:11.700783 [Node 1] event=SEND_REQUEST to=2 lamport=111
2026/11/08 21:49:11.700783 [Node 1] event=RECEIVE_REPLY from=3 lamport=113 local=114 replies=1/2
2026/11/08 21:49:11.943782 [Node 1] event=RECEIVE_REPLY from=2 lamport=115 local=116 replies=2/2
2026/11/08 21:49:11.943782 [Node 1] event=ENTER_CS lamport=116
2026/11/08 21:49:12.943808 [Node 1] event=EXIT_CS lamport=116 deferred=0
2026/11/08 21:49:13.886785 [Node 1] event=RECEIVE_REQUEST from=3 lamport=114 local=117 state=RELEASED
2026/11/08 21:49:13.886785 [Node 1] event=SEND_REPLY to=3 lamport=118
2026/11/08 21:49:14.228215 [Node 1] event=RECEIVE_REQUEST from=2 lamport=118 local=119 state=RELEASED
2026/11/08 21:49:14.228215 [Node 1] event=SEND_REPLY to=2 lamport=120
2026/11/08 21:49:14.473214 [Node 1] event=REQUEST_CS lamport=121
2026/11/08 21:49:14.473214 [Node 1] event=SEND_REQUEST to=2 lamport=121
2026/11/08 21:49:14.473214 [Node 1] event=SEND_REQUEST to=3 lamport=121
2026/11/08 21:49:14.889214 [Node 1] event=RECEIVE_REPLY from=3 lamport=123 local=124 replies=1/2
2026/11/08 21:49:15.891001 [Node 1] event=RECEIVE_REPLY from=2 lamport=126 local=127 replies=2/2
2026/11/08 21:49:15.891001 [Node 1] event=ENTER_CS lamport=127
2026/11/08 21:49:16.891991 [Node 1] event=EXIT_CS lamport=127 deferred=0
2026/11/08 21:49:17.161578 [Node 1] event=RECEIVE_REQUEST from=3 lamport=125 local=128 state=RELEASED
2026/11/08 21:49:17.161578 [Node 1] event=SEND_REPLY to=3 lamport=129
2026/11/08 21:49:19.579330 [Node 1] event=REQUEST_CS lamport=130
2026/11/08 21:49:19.579330 [Node 1] event=SEND_REQUEST to=2 lamport=130
2026/11/08 21:49:19.579330 [Node 1] event=SEND_REQUEST to=3 lamport=130
2026/11/08 21:49:19.579330 [Node 1] event=RECEIVE_REPLY from=2 lamport=132 local=133 replies=1/2
2026/11/08 21:49:19.580454 [Node 1] event=RECEIVE_REPLY from=3 lamport=132 local=134 replies=2/2
2026/11/08 21:49:19.580454 [Node 1] event=ENTER_CS lamport=134
2026/11/08 21:49:19.772539 [Node 1] event=RECEIVE_REQUEST from=2 lamport=133 local=135 state=HELD
2026/11/08 21:49:19.772539 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:20.580955 [Node 1] event=EXIT_CS lamport=135 deferred=1
2026/11/08 21:49:20.580955 [Node 1] event=SEND_REPLY to=2 lamport=136
2026/11/08 21:49:20.948917 [Node 1] event=RECEIVE_REQUEST from=3 lamport=136 local=137 state=RELEASED
2026/11/08 21:49:20.948917 [Node 1] event=SEND_REPLY to=3 lamport=138
2026/11/08 21:49:23.396838 [Node 1] event=REQUEST_CS lamport=139
2026/11/08 21:49:23.396838 [Node 1] event=SEND_REQUEST to=2 lamport=139
2026/11/08 21:49:23.396838 [Node 1] event=SEND_REQUEST to=3 lamport=139
2026/11/08 21:49:23.397455 [Node 1] event=RECEIVE_REPLY from=2 lamport=141 local=142 replies=1/2
2026/11/08 21:49:23.397455 [Node 1] event=RECEIVE_REPLY from=3 lamport=142 local=143 replies=2/2
2026/11/08 21:49:23.397455 [Node 1] event=ENTER_CS lamport=143
2026/11/08 21:49:23.883087 [Node 1] event=RECEIVE_REQUEST from=3 lamport=143 local=144 state=HELD
2026/11/08 21:49:23.883087 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:49:23.897211 [Node 1] event=RECEIVE_REQUEST from=2 lamport=146 local=147 state=HELD
2026/11/08 21:49:23.897211 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:24.398741 [Node 1] event=EXIT_CS lamport=147 deferred=2
2026/11/08 21:49:24.398741 [Node 1] event=SEND_REPLY to=2 lamport=148
2026/11/08 21:49:24.398741 [Node 1] event=SEND_REPLY to=3 lamport=149
2026/11/08 21:49:27.834612 [Node 1] event=RECEIVE_REQUEST from=3 lamport=152 local=153 state=RELEASED
2026/11/08 21:49:27.834612 [Node 1] event=SEND_REPLY to=3 lamport=154
2026/11/08 21:49:28.170386 [Node 1] event=REQUEST_CS lamport=155
2026/11/08 21:49:28.170386 [Node 1] event=SEND_REQUEST to=2 lamport=155
2026/11/08 21:49:28.170386 [Node 1] event=SEND_REQUEST to=3 lamport=155
2026/11/08 21:49:28.170386 [Node 1] event=RECEIVE_REPLY from=2 lamport=157 local=158 replies=1/2
2026/11/08 21:49:28.836313 [Node 1] event=RECEIVE_REPLY from=3 lamport=158 local=159 replies=2/2
2026/11/08 21:49:28.836313 [Node 1] event=ENTER_CS lamport=159
2026/11/08 21:49:29.309702 [Node 1] event=RECEIVE_REQUEST from=2 lamport=158 local=160 state=HELD
2026/11/08 21:49:29.309702 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:29.837702 [Node 1] event=EXIT_CS lamport=160 deferred=1
2026/11/08 21:49:29.837702 [Node 1] event=SEND_REPLY to=2 lamport=161
2026/11/08 21:49:31.660819 [Node 1] event=RECEIVE_REQUEST from=3 lamport=161 local=162 state=RELEASED
2026/11/08 21:49:31.660819 [Node 1] event=SEND_REPLY to=3 lamport=163
2026/11/08 21:49:33.609940 [Node 1] event=REQUEST_CS lamport=164
2026/11/08 21:49:33.609940 [Node 1] event=SEND_REQUEST to=2 lamport=164
2026/11/08 21:49:33.609940 [Node 1] event=SEND_REQUEST to=3 lamport=164
2026/11/08 21:49:33.609940 [Node 1] event=RECEIVE_REPLY from=2 lamport=166 local=167 replies=1/2
2026/11/08 21:49:33.610924 [Node 1] event=RECEIVE_REPLY from=3 lamport=167 local=168 replies=2/2
2026/11/08 21:49:33.610947 [Node 1] event=ENTER_CS lamport=168
2026/11/08 21:49:34.316255 [Node 1] event=RECEIVE_REQUEST from=2 lamport=167 local=169 state=HELD
2026/11/08 21:49:34.316255 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:34.611253 [Node 1] event=EXIT_CS lamport=169 deferred=1
2026/11/08 21:49:34.611253 [Node 1] event=SEND_REPLY to=2 lamport=170
2026/11/08 21:49:35.650781 [Node 1] event=RECEIVE_REQUEST from=3 lamport=170 local=171 state=RELEASED
2026/11/08 21:49:35.650781 [Node 1] event=SEND_REPLY to=3 lamport=172
2026/11/08 21:49:36.405970 [Node 1] event=REQUEST_CS lamport=173
2026/11/08 21:49:36.405970 [Node 1] event=SEND_REQUEST to=2 lamport=173
2026/11/08 21:49:36.405970 [Node 1] event=SEND_REQUEST to=3 lamport=173
2026/11/08 21:49:36.405970 [Node 1] event=RECEIVE_REPLY from=2 lamport=175 local=176 replies=1/2
2026/11/08 21:49:36.652969 [Node 1] event=RECEIVE_REPLY from=3 lamport=177 local=178 replies=2/2
2026/11/08 21:49:36.652969 [Node 1] event=ENTER_CS lamport=178
2026/11/08 21:49:37.652971 [Node 1] event=EXIT_CS lamport=178 deferred=0
2026/11/08 21:49:38.729037 [Node 1] event=REQUEST_CS lamport=179
2026/11/08 21:49:38.729037 [Node 1] event=SEND_REQUEST to=2 lamport=179
2026/11/08 21:49:38.729037 [Node 1] event=SEND_REQUEST to=3 lamport=179
2026/11/08 21:49:38.729037 [Node 1] event=RECEIVE_REPLY from=2 lamport=181 local=182 replies=1/2
2026/11/08 21:49:38.729037 [Node 1] event=RECEIVE_REPLY from=3 lamport=181 local=183 replies=2/2
2026/11/08 21:49:38.729037 [Node 1] event=ENTER_CS lamport=183
2026/11/08 21:49:39.228496 [Node 1] event=RECEIVE_REQUEST from=2 lamport=182 local=184 state=HELD
2026/11/08 21:49:39.228496 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:39.729545 [Node 1] event=EXIT_CS lamport=184 deferred=1
2026/11/08 21:49:39.729545 [Node 1] event=SEND_REPLY to=2 lamport=185
2026/11/08 21:49:39.780758 [Node 1] event=RECEIVE_REQUEST from=3 lamport=185 local=186 state=RELEASED
2026/11/08 21:49:39.780758 [Node 1] event=SEND_REPLY to=3 lamport=187
2026/11/08 21:49:42.308091 [Node 1] event=REQUEST_CS lamport=188
2026/11/08 21:49:42.308091 [Node 1] event=SEND_REQUEST to=3 lamport=188
2026/11/08 21:49:42.308091 [Node 1] event=SEND_REQUEST to=2 lamport=188
2026/11/08 21:49:42.308091 [Node 1] event=RECEIVE_REPLY from=2 lamport=190 local=191 replies=1/2
2026/11/08 21:49:42.308091 [Node 1] event=RECEIVE_REPLY from=3 lamport=191 local=192 replies=2/2
2026/11/08 21:49:42.309091 [Node 1] event=ENTER_CS lamport=192
2026/11/08 21:49:43.309729 [Node 1] event=EXIT_CS lamport=192 deferred=0
2026/11/08 21:49:43.605729 [Node 1] event=RECEIVE_REQUEST from=2 lamport=191 local=193 state=RELEASED
2026/11/08 21:49:43.605729 [Node 1] event=SEND_REPLY to=2 lamport=194
2026/11/08 21:49:44.261728 [Node 1] event=RECEIVE_REQUEST from=3 lamport=194 local=195 state=RELEASED
2026/11/08 21:49:44.261728 [Node 1] event=SEND_REPLY to=3 lamport=196
2026/11/08 21:49:44.914259 [Node 1] event=REQUEST_CS lamport=197
2026/11/08 21:49:44.914259 [Node 1] event=SEND_REQUEST to=2 lamport=197
2026/11/08 21:49:44.914259 [Node 1] event=SEND_REQUEST to=3 lamport=197
2026/11/08 21:49:44.914259 [Node 1] event=RECEIVE_REPLY from=2 lamport=200 local=201 replies=1/2
2026/11/08 21:49:45.608261 [Node 1] event=RECEIVE_REPLY from=3 lamport=201 local=202 replies=2/2
2026/11/08 21:49:45.608261 [Node 1] event=ENTER_CS lamport=202
2026/11/08 21:49:46.608317 [Node 1] event=EXIT_CS lamport=202 deferred=0
2026/11/08 21:49:46.695382 [Node 1] event=RECEIVE_REQUEST from=3 lamport=202 local=203 state=RELEASED
2026/11/08 21:49:46.695382 [Node 1] event=SEND_REPLY to=3 lamport=204
2026/11/08 21:49:47.403553 [Node 1] event=RECEIVE_REQUEST from=2 lamport=205 local=206 state=RELEASED
2026/11/08 21:49:47.403553 [Node 1] event=SEND_REPLY to=2 lamport=207
2026/11/08 21:49:47.912668 [Node 1] event=REQUEST_CS lamport=208
2026/11/08 21:49:47.912668 [Node 1] event=SEND_REQUEST to=2 lamport=208
2026/11/08 21:49:47.912668 [Node 1] event=SEND_REQUEST to=3 lamport=208
2026/11/08 21:49:47.913316 [Node 1] event=RECEIVE_REPLY from=3 lamport=210 local=211 replies=1/2
2026/11/08 21:49:48.697503 [Node 1] event=RECEIVE_REPLY from=2 lamport=211 local=212 replies=2/2
2026/11/08 21:49:48.697503 [Node 1] event=ENTER_CS lamport=212
2026/11/08 21:49:49.698021 [Node 1] event=EXIT_CS lamport=212 deferred=0
2026/11/08 21:49:50.583648 [Node 1] event=RECEIVE_REQUEST from=3 lamport=211 local=213 state=RELEASED
2026/11/08 21:49:50.583648 [Node 1] event=SEND_REPLY to=3 lamport=214
2026/11/08 21:49:51.439324 [Node 1] event=RECEIVE_REQUEST from=2 lamport=214 local=215 state=RELEASED
2026/11/08 21:49:51.439324 [Node 1] event=SEND_REPLY to=2 lamport=216
2026/11/08 21:49:52.155451 [Node 1] event=REQUEST_CS lamport=217
2026/11/08 21:49:52.155451 [Node 1] event=SEND_REQUEST to=3 lamport=217
2026/11/08 21:49:52.155451 [Node 1] event=SEND_REQUEST to=2 lamport=217
2026/11/08 21:49:52.155451 [Node 1] event=RECEIVE_REPLY from=3 lamport=220 local=221 replies=1/2
2026/11/08 21:49:52.586453 [Node 1] event=RECEIVE_REPLY from=2 lamport=221 local=222 replies=2/2
2026/11/08 21:49:52.586453 [Node 1] event=ENTER_CS lamport=222
2026/11/08 21:49:53.587451 [Node 1] event=EXIT_CS lamport=222 deferred=0
2026/11/08 21:49:54.621989 [Node 1] event=RECEIVE_REQUEST from=2 lamport=222 local=223 state=RELEASED
2026/11/08 21:49:54.621989 [Node 1] event=SEND_REPLY to=2 lamport=224
2026/11/08 21:49:54.935990 [Node 1] event=RECEIVE_REQUEST from=3 lamport=225 local=226 state=RELEASED
2026/11/08 21:49:54.935990 [Node 1] event=SEND_REPLY to=3 lamport=227
2026/11/08 21:49:57.227981 [Node 1] event=REQUEST_CS lamport=228
2026/11/08 21:49:57.227981 [Node 1] event=SEND_REQUEST to=2 lamport=228
2026/11/08 21:49:57.227981 [Node 1] event=SEND_REQUEST to=3 lamport=228
2026/11/08 21:49:57.228903 [Node 1] event=RECEIVE_REPLY from=3 lamport=231 local=232 replies=1/2
2026/11/08 21:49:57.228903 [Node 1] event=RECEIVE_REPLY from=2 lamport=230 local=233 replies=2/2
2026/11/08 21:49:57.228903 [Node 1] event=ENTER_CS lamport=233
2026/11/08 21:49:58.135041 [Node 1] event=RECEIVE_REQUEST from=2 lamport=231 local=234 state=HELD
2026/11/08 21:49:58.135041 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:49:58.229905 [Node 1] event=EXIT_CS lamport=234 deferred=1
2026/11/08 21:49:58.229905 [Node 1] event=SEND_REPLY to=2 lamport=235
2026/11/08 21:49:59.739551 [Node 1] event=RECEIVE_REQUEST from=3 lamport=234 local=236 state=RELEASED
2026/11/08 21:49:59.739551 [Node 1] event=SEND_REPLY to=3 lamport=237
2026/11/08 21:50:01.475551 [Node 1] event=REQUEST_CS lamport=238
2026/11/08 21:50:01.475551 [Node 1] event=SEND_REQUEST to=2 lamport=238
2026/11/08 21:50:01.475551 [Node 1] event=SEND_REQUEST to=3 lamport=238
2026/11/08 21:50:01.476551 [Node 1] event=RECEIVE_REPLY from=2 lamport=240 local=241 replies=1/2
2026/11/08 21:50:01.476551 [Node 1] event=RECEIVE_REPLY from=3 lamport=242 local=243 replies=2/2
2026/11/08 21:50:01.476551 [Node 1] event=ENTER_CS lamport=243
2026/11/08 21:50:02.384094 [Node 1] event=RECEIVE_REQUEST from=3 lamport=243 local=244 state=HELD
2026/11/08 21:50:02.384094 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:50:02.477195 [Node 1] event=EXIT_CS lamport=244 deferred=1
2026/11/08 21:50:02.477253 [Node 1] event=SEND_REPLY to=3 lamport=245
2026/11/08 21:50:02.756389 [Node 1] event=RECEIVE_REQUEST from=2 lamport=246 local=247 state=RELEASED
2026/11/08 21:50:02.756389 [Node 1] event=SEND_REPLY to=2 lamport=248
2026/11/08 21:50:03.561384 [Node 1] event=REQUEST_CS lamport=249
2026/11/08 21:50:03.561384 [Node 1] event=SEND_REQUEST to=2 lamport=249
2026/11/08 21:50:03.561384 [Node 1] event=SEND_REQUEST to=3 lamport=249
2026/11/08 21:50:03.561384 [Node 1] event=RECEIVE_REPLY from=3 lamport=251 local=252 replies=1/2
2026/11/08 21:50:04.478686 [Node 1] event=RECEIVE_REPLY from=2 lamport=252 local=253 replies=2/2
2026/11/08 21:50:04.478686 [Node 1] event=ENTER_CS lamport=253
2026/11/08 21:50:05.478985 [Node 1] event=EXIT_CS lamport=253 deferred=0
2026/11/08 21:50:06.191544 [Node 1] event=RECEIVE_REQUEST from=3 lamport=252 local=254 state=RELEASED
2026/11/08 21:50:06.191544 [Node 1] event=SEND_REPLY to=3 lamport=255
2026/11/08 21:50:07.034858 [Node 1] event=RECEIVE_REQUEST from=2 lamport=255 local=256 state=RELEASED
2026/11/08 21:50:07.034858 [Node 1] event=SEND_REPLY to=2 lamport=257
2026/11/08 21:50:09.472858 [Node 1] event=REQUEST_CS lamport=258
2026/11/08 21:50:09.472858 [Node 1] event=SEND_REQUEST to=2 lamport=258
2026/11/08 21:50:09.472858 [Node 1] event=SEND_REQUEST to=3 lamport=258
2026/11/08 21:50:09.472858 [Node 1] event=RECEIVE_REPLY from=2 lamport=261 local=262 replies=1/2
2026/11/08 21:50:09.472858 [Node 1] event=RECEIVE_REPLY from=3 lamport=260 local=263 replies=2/2
2026/11/08 21:50:09.472858 [Node 1] event=ENTER_CS lamport=263
2026/11/08 21:50:10.473387 [Node 1] event=EXIT_CS lamport=263 deferred=0
2026/11/08 21:50:11.005396 [Node 1] event=RECEIVE_REQUEST from=3 lamport=261 local=264 state=RELEASED
2026/11/08 21:50:11.005396 [Node 1] event=SEND_REPLY to=3 lamport=265
2026/11/08 21:50:11.893394 [Node 1] event=RECEIVE_REQUEST from=2 lamport=264 local=266 state=RELEASED
2026/11/08 21:50:11.893394 [Node 1] event=SEND_REPLY to=2 lamport=267
2026/11/08 21:50:13.579938 [Node 1] event=RECEIVE_REQUEST from=3 lamport=270 local=271 state=RELEASED
2026/11/08 21:50:13.579938 [Node 1] event=SEND_REPLY to=3 lamport=272
2026/11/08 21:50:14.056938 [Node 1] event=REQUEST_CS lamport=273
2026/11/08 21:50:14.056938 [Node 1] event=SEND_REQUEST to=2 lamport=273
2026/11/08 21:50:14.056938 [Node 1] event=SEND_REQUEST to=3 lamport=273
2026/11/08 21:50:14.057940 [Node 1] event=RECEIVE_REPLY from=2 lamport=275 local=276 replies=1/2
2026/11/08 21:50:14.581496 [Node 1] event=RECEIVE_REPLY from=3 lamport=276 local=277 replies=2/2
2026/11/08 21:50:14.581496 [Node 1] event=ENTER_CS lamport=277
2026/11/08 21:50:15.583069 [Node 1] event=EXIT_CS lamport=277 deferred=0
2026/11/08 21:50:15.910069 [Node 1] event=RECEIVE_REQUEST from=2 lamport=276 local=278 state=RELEASED
2026/11/08 21:50:15.910069 [Node 1] event=SEND_REPLY to=2 lamport=279
2026/11/08 21:50:16.652068 [Node 1] event=REQUEST_CS lamport=280
2026/11/08 21:50:16.652068 [Node 1] event=SEND_REQUEST to=2 lamport=280
2026/11/08 21:50:16.652068 [Node 1] event=SEND_REQUEST to=3 lamport=280
2026/11/08 21:50:16.652068 [Node 1] event=RECEIVE_REPLY from=3 lamport=282 local=283 replies=1/2
2026/11/08 21:50:16.912069 [Node 1] event=RECEIVE_REPLY from=2 lamport=282 local=284 replies=2/2
2026/11/08 21:50:16.912069 [Node 1] event=ENTER_CS lamport=284
2026/11/08 21:50:17.913068 [Node 1] event=EXIT_CS lamport=284 deferred=0
2026/11/08 21:50:18.016069 [Node 1] event=RECEIVE_REQUEST from=2 lamport=283 local=285 state=RELEASED
2026/11/08 21:50:18.016069 [Node 1] event=SEND_REPLY to=2 lamport=286
2026/11/08 21:50:18.124070 [Node 1] event=RECEIVE_REQUEST from=3 lamport=286 local=287 state=RELEASED
2026/11/08 21:50:18.124070 [Node 1] event=SEND_REPLY to=3 lamport=288
2026/11/08 21:50:21.542646 [Node 1] event=REQUEST_CS lamport=289
2026/11/08 21:50:21.542646 [Node 1] event=SEND_REQUEST to=2 lamport=289
2026/11/08 21:50:21.542646 [Node 1] event=SEND_REQUEST to=3 lamport=289
2026/11/08 21:50:21.542646 [Node 1] event=RECEIVE_REPLY from=3 lamport=293 local=294 replies=1/2
2026/11/08 21:50:21.543646 [Node 1] event=RECEIVE_REPLY from=2 lamport=292 local=295 replies=2/2
2026/11/08 21:50:21.543646 [Node 1] event=ENTER_CS lamport=295
2026/11/08 21:50:21.758495 [Node 1] event=RECEIVE_REQUEST from=3 lamport=294 local=296 state=HELD
2026/11/08 21:50:21.758495 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:50:22.380813 [Node 1] event=RECEIVE_REQUEST from=2 lamport=297 local=298 state=HELD
2026/11/08 21:50:22.380813 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:50:22.543813 [Node 1] event=EXIT_CS lamport=298 deferred=2
2026/11/08 21:50:22.543813 [Node 1] event=SEND_REPLY to=2 lamport=299
2026/11/08 21:50:22.543813 [Node 1] event=SEND_REPLY to=3 lamport=300
2026/11/08 21:50:23.561760 [Node 1] event=REQUEST_CS lamport=301
2026/11/08 21:50:23.561760 [Node 1] event=SEND_REQUEST to=2 lamport=301
2026/11/08 21:50:23.561760 [Node 1] event=SEND_REQUEST to=3 lamport=301
2026/11/08 21:50:23.562374 [Node 1] event=RECEIVE_REPLY from=3 lamport=304 local=305 replies=1/2
2026/11/08 21:50:24.545549 [Node 1] event=RECEIVE_REPLY from=2 lamport=305 local=306 replies=2/2
2026/11/08 21:50:24.545549 [Node 1] event=ENTER_CS lamport=306
2026/11/08 21:50:25.545985 [Node 1] event=EXIT_CS lamport=306 deferred=0
2026/11/08 21:50:26.408605 [Node 1] event=RECEIVE_REQUEST from=3 lamport=305 local=307 state=RELEASED
2026/11/08 21:50:26.408605 [Node 1] event=SEND_REPLY to=3 lamport=308
2026/11/08 21:50:26.701148 [Node 1] event=RECEIVE_REQUEST from=2 lamport=308 local=309 state=RELEASED
2026/11/08 21:50:26.701148 [Node 1] event=SEND_REPLY to=2 lamport=310
2026/11/08 21:50:28.363709 [Node 1] event=REQUEST_CS lamport=311
2026/11/08 21:50:28.363709 [Node 1] event=SEND_REQUEST to=2 lamport=311
2026/11/08 21:50:28.363709 [Node 1] event=SEND_REQUEST to=3 lamport=311
2026/11/08 21:50:28.363709 [Node 1] event=RECEIVE_REPLY from=3 lamport=313 local=314 replies=1/2
2026/11/08 21:50:28.410709 [Node 1] event=RECEIVE_REPLY from=2 lamport=314 local=315 replies=2/2
2026/11/08 21:50:28.410709 [Node 1] event=ENTER_CS lamport=315
2026/11/08 21:50:28.448709 [Node 1] event=RECEIVE_REQUEST from=3 lamport=314 local=316 state=HELD
2026/11/08 21:50:28.448709 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:50:29.411345 [Node 1] event=EXIT_CS lamport=316 deferred=1
2026/11/08 21:50:29.411345 [Node 1] event=SEND_REPLY to=3 lamport=317
2026/11/08 21:50:29.934347 [Node 1] event=RECEIVE_REQUEST from=2 lamport=317 local=318 state=RELEASED
2026/11/08 21:50:29.934347 [Node 1] event=SEND_REPLY to=2 lamport=319
2026/11/08 21:50:30.714556 [Node 1] event=REQUEST_CS lamport=320
2026/11/08 21:50:30.714556 [Node 1] event=SEND_REQUEST to=2 lamport=320
2026/11/08 21:50:30.714556 [Node 1] event=SEND_REQUEST to=3 lamport=320
2026/11/08 21:50:30.714556 [Node 1] event=RECEIVE_REPLY from=3 lamport=322 local=323 replies=1/2
2026/11/08 21:50:31.412735 [Node 1] event=RECEIVE_REPLY from=2 lamport=323 local=324 replies=2/2
2026/11/08 21:50:31.412735 [Node 1] event=ENTER_CS lamport=324
2026/11/08 21:50:32.413519 [Node 1] event=EXIT_CS lamport=324 deferred=0
2026/11/08 21:50:32.501116 [Node 1] event=RECEIVE_REQUEST from=2 lamport=324 local=325 state=RELEASED
2026/11/08 21:50:32.501116 [Node 1] event=SEND_REPLY to=2 lamport=326
2026/11/08 21:50:32.508253 [Node 1] event=RECEIVE_REQUEST from=3 lamport=327 local=328 state=RELEASED
2026/11/08 21:50:32.508253 [Node 1] event=SEND_REPLY to=3 lamport=329
2026/11/08 21:50:34.506439 [Node 1] event=REQUEST_CS lamport=330
2026/11/08 21:50:34.506439 [Node 1] event=SEND_REQUEST to=2 lamport=330
2026/11/08 21:50:34.506439 [Node 1] event=SEND_REQUEST to=3 lamport=330
2026/11/08 21:50:34.507029 [Node 1] event=RECEIVE_REPLY from=2 lamport=332 local=333 replies=1/2
2026/11/08 21:50:34.507029 [Node 1] event=RECEIVE_REPLY from=3 lamport=333 local=334 replies=2/2
2026/11/08 21:50:34.507029 [Node 1] event=ENTER_CS lamport=334
2026/11/08 21:50:35.507350 [Node 1] event=EXIT_CS lamport=334 deferred=0
2026/11/08 21:50:35.699017 [Node 1] event=RECEIVE_REQUEST from=3 lamport=334 local=335 state=RELEASED
2026/11/08 21:50:35.699593 [Node 1] event=SEND_REPLY to=3 lamport=336
2026/11/08 21:50:36.498921 [Node 1] event=RECEIVE_REQUEST from=2 lamport=337 local=338 state=RELEASED
2026/11/08 21:50:36.499443 [Node 1] event=SEND_REPLY to=2 lamport=339
2026/11/08 21:50:36.931396 [Node 1] event=REQUEST_CS lamport=340
2026/11/08 21:50:36.931396 [Node 1] event=SEND_REQUEST to=2 lamport=340
2026/11/08 21:50:36.931396 [Node 1] event=SEND_REQUEST to=3 lamport=340
2026/11/08 21:50:36.931926 [Node 1] event=RECEIVE_REPLY from=3 lamport=342 local=343 replies=1/2
2026/11/08 21:50:37.700641 [Node 1] event=RECEIVE_REPLY from=2 lamport=343 local=344 replies=2/2
2026/11/08 21:50:37.701168 [Node 1] event=ENTER_CS lamport=344
2026/11/08 21:50:38.702143 [Node 1] event=EXIT_CS lamport=344 deferred=0
2026/11/08 21:50:39.572932 [Node 1] event=RECEIVE_REQUEST from=3 lamport=343 local=345 state=RELEASED
2026/11/08 21:50:39.573491 [Node 1] event=SEND_REPLY to=3 lamport=346
2026/11/08 21:50:39.882745 [Node 1] event=RECEIVE_REQUEST from=2 lamport=346 local=347 state=RELEASED
2026/11/08 21:50:39.882745 [Node 1] event=SEND_REPLY to=2 lamport=348
2026/11/08 21:50:41.881401 [Node 1] event=REQUEST_CS lamport=349
2026/11/08 21:50:41.881401 [Node 1] event=SEND_REQUEST to=2 lamport=349
2026/11/08 21:50:41.881401 [Node 1] event=SEND_REQUEST to=3 lamport=349
2026/11/08 21:50:41.881978 [Node 1] event=RECEIVE_REPLY from=2 lamport=352 local=353 replies=1/2
2026/11/08 21:50:41.881978 [Node 1] event=RECEIVE_REPLY from=3 lamport=351 local=354 replies=2/2
2026/11/08 21:50:41.881978 [Node 1] event=ENTER_CS lamport=354
2026/11/08 21:50:42.882363 [Node 1] event=EXIT_CS lamport=354 deferred=0
2026/11/08 21:50:43.829352 [Node 1] event=RECEIVE_REQUEST from=3 lamport=352 local=355 state=RELEASED
2026/11/08 21:50:43.829352 [Node 1] event=SEND_REPLY to=3 lamport=356
2026/11/08 21:50:44.588719 [Node 1] event=RECEIVE_REQUEST from=2 lamport=355 local=357 state=RELEASED
2026/11/08 21:50:44.588719 [Node 1] event=SEND_REPLY to=2 lamport=358
2026/11/08 21:50:45.122762 [Node 1] event=REQUEST_CS lamport=359
2026/11/08 21:50:45.122762 [Node 1] event=SEND_REQUEST to=2 lamport=359
2026/11/08 21:50:45.122762 [Node 1] event=SEND_REQUEST to=3 lamport=359
2026/11/08 21:50:45.123333 [Node 1] event=RECEIVE_REPLY from=3 lamport=361 local=362 replies=1/2
2026/11/08 21:50:45.831263 [Node 1] event=RECEIVE_REPLY from=2 lamport=362 local=363 replies=2/2
2026/11/08 21:50:45.831263 [Node 1] event=ENTER_CS lamport=363
2026/11/08 21:50:46.831693 [Node 1] event=EXIT_CS lamport=363 deferred=0
2026/11/08 21:50:47.287963 [Node 1] event=RECEIVE_REQUEST from=3 lamport=362 local=364 state=RELEASED
2026/11/08 21:50:47.287988 [Node 1] event=SEND_REPLY to=3 lamport=365
2026/11/08 21:50:48.178033 [Node 1] event=RECEIVE_REQUEST from=2 lamport=365 local=366 state=RELEASED
2026/11/08 21:50:48.178556 [Node 1] event=SEND_REPLY to=2 lamport=367
2026/11/08 21:50:49.050316 [Node 1] event=REQUEST_CS lamport=368
2026/11/08 21:50:49.050316 [Node 1] event=SEND_REQUEST to=3 lamport=368
2026/11/08 21:50:49.050316 [Node 1] event=SEND_REQUEST to=2 lamport=368
2026/11/08 21:50:49.050927 [Node 1] event=RECEIVE_REPLY from=3 lamport=370 local=371 replies=1/2
2026/11/08 21:50:49.289473 [Node 1] event=RECEIVE_REPLY from=2 lamport=371 local=372 replies=2/2
2026/11/08 21:50:49.289473 [Node 1] event=ENTER_CS lamport=372
2026/11/08 21:50:49.849200 [Node 1] event=RECEIVE_REQUEST from=3 lamport=371 local=373 state=HELD
2026/11/08 21:50:49.849200 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:50:50.289501 [Node 1] event=EXIT_CS lamport=373 deferred=1
2026/11/08 21:50:50.289501 [Node 1] event=SEND_REPLY to=3 lamport=374
2026/11/08 21:50:50.317103 [Node 1] event=RECEIVE_REQUEST from=2 lamport=374 local=375 state=RELEASED
2026/11/08 21:50:50.317103 [Node 1] event=SEND_REPLY to=2 lamport=376
2026/11/08 21:50:52.644315 [Node 1] event=REQUEST_CS lamport=377
2026/11/08 21:50:52.644315 [Node 1] event=SEND_REQUEST to=2 lamport=377
2026/11/08 21:50:52.644315 [Node 1] event=SEND_REQUEST to=3 lamport=377
2026/11/08 21:50:52.645262 [Node 1] event=RECEIVE_REPLY from=3 lamport=379 local=380 replies=1/2
2026/11/08 21:50:52.645262 [Node 1] event=RECEIVE_REPLY from=2 lamport=380 local=381 replies=2/2
2026/11/08 21:50:52.645262 [Node 1] event=ENTER_CS lamport=381
2026/11/08 21:50:53.374818 [Node 1] event=RECEIVE_REQUEST from=3 lamport=380 local=382 state=HELD
2026/11/08 21:50:53.374818 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:50:53.622034 [Node 1] event=RECEIVE_REQUEST from=2 lamport=383 local=384 state=HELD
2026/11/08 21:50:53.622034 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:50:53.645544 [Node 1] event=EXIT_CS lamport=384 deferred=2
2026/11/08 21:50:53.645544 [Node 1] event=SEND_REPLY to=2 lamport=385
2026/11/08 21:50:53.645544 [Node 1] event=SEND_REPLY to=3 lamport=386
2026/11/08 21:50:55.950152 [Node 1] event=RECEIVE_REQUEST from=3 lamport=389 local=390 state=RELEASED
2026/11/08 21:50:55.950152 [Node 1] event=SEND_REPLY to=3 lamport=391
2026/11/08 21:50:57.272979 [Node 1] event=RECEIVE_REQUEST from=2 lamport=392 local=393 state=RELEASED
2026/11/08 21:50:57.272979 [Node 1] event=SEND_REPLY to=2 lamport=394
2026/11/08 21:50:57.501619 [Node 1] event=REQUEST_CS lamport=395
2026/11/08 21:50:57.501653 [Node 1] event=SEND_REQUEST to=2 lamport=395
2026/11/08 21:50:57.501653 [Node 1] event=SEND_REQUEST to=3 lamport=395
2026/11/08 21:50:57.502180 [Node 1] event=RECEIVE_REPLY from=3 lamport=397 local=398 replies=1/2
2026/11/08 21:50:58.274523 [Node 1] event=RECEIVE_REPLY from=2 lamport=398 local=399 replies=2/2
2026/11/08 21:50:58.275078 [Node 1] event=ENTER_CS lamport=399
2026/11/08 21:50:58.818627 [Node 1] event=RECEIVE_REQUEST from=3 lamport=398 local=400 state=HELD
2026/11/08 21:50:58.818627 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:50:59.276239 [Node 1] event=EXIT_CS lamport=400 deferred=1
2026/11/08 21:50:59.276286 [Node 1] event=SEND_REPLY to=3 lamport=401
2026/11/08 21:50:59.568014 [Node 1] event=RECEIVE_REQUEST from=2 lamport=401 local=402 state=RELEASED
2026/11/08 21:50:59.568023 [Node 1] event=SEND_REPLY to=2 lamport=403
2026/11/08 21:51:01.674312 [Node 1] event=RECEIVE_REQUEST from=3 lamport=405 local=406 state=RELEASED
2026/11/08 21:51:01.674312 [Node 1] event=SEND_REPLY to=3 lamport=407
2026/11/08 21:51:02.599159 [Node 1] event=REQUEST_CS lamport=408
2026/11/08 21:51:02.599159 [Node 1] event=SEND_REQUEST to=3 lamport=408
2026/11/08 21:51:02.599159 [Node 1] event=SEND_REQUEST to=2 lamport=408
2026/11/08 21:51:02.600811 [Node 1] event=RECEIVE_REPLY from=2 lamport=410 local=411 replies=1/2
2026/11/08 21:51:02.675822 [Node 1] event=RECEIVE_REPLY from=3 lamport=411 local=412 replies=2/2
2026/11/08 21:51:02.675822 [Node 1] event=ENTER_CS lamport=412
2026/11/08 21:51:03.676874 [Node 1] event=EXIT_CS lamport=412 deferred=0
2026/11/08 21:51:04.803489 [Node 1] event=RECEIVE_REQUEST from=3 lamport=412 local=413 state=RELEASED
2026/11/08 21:51:04.803489 [Node 1] event=SEND_REPLY to=3 lamport=414
2026/11/08 21:51:05.067372 [Node 1] event=RECEIVE_REQUEST from=2 lamport=415 local=416 state=RELEASED
2026/11/08 21:51:05.067411 [Node 1] event=SEND_REPLY to=2 lamport=417
2026/11/08 21:51:07.492166 [Node 1] event=REQUEST_CS lamport=418
2026/11/08 21:51:07.492235 [Node 1] event=SEND_REQUEST to=2 lamport=418
2026/11/08 21:51:07.492235 [Node 1] event=SEND_REQUEST to=3 lamport=418
2026/11/08 21:51:07.492757 [Node 1] event=RECEIVE_REPLY from=2 lamport=421 local=422 replies=1/2
2026/11/08 21:51:07.493370 [Node 1] event=RECEIVE_REPLY from=3 lamport=420 local=423 replies=2/2
2026/11/08 21:51:07.493370 [Node 1] event=ENTER_CS lamport=423
2026/11/08 21:51:08.494322 [Node 1] event=EXIT_CS lamport=423 deferred=0
2026/11/08 21:51:08.831648 [Node 1] event=RECEIVE_REQUEST from=3 lamport=421 local=424 state=RELEASED
2026/11/08 21:51:08.831672 [Node 1] event=SEND_REPLY to=3 lamport=425
2026/11/08 21:51:09.863016 [Node 1] event=RECEIVE_REQUEST from=2 lamport=424 local=426 state=RELEASED
2026/11/08 21:51:09.863016 [Node 1] event=SEND_REPLY to=2 lamport=427
2026/11/08 21:51:12.243279 [Node 1] event=REQUEST_CS lamport=428
2026/11/08 21:51:12.243279 [Node 1] event=SEND_REQUEST to=2 lamport=428
2026/11/08 21:51:12.243279 [Node 1] event=SEND_REQUEST to=3 lamport=428
2026/11/08 21:51:12.243877 [Node 1] event=RECEIVE_REPLY from=3 lamport=430 local=431 replies=1/2
2026/11/08 21:51:12.243877 [Node 1] event=RECEIVE_REPLY from=2 lamport=432 local=433 replies=2/2
2026/11/08 21:51:12.243877 [Node 1] event=ENTER_CS lamport=433
2026/11/08 21:51:12.797544 [Node 1] event=RECEIVE_REQUEST from=3 lamport=431 local=434 state=HELD
2026/11/08 21:51:12.797544 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:51:13.106345 [Node 1] event=RECEIVE_REQUEST from=2 lamport=435 local=436 state=HELD
2026/11/08 21:51:13.106913 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:51:13.244530 [Node 1] event=EXIT_CS lamport=436 deferred=2
2026/11/08 21:51:13.244530 [Node 1] event=SEND_REPLY to=2 lamport=437
2026/11/08 21:51:13.244530 [Node 1] event=SEND_REPLY to=3 lamport=438
2026/11/08 21:51:15.760122 [Node 1] event=REQUEST_CS lamport=439
2026/11/08 21:51:15.760122 [Node 1] event=SEND_REQUEST to=2 lamport=439
2026/11/08 21:51:15.760122 [Node 1] event=SEND_REQUEST to=3 lamport=439
2026/11/08 21:51:15.760727 [Node 1] event=RECEIVE_REPLY from=2 lamport=443 local=444 replies=1/2
2026/11/08 21:51:15.761234 [Node 1] event=RECEIVE_REPLY from=3 lamport=442 local=445 replies=2/2
2026/11/08 21:51:15.761255 [Node 1] event=ENTER_CS lamport=445
2026/11/08 21:51:16.762053 [Node 1] event=EXIT_CS lamport=445 deferred=0
2026/11/08 21:51:17.267352 [Node 1] event=RECEIVE_REQUEST from=2 lamport=444 local=446 state=RELEASED
2026/11/08 21:51:17.267352 [Node 1] event=SEND_REPLY to=2 lamport=447
2026/11/08 21:51:17.930584 [Node 1] event=RECEIVE_REQUEST from=3 lamport=447 local=448 state=RELEASED
2026/11/08 21:51:17.931152 [Node 1] event=SEND_REPLY to=3 lamport=449
2026/11/08 21:51:18.685583 [Node 1] event=REQUEST_CS lamport=450
2026/11/08 21:51:18.685583 [Node 1] event=SEND_REQUEST to=2 lamport=450
2026/11/08 21:51:18.685583 [Node 1] event=SEND_REQUEST to=3 lamport=450
2026/11/08 21:51:18.686164 [Node 1] event=RECEIVE_REPLY from=2 lamport=452 local=453 replies=1/2
2026/11/08 21:51:19.269432 [Node 1] event=RECEIVE_REPLY from=3 lamport=453 local=454 replies=2/2
2026/11/08 21:51:19.269432 [Node 1] event=ENTER_CS lamport=454
2026/11/08 21:51:19.718281 [Node 1] event=RECEIVE_REQUEST from=2 lamport=453 local=455 state=HELD
2026/11/08 21:51:19.718281 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:51:20.270859 [Node 1] event=EXIT_CS lamport=455 deferred=1
2026/11/08 21:51:20.270859 [Node 1] event=SEND_REPLY to=2 lamport=456
2026/11/08 21:51:20.709725 [Node 1] event=RECEIVE_REQUEST from=3 lamport=456 local=457 state=RELEASED
2026/11/08 21:51:20.709748 [Node 1] event=SEND_REPLY to=3 lamport=458
2026/11/08 21:51:22.844979 [Node 1] event=RECEIVE_REQUEST from=2 lamport=460 local=461 state=RELEASED
2026/11/08 21:51:22.845716 [Node 1] event=SEND_REPLY to=2 lamport=462
2026/11/08 21:51:23.731975 [Node 1] event=REQUEST_CS lamport=463
2026/11/08 21:51:23.731975 [Node 1] event=SEND_REQUEST to=2 lamport=463
2026/11/08 21:51:23.731975 [Node 1] event=SEND_REQUEST to=3 lamport=463
2026/11/08 21:51:23.732557 [Node 1] event=RECEIVE_REPLY from=3 lamport=465 local=466 replies=1/2
2026/11/08 21:51:23.846732 [Node 1] event=RECEIVE_REPLY from=2 lamport=466 local=467 replies=2/2
2026/11/08 21:51:23.846732 [Node 1] event=ENTER_CS lamport=467
2026/11/08 21:51:24.847153 [Node 1] event=EXIT_CS lamport=467 deferred=0
2026/11/08 21:51:25.016300 [Node 1] event=RECEIVE_REQUEST from=3 lamport=466 local=468 state=RELEASED
2026/11/08 21:51:25.016300 [Node 1] event=SEND_REPLY to=3 lamport=469
2026/11/08 21:51:27.571080 [Node 1] event=RECEIVE_REQUEST from=2 lamport=469 local=470 state=RELEASED
2026/11/08 21:51:27.571080 [Node 1] event=SEND_REPLY to=2 lamport=471
2026/11/08 21:51:28.133799 [Node 1] event=REQUEST_CS lamport=472
2026/11/08 21:51:28.133799 [Node 1] event=SEND_REQUEST to=2 lamport=472
2026/11/08 21:51:28.133799 [Node 1] event=SEND_REQUEST to=3 lamport=472
2026/11/08 21:51:28.134388 [Node 1] event=RECEIVE_REPLY from=3 lamport=474 local=475 replies=1/2
2026/11/08 21:51:28.233981 [Node 1] event=RECEIVE_REQUEST from=3 lamport=475 local=476 state=WANTED
2026/11/08 21:51:28.233981 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:51:28.572385 [Node 1] event=RECEIVE_REPLY from=2 lamport=478 local=479 replies=2/2
2026/11/08 21:51:28.572385 [Node 1] event=ENTER_CS lamport=479
2026/11/08 21:51:29.573860 [Node 1] event=EXIT_CS lamport=479 deferred=1
2026/11/08 21:51:29.573913 [Node 1] event=SEND_REPLY to=3 lamport=480
2026/11/08 21:51:32.454762 [Node 1] event=RECEIVE_REQUEST from=2 lamport=479 local=481 state=RELEASED
2026/11/08 21:51:32.455582 [Node 1] event=SEND_REPLY to=2 lamport=482
2026/11/08 21:51:32.485501 [Node 1] event=RECEIVE_REQUEST from=3 lamport=484 local=485 state=RELEASED
2026/11/08 21:51:32.485501 [Node 1] event=SEND_REPLY to=3 lamport=486
2026/11/08 21:51:32.567111 [Node 1] event=REQUEST_CS lamport=487
2026/11/08 21:51:32.567111 [Node 1] event=SEND_REQUEST to=3 lamport=487
2026/11/08 21:51:32.567111 [Node 1] event=SEND_REQUEST to=2 lamport=487
2026/11/08 21:51:33.456325 [Node 1] event=RECEIVE_REPLY from=2 lamport=489 local=490 replies=1/2
2026/11/08 21:51:34.457984 [Node 1] event=RECEIVE_REPLY from=3 lamport=492 local=493 replies=2/2
2026/11/08 21:51:34.457984 [Node 1] event=ENTER_CS lamport=493
2026/11/08 21:51:34.644542 [Node 1] event=RECEIVE_REQUEST from=2 lamport=491 local=494 state=HELD
2026/11/08 21:51:34.644542 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:51:35.458799 [Node 1] event=EXIT_CS lamport=494 deferred=1
2026/11/08 21:51:35.458847 [Node 1] event=SEND_REPLY to=2 lamport=495
2026/11/08 21:51:36.411855 [Node 1] event=RECEIVE_REQUEST from=3 lamport=495 local=496 state=RELEASED
2026/11/08 21:51:36.411855 [Node 1] event=SEND_REPLY to=3 lamport=497
2026/11/08 21:51:36.864310 [Node 1] event=REQUEST_CS lamport=498
2026/11/08 21:51:36.864310 [Node 1] event=SEND_REQUEST to=2 lamport=498
2026/11/08 21:51:36.864310 [Node 1] event=SEND_REQUEST to=3 lamport=498
2026/11/08 21:51:36.864993 [Node 1] event=RECEIVE_REPLY from=2 lamport=500 local=501 replies=1/2
2026/11/08 21:51:37.460228 [Node 1] event=RECEIVE_REPLY from=3 lamport=501 local=502 replies=2/2
2026/11/08 21:51:37.460228 [Node 1] event=ENTER_CS lamport=502
2026/11/08 21:51:38.460707 [Node 1] event=EXIT_CS lamport=502 deferred=0
2026/11/08 21:51:39.849876 [Node 1] event=RECEIVE_REQUEST from=2 lamport=501 local=503 state=RELEASED
2026/11/08 21:51:39.849876 [Node 1] event=SEND_REPLY to=2 lamport=504
2026/11/08 21:51:40.152113 [Node 1] event=REQUEST_CS lamport=505
2026/11/08 21:51:40.152163 [Node 1] event=SEND_REQUEST to=2 lamport=505
2026/11/08 21:51:40.152163 [Node 1] event=SEND_REQUEST to=3 lamport=505
2026/11/08 21:51:40.152746 [Node 1] event=RECEIVE_REPLY from=3 lamport=507 local=508 replies=1/2
2026/11/08 21:51:40.852067 [Node 1] event=RECEIVE_REPLY from=2 lamport=508 local=509 replies=2/2
2026/11/08 21:51:40.852067 [Node 1] event=ENTER_CS lamport=509
2026/11/08 21:51:41.034166 [Node 1] event=RECEIVE_REQUEST from=3 lamport=508 local=510 state=HELD
2026/11/08 21:51:41.034166 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:51:41.853036 [Node 1] event=EXIT_CS lamport=510 deferred=1
2026/11/08 21:51:41.853036 [Node 1] event=SEND_REPLY to=3 lamport=511
2026/11/08 21:51:43.045604 [Node 1] event=REQUEST_CS lamport=512
2026/11/08 21:51:43.045604 [Node 1] event=SEND_REQUEST to=2 lamport=512
2026/11/08 21:51:43.045604 [Node 1] event=SEND_REQUEST to=3 lamport=512
2026/11/08 21:51:43.046191 [Node 1] event=RECEIVE_REPLY from=3 lamport=514 local=515 replies=1/2
2026/11/08 21:51:43.046191 [Node 1] event=RECEIVE_REPLY from=2 lamport=514 local=516 replies=2/2
2026/11/08 21:51:43.046191 [Node 1] event=ENTER_CS lamport=516
2026/11/08 21:51:43.783829 [Node 1] event=RECEIVE_REQUEST from=2 lamport=515 local=517 state=HELD
2026/11/08 21:51:43.783829 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:51:44.046782 [Node 1] event=EXIT_CS lamport=517 deferred=1
2026/11/08 21:51:44.046817 [Node 1] event=SEND_REPLY to=2 lamport=518
2026/11/08 21:51:46.319517 [Node 1] event=REQUEST_CS lamport=519
2026/11/08 21:51:46.319517 [Node 1] event=SEND_REQUEST to=3 lamport=519
2026/11/08 21:51:46.319517 [Node 1] event=SEND_REQUEST to=2 lamport=519
2026/11/08 21:51:46.320235 [Node 1] event=RECEIVE_REPLY from=3 lamport=521 local=522 replies=1/2
2026/11/08 21:51:46.320235 [Node 1] event=RECEIVE_REPLY from=2 lamport=521 local=523 replies=2/2
2026/11/08 21:51:46.320235 [Node 1] event=ENTER_CS lamport=523
2026/11/08 21:51:46.596278 [Node 1] event=RECEIVE_REQUEST from=3 lamport=522 local=524 state=HELD
2026/11/08 21:51:46.596278 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:51:47.320465 [Node 1] event=EXIT_CS lamport=524 deferred=1
2026/11/08 21:51:47.320465 [Node 1] event=SEND_REPLY to=3 lamport=525
2026/11/08 21:51:47.594630 [Node 1] event=RECEIVE_REQUEST from=2 lamport=525 local=526 state=RELEASED
2026/11/08 21:51:47.594630 [Node 1] event=SEND_REPLY to=2 lamport=527
2026/11/08 21:51:50.668144 [Node 1] event=REQUEST_CS lamport=528
2026/11/08 21:51:50.668144 [Node 1] event=SEND_REQUEST to=2 lamport=528
2026/11/08 21:51:50.668144 [Node 1] event=SEND_REQUEST to=3 lamport=528
2026/11/08 21:51:50.668727 [Node 1] event=RECEIVE_REPLY from=2 lamport=531 local=532 replies=1/2
2026/11/08 21:51:50.668727 [Node 1] event=RECEIVE_REPLY from=3 lamport=530 local=533 replies=2/2
2026/11/08 21:51:50.668727 [Node 1] event=ENTER_CS lamport=533
2026/11/08 21:51:51.668997 [Node 1] event=EXIT_CS lamport=533 deferred=0
2026/11/08 21:51:51.744163 [Node 1] event=RECEIVE_REQUEST from=3 lamport=531 local=534 state=RELEASED
2026/11/08 21:51:51.744163 [Node 1] event=SEND_REPLY to=3 lamport=535
2026/11/08 21:51:52.798728 [Node 1] event=RECEIVE_REQUEST from=2 lamport=534 local=536 state=RELEASED
2026/11/08 21:51:52.798728 [Node 1] event=SEND_REPLY to=2 lamport=537
2026/11/08 21:51:53.321526 [Node 1] event=REQUEST_CS lamport=538
2026/11/08 21:51:53.321526 [Node 1] event=SEND_REQUEST to=2 lamport=538
2026/11/08 21:51:53.321526 [Node 1] event=SEND_REQUEST to=3 lamport=538
2026/11/08 21:51:53.322103 [Node 1] event=RECEIVE_REPLY from=3 lamport=541 local=542 replies=1/2
2026/11/08 21:51:53.799853 [Node 1] event=RECEIVE_REPLY from=2 lamport=542 local=543 replies=2/2
2026/11/08 21:51:53.799853 [Node 1] event=ENTER_CS lamport=543
2026/11/08 21:51:54.801010 [Node 1] event=EXIT_CS lamport=543 deferred=0
2026/11/08 21:51:55.014655 [Node 1] event=RECEIVE_REQUEST from=3 lamport=542 local=544 state=RELEASED
2026/11/08 21:51:55.014655 [Node 1] event=SEND_REPLY to=3 lamport=545
2026/11/08 21:51:55.663869 [Node 1] event=RECEIVE_REQUEST from=2 lamport=545 local=546 state=RELEASED
2026/11/08 21:51:55.663869 [Node 1] event=SEND_REPLY to=2 lamport=547
2026/11/08 21:51:57.997063 [Node 1] event=REQUEST_CS lamport=548
2026/11/08 21:51:57.997063 [Node 1] event=SEND_REQUEST to=2 lamport=548
2026/11/08 21:51:57.997063 [Node 1] event=SEND_REQUEST to=3 lamport=548
2026/11/08 21:51:57.997651 [Node 1] event=RECEIVE_REPLY from=3 lamport=551 local=552 replies=1/2
2026/11/08 21:51:57.997651 [Node 1] event=RECEIVE_REPLY from=2 lamport=552 local=553 replies=2/2
2026/11/08 21:51:57.997651 [Node 1] event=ENTER_CS lamport=553
2026/11/08 21:51:58.998705 [Node 1] event=EXIT_CS lamport=553 deferred=0
2026/11/08 21:51:59.362568 [Node 1] event=RECEIVE_REQUEST from=2 lamport=553 local=554 state=RELEASED
2026/11/08 21:51:59.362568 [Node 1] event=SEND_REPLY to=2 lamport=555
2026/11/08 21:51:59.921350 [Node 1] event=RECEIVE_REQUEST from=3 lamport=556 local=557 state=RELEASED
2026/11/08 21:51:59.921350 [Node 1] event=SEND_REPLY to=3 lamport=558
2026/11/08 21:52:00.527440 [Node 1] event=REQUEST_CS lamport=559
2026/11/08 21:52:00.527440 [Node 1] event=SEND_REQUEST to=2 lamport=559
2026/11/08 21:52:00.527440 [Node 1] event=SEND_REQUEST to=3 lamport=559
2026/11/08 21:52:00.527440 [Node 1] event=RECEIVE_REPLY from=2 lamport=561 local=562 replies=1/2
2026/11/08 21:52:01.365439 [Node 1] event=RECEIVE_REPLY from=3 lamport=562 local=563 replies=2/2
2026/11/08 21:52:01.365439 [Node 1] event=ENTER_CS lamport=563
2026/11/08 21:52:02.366988 [Node 1] event=EXIT_CS lamport=563 deferred=0
2026/11/08 21:52:03.382988 [Node 1] event=RECEIVE_REQUEST from=3 lamport=563 local=564 state=RELEASED
2026/11/08 21:52:03.382988 [Node 1] event=SEND_REPLY to=3 lamport=565
2026/11/08 21:52:03.769988 [Node 1] event=RECEIVE_REQUEST from=2 lamport=566 local=567 state=RELEASED
2026/11/08 21:52:03.769988 [Node 1] event=SEND_REPLY to=2 lamport=568
2026/11/08 21:52:05.629229 [Node 1] event=REQUEST_CS lamport=569
2026/11/08 21:52:05.629229 [Node 1] event=SEND_REQUEST to=2 lamport=569
2026/11/08 21:52:05.629229 [Node 1] event=SEND_REQUEST to=3 lamport=569
2026/11/08 21:52:05.629229 [Node 1] event=RECEIVE_REPLY from=2 lamport=572 local=573 replies=1/2
2026/11/08 21:52:05.629229 [Node 1] event=RECEIVE_REPLY from=3 lamport=571 local=574 replies=2/2
2026/11/08 21:52:05.629229 [Node 1] event=ENTER_CS lamport=574
2026/11/08 21:52:06.630727 [Node 1] event=EXIT_CS lamport=574 deferred=0
2026/11/08 21:52:07.949735 [Node 1] event=RECEIVE_REQUEST from=3 lamport=572 local=575 state=RELEASED
2026/11/08 21:52:07.949735 [Node 1] event=SEND_REPLY to=3 lamport=576
2026/11/08 21:52:08.476735 [Node 1] event=RECEIVE_REQUEST from=2 lamport=575 local=577 state=RELEASED
2026/11/08 21:52:08.476735 [Node 1] event=SEND_REPLY to=2 lamport=578
2026/11/08 21:52:08.530735 [Node 1] event=REQUEST_CS lamport=579
2026/11/08 21:52:08.530735 [Node 1] event=SEND_REQUEST to=2 lamport=579
2026/11/08 21:52:08.530735 [Node 1] event=SEND_REQUEST to=3 lamport=579
2026/11/08 21:52:08.951735 [Node 1] event=RECEIVE_REPLY from=3 lamport=581 local=582 replies=1/2
2026/11/08 21:52:09.952272 [Node 1] event=RECEIVE_REPLY from=2 lamport=584 local=585 replies=2/2
2026/11/08 21:52:09.952272 [Node 1] event=ENTER_CS lamport=585
2026/11/08 21:52:10.953269 [Node 1] event=EXIT_CS lamport=585 deferred=0
2026/11/08 21:52:11.223268 [Node 1] event=RECEIVE_REQUEST from=2 lamport=585 local=586 state=RELEASED
2026/11/08 21:52:11.223268 [Node 1] event=SEND_REPLY to=2 lamport=587
2026/11/08 21:52:12.342578 [Node 1] event=REQUEST_CS lamport=588
2026/11/08 21:52:12.342578 [Node 1] event=SEND_REQUEST to=2 lamport=588
2026/11/08 21:52:12.342578 [Node 1] event=SEND_REQUEST to=3 lamport=588
2026/11/08 21:52:12.343802 [Node 1] event=RECEIVE_REPLY from=3 lamport=590 local=591 replies=1/2
2026/11/08 21:52:12.343802 [Node 1] event=RECEIVE_REPLY from=2 lamport=591 local=592 replies=2/2
2026/11/08 21:52:12.343802 [Node 1] event=ENTER_CS lamport=592
2026/11/08 21:52:12.927316 [Node 1] event=RECEIVE_REQUEST from=3 lamport=591 local=593 state=HELD
2026/11/08 21:52:12.927316 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:52:13.344254 [Node 1] event=EXIT_CS lamport=593 deferred=1
2026/11/08 21:52:13.344254 [Node 1] event=SEND_REPLY to=3 lamport=594
2026/11/08 21:52:15.474377 [Node 1] event=REQUEST_CS lamport=595
2026/11/08 21:52:15.474377 [Node 1] event=SEND_REQUEST to=2 lamport=595
2026/11/08 21:52:15.474377 [Node 1] event=SEND_REQUEST to=3 lamport=595
2026/11/08 21:52:15.474377 [Node 1] event=RECEIVE_REPLY from=3 lamport=597 local=598 replies=1/2
2026/11/08 21:52:15.475379 [Node 1] event=RECEIVE_REPLY from=2 lamport=597 local=599 replies=2/2
2026/11/08 21:52:15.475379 [Node 1] event=ENTER_CS lamport=599
2026/11/08 21:52:15.783378 [Node 1] event=RECEIVE_REQUEST from=2 lamport=598 local=600 state=HELD
2026/11/08 21:52:15.783378 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:52:16.241378 [Node 1] event=RECEIVE_REQUEST from=3 lamport=601 local=602 state=HELD
2026/11/08 21:52:16.241378 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:52:16.476377 [Node 1] event=EXIT_CS lamport=602 deferred=2
2026/11/08 21:52:16.476377 [Node 1] event=SEND_REPLY to=3 lamport=603
2026/11/08 21:52:16.476377 [Node 1] event=SEND_REPLY to=2 lamport=604
2026/11/08 21:52:18.596911 [Node 1] event=REQUEST_CS lamport=605
2026/11/08 21:52:18.596911 [Node 1] event=SEND_REQUEST to=2 lamport=605
2026/11/08 21:52:18.596911 [Node 1] event=SEND_REQUEST to=3 lamport=605
2026/11/08 21:52:18.596911 [Node 1] event=RECEIVE_REPLY from=3 lamport=609 local=610 replies=1/2
2026/11/08 21:52:18.596911 [Node 1] event=RECEIVE_REPLY from=2 lamport=608 local=611 replies=2/2
2026/11/08 21:52:18.596911 [Node 1] event=ENTER_CS lamport=611
2026/11/08 21:52:19.597912 [Node 1] event=EXIT_CS lamport=611 deferred=0
2026/11/08 21:52:19.843921 [Node 1] event=RECEIVE_REQUEST from=2 lamport=609 local=612 state=RELEASED
2026/11/08 21:52:19.843921 [Node 1] event=SEND_REPLY to=2 lamport=613
2026/11/08 21:52:21.510580 [Node 1] event=RECEIVE_REQUEST from=3 lamport=612 local=614 state=RELEASED
2026/11/08 21:52:21.510580 [Node 1] event=SEND_REPLY to=3 lamport=615
2026/11/08 21:52:22.226676 [Node 1] event=REQUEST_CS lamport=616
2026/11/08 21:52:22.226676 [Node 1] event=SEND_REQUEST to=3 lamport=616
2026/11/08 21:52:22.226676 [Node 1] event=SEND_REQUEST to=2 lamport=616
2026/11/08 21:52:22.227947 [Node 1] event=RECEIVE_REPLY from=2 lamport=619 local=620 replies=1/2
2026/11/08 21:52:22.511464 [Node 1] event=RECEIVE_REPLY from=3 lamport=620 local=621 replies=2/2
2026/11/08 21:52:22.511464 [Node 1] event=ENTER_CS lamport=621
2026/11/08 21:52:23.303831 [Node 1] event=RECEIVE_REQUEST from=2 lamport=620 local=622 state=HELD
2026/11/08 21:52:23.303831 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:52:23.511831 [Node 1] event=EXIT_CS lamport=622 deferred=1
2026/11/08 21:52:23.511831 [Node 1] event=SEND_REPLY to=2 lamport=623
2026/11/08 21:52:25.432832 [Node 1] event=REQUEST_CS lamport=624
2026/11/08 21:52:25.432832 [Node 1] event=SEND_REQUEST to=2 lamport=624
2026/11/08 21:52:25.432832 [Node 1] event=SEND_REQUEST to=3 lamport=624
2026/11/08 21:52:25.433837 [Node 1] event=RECEIVE_REPLY from=3 lamport=626 local=627 replies=1/2
2026/11/08 21:52:25.433837 [Node 1] event=RECEIVE_REPLY from=2 lamport=626 local=628 replies=2/2
2026/11/08 21:52:25.433837 [Node 1] event=ENTER_CS lamport=628
2026/11/08 21:52:26.044832 [Node 1] event=RECEIVE_REQUEST from=3 lamport=627 local=629 state=HELD
2026/11/08 21:52:26.044832 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:52:26.434831 [Node 1] event=EXIT_CS lamport=629 deferred=1
2026/11/08 21:52:26.434831 [Node 1] event=SEND_REPLY to=3 lamport=630
2026/11/08 21:52:28.171637 [Node 1] event=RECEIVE_REQUEST from=2 lamport=630 local=631 state=RELEASED
2026/11/08 21:52:28.171637 [Node 1] event=SEND_REPLY to=2 lamport=632
2026/11/08 21:52:30.430493 [Node 1] event=REQUEST_CS lamport=633
2026/11/08 21:52:30.430493 [Node 1] event=SEND_REQUEST to=2 lamport=633
2026/11/08 21:52:30.430493 [Node 1] event=SEND_REQUEST to=3 lamport=633
2026/11/08 21:52:30.431494 [Node 1] event=RECEIVE_REPLY from=3 lamport=635 local=636 replies=1/2
2026/11/08 21:52:30.431494 [Node 1] event=RECEIVE_REPLY from=2 lamport=636 local=637 replies=2/2
2026/11/08 21:52:30.431494 [Node 1] event=ENTER_CS lamport=637
2026/11/08 21:52:30.678492 [Node 1] event=RECEIVE_REQUEST from=3 lamport=636 local=638 state=HELD
2026/11/08 21:52:30.678492 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:52:31.432133 [Node 1] event=EXIT_CS lamport=638 deferred=1
2026/11/08 21:52:31.432133 [Node 1] event=SEND_REPLY to=3 lamport=639
2026/11/08 21:52:32.514293 [Node 1] event=RECEIVE_REQUEST from=2 lamport=639 local=640 state=RELEASED
2026/11/08 21:52:32.514293 [Node 1] event=SEND_REPLY to=2 lamport=641
2026/11/08 21:52:33.019374 [Node 1] event=REQUEST_CS lamport=642
2026/11/08 21:52:33.019374 [Node 1] event=SEND_REQUEST to=2 lamport=642
2026/11/08 21:52:33.019374 [Node 1] event=SEND_REQUEST to=3 lamport=642
2026/11/08 21:52:33.020286 [Node 1] event=RECEIVE_REPLY from=3 lamport=644 local=645 replies=1/2
2026/11/08 21:52:33.515859 [Node 1] event=RECEIVE_REPLY from=2 lamport=646 local=647 replies=2/2
2026/11/08 21:52:33.515859 [Node 1] event=ENTER_CS lamport=647
2026/11/08 21:52:33.893090 [Node 1] event=RECEIVE_REQUEST from=3 lamport=645 local=648 state=HELD
2026/11/08 21:52:33.893113 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:52:34.517144 [Node 1] event=EXIT_CS lamport=648 deferred=1
2026/11/08 21:52:34.517144 [Node 1] event=SEND_REPLY to=3 lamport=649
2026/11/08 21:52:35.552890 [Node 1] event=RECEIVE_REQUEST from=2 lamport=649 local=650 state=RELEASED
2026/11/08 21:52:35.552890 [Node 1] event=SEND_REPLY to=2 lamport=651
2026/11/08 21:52:36.322708 [Node 1] event=REQUEST_CS lamport=652
2026/11/08 21:52:36.322708 [Node 1] event=SEND_REQUEST to=2 lamport=652
2026/11/08 21:52:36.322708 [Node 1] event=SEND_REQUEST to=3 lamport=652
2026/11/08 21:52:36.323292 [Node 1] event=RECEIVE_REPLY from=3 lamport=654 local=655 replies=1/2
2026/11/08 21:52:36.554402 [Node 1] event=RECEIVE_REPLY from=2 lamport=656 local=657 replies=2/2
2026/11/08 21:52:36.554402 [Node 1] event=ENTER_CS lamport=657
2026/11/08 21:52:37.555304 [Node 1] event=EXIT_CS lamport=657 deferred=0
2026/11/08 21:52:38.828896 [Node 1] event=RECEIVE_REQUEST from=3 lamport=655 local=658 state=RELEASED
2026/11/08 21:52:38.828896 [Node 1] event=SEND_REPLY to=3 lamport=659
2026/11/08 21:52:39.296317 [Node 1] event=REQUEST_CS lamport=660
2026/11/08 21:52:39.296317 [Node 1] event=SEND_REQUEST to=2 lamport=660
2026/11/08 21:52:39.296317 [Node 1] event=SEND_REQUEST to=3 lamport=660
2026/11/08 21:52:39.296944 [Node 1] event=RECEIVE_REPLY from=2 lamport=662 local=663 replies=1/2
2026/11/08 21:52:39.830804 [Node 1] event=RECEIVE_REPLY from=3 lamport=662 local=664 replies=2/2
2026/11/08 21:52:39.830804 [Node 1] event=ENTER_CS lamport=664
2026/11/08 21:52:40.133315 [Node 1] event=RECEIVE_REQUEST from=2 lamport=663 local=665 state=HELD
2026/11/08 21:52:40.133315 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:52:40.831666 [Node 1] event=EXIT_CS lamport=665 deferred=1
2026/11/08 21:52:40.831666 [Node 1] event=SEND_REPLY to=2 lamport=666
2026/11/08 21:52:43.336352 [Node 1] event=RECEIVE_REQUEST from=2 lamport=668 local=669 state=RELEASED
2026/11/08 21:52:43.336876 [Node 1] event=SEND_REPLY to=2 lamport=670
2026/11/08 21:52:43.813437 [Node 1] event=RECEIVE_REQUEST from=3 lamport=671 local=672 state=RELEASED
2026/11/08 21:52:43.813437 [Node 1] event=SEND_REPLY to=3 lamport=673
2026/11/08 21:52:43.825579 [Node 1] event=REQUEST_CS lamport=674
2026/11/08 21:52:43.825579 [Node 1] event=SEND_REQUEST to=2 lamport=674
2026/11/08 21:52:43.825579 [Node 1] event=SEND_REQUEST to=3 lamport=674
2026/11/08 21:52:44.338659 [Node 1] event=RECEIVE_REPLY from=2 lamport=676 local=677 replies=1/2
2026/11/08 21:52:45.339623 [Node 1] event=RECEIVE_REPLY from=3 lamport=679 local=680 replies=2/2
2026/11/08 21:52:45.339623 [Node 1] event=ENTER_CS lamport=680
2026/11/08 21:52:46.340765 [Node 1] event=EXIT_CS lamport=680 deferred=0
2026/11/08 21:52:47.784360 [Node 1] event=RECEIVE_REQUEST from=2 lamport=678 local=681 state=RELEASED
2026/11/08 21:52:47.784360 [Node 1] event=SEND_REPLY to=2 lamport=682
2026/11/08 21:52:48.388782 [Node 1] event=REQUEST_CS lamport=683
2026/11/08 21:52:48.388782 [Node 1] event=SEND_REQUEST to=3 lamport=683
2026/11/08 21:52:48.388782 [Node 1] event=SEND_REQUEST to=2 lamport=683
2026/11/08 21:52:48.389434 [Node 1] event=RECEIVE_REPLY from=3 lamport=685 local=686 replies=1/2
2026/11/08 21:52:48.402490 [Node 1] event=RECEIVE_REQUEST from=3 lamport=686 local=687 state=WANTED
2026/11/08 21:52:48.402490 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:52:48.785579 [Node 1] event=RECEIVE_REPLY from=2 lamport=689 local=690 replies=2/2
2026/11/08 21:52:48.785579 [Node 1] event=ENTER_CS lamport=690
2026/11/08 21:52:49.785985 [Node 1] event=EXIT_CS lamport=690 deferred=1
2026/11/08 21:52:49.785985 [Node 1] event=SEND_REPLY to=3 lamport=691
2026/11/08 21:52:51.133723 [Node 1] event=RECEIVE_REQUEST from=2 lamport=690 local=692 state=RELEASED
2026/11/08 21:52:51.133723 [Node 1] event=SEND_REPLY to=2 lamport=693
2026/11/08 21:52:52.198911 [Node 1] event=RECEIVE_REQUEST from=3 lamport=695 local=696 state=RELEASED
2026/11/08 21:52:52.198917 [Node 1] event=SEND_REPLY to=3 lamport=697
2026/11/08 21:52:52.881008 [Node 1] event=REQUEST_CS lamport=698
2026/11/08 21:52:52.881008 [Node 1] event=SEND_REQUEST to=2 lamport=698
2026/11/08 21:52:52.881008 [Node 1] event=SEND_REQUEST to=3 lamport=698
2026/11/08 21:52:52.881576 [Node 1] event=RECEIVE_REPLY from=2 lamport=700 local=701 replies=1/2
2026/11/08 21:52:53.199962 [Node 1] event=RECEIVE_REPLY from=3 lamport=701 local=702 replies=2/2
2026/11/08 21:52:53.199962 [Node 1] event=ENTER_CS lamport=702
2026/11/08 21:52:54.200216 [Node 1] event=EXIT_CS lamport=702 deferred=0
2026/11/08 21:52:55.466536 [Node 1] event=RECEIVE_REQUEST from=2 lamport=701 local=703 state=RELEASED
2026/11/08 21:52:55.466536 [Node 1] event=SEND_REPLY to=2 lamport=704
2026/11/08 21:52:56.993302 [Node 1] event=RECEIVE_REQUEST from=3 lamport=704 local=705 state=RELEASED
2026/11/08 21:52:56.993909 [Node 1] event=SEND_REPLY to=3 lamport=706
2026/11/08 21:52:57.829182 [Node 1] event=REQUEST_CS lamport=707
2026/11/08 21:52:57.829182 [Node 1] event=SEND_REQUEST to=2 lamport=707
2026/11/08 21:52:57.829182 [Node 1] event=SEND_REQUEST to=3 lamport=707
2026/11/08 21:52:57.829766 [Node 1] event=RECEIVE_REPLY from=2 lamport=709 local=710 replies=1/2
2026/11/08 21:52:57.994282 [Node 1] event=RECEIVE_REPLY from=3 lamport=711 local=712 replies=2/2
2026/11/08 21:52:57.994282 [Node 1] event=ENTER_CS lamport=712
2026/11/08 21:52:58.995440 [Node 1] event=EXIT_CS lamport=712 deferred=0
2026/11/08 21:53:00.099526 [Node 1] event=REQUEST_CS lamport=713
2026/11/08 21:53:00.099526 [Node 1] event=SEND_REQUEST to=2 lamport=713
2026/11/08 21:53:00.099526 [Node 1] event=SEND_REQUEST to=3 lamport=713
2026/11/08 21:53:00.100180 [Node 1] event=RECEIVE_REPLY from=3 lamport=715 local=716 replies=1/2
2026/11/08 21:53:00.100180 [Node 1] event=RECEIVE_REPLY from=2 lamport=715 local=717 replies=2/2
2026/11/08 21:53:00.100180 [Node 1] event=ENTER_CS lamport=717
2026/11/08 21:53:00.170818 [Node 1] event=RECEIVE_REQUEST from=3 lamport=716 local=718 state=HELD
2026/11/08 21:53:00.170818 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:53:00.358929 [Node 1] event=RECEIVE_REQUEST from=2 lamport=719 local=720 state=HELD
2026/11/08 21:53:00.358929 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:53:01.100212 [Node 1] event=EXIT_CS lamport=720 deferred=2
2026/11/08 21:53:01.100212 [Node 1] event=SEND_REPLY to=2 lamport=721
2026/11/08 21:53:01.100212 [Node 1] event=SEND_REPLY to=3 lamport=722
2026/11/08 21:53:03.404129 [Node 1] event=REQUEST_CS lamport=723
2026/11/08 21:53:03.404129 [Node 1] event=SEND_REQUEST to=2 lamport=723
2026/11/08 21:53:03.404129 [Node 1] event=SEND_REQUEST to=3 lamport=723
2026/11/08 21:53:03.404792 [Node 1] event=RECEIVE_REPLY from=3 lamport=726 local=727 replies=1/2
2026/11/08 21:53:03.404792 [Node 1] event=RECEIVE_REPLY from=2 lamport=727 local=728 replies=2/2
2026/11/08 21:53:03.404792 [Node 1] event=ENTER_CS lamport=728
2026/11/08 21:53:04.405151 [Node 1] event=EXIT_CS lamport=728 deferred=0
2026/11/08 21:53:05.356513 [Node 1] event=RECEIVE_REQUEST from=3 lamport=727 local=729 state=RELEASED
2026/11/08 21:53:05.356513 [Node 1] event=SEND_REPLY to=3 lamport=730
2026/11/08 21:53:05.943536 [Node 1] event=RECEIVE_REQUEST from=2 lamport=730 local=731 state=RELEASED
2026/11/08 21:53:05.944077 [Node 1] event=SEND_REPLY to=2 lamport=732
2026/11/08 21:53:07.933741 [Node 1] event=REQUEST_CS lamport=733
2026/11/08 21:53:07.933741 [Node 1] event=SEND_REQUEST to=2 lamport=733
2026/11/08 21:53:07.933741 [Node 1] event=SEND_REQUEST to=3 lamport=733
2026/11/08 21:53:07.934340 [Node 1] event=RECEIVE_REPLY from=3 lamport=735 local=736 replies=1/2
2026/11/08 21:53:07.934865 [Node 1] event=RECEIVE_REPLY from=2 lamport=736 local=737 replies=2/2
2026/11/08 21:53:07.934865 [Node 1] event=ENTER_CS lamport=737
2026/11/08 21:53:08.935691 [Node 1] event=EXIT_CS lamport=737 deferred=0
2026/11/08 21:53:09.378339 [Node 1] event=RECEIVE_REQUEST from=3 lamport=736 local=738 state=RELEASED
2026/11/08 21:53:09.378339 [Node 1] event=SEND_REPLY to=3 lamport=739
2026/11/08 21:53:09.676010 [Node 1] event=RECEIVE_REQUEST from=2 lamport=739 local=740 state=RELEASED
2026/11/08 21:53:09.676010 [Node 1] event=SEND_REPLY to=2 lamport=741
2026/11/08 21:53:11.570234 [Node 1] event=RECEIVE_REQUEST from=3 lamport=744 local=745 state=RELEASED
2026/11/08 21:53:11.570234 [Node 1] event=SEND_REPLY to=3 lamport=746
2026/11/08 21:53:12.308074 [Node 1] event=REQUEST_CS lamport=747
2026/11/08 21:53:12.308074 [Node 1] event=SEND_REQUEST to=3 lamport=747
2026/11/08 21:53:12.308074 [Node 1] event=SEND_REQUEST to=2 lamport=747
2026/11/08 21:53:12.308696 [Node 1] event=RECEIVE_REPLY from=2 lamport=749 local=750 replies=1/2
2026/11/08 21:53:12.571332 [Node 1] event=RECEIVE_REPLY from=3 lamport=750 local=751 replies=2/2
2026/11/08 21:53:12.571332 [Node 1] event=ENTER_CS lamport=751
2026/11/08 21:53:13.571528 [Node 1] event=EXIT_CS lamport=751 deferred=0
2026/11/08 21:53:14.241356 [Node 1] event=RECEIVE_REQUEST from=3 lamport=751 local=752 state=RELEASED
2026/11/08 21:53:14.241356 [Node 1] event=SEND_REPLY to=3 lamport=753
2026/11/08 21:53:14.648115 [Node 1] event=RECEIVE_REQUEST from=2 lamport=754 local=755 state=RELEASED
2026/11/08 21:53:14.648115 [Node 1] event=SEND_REPLY to=2 lamport=756
2026/11/08 21:53:16.295348 [Node 1] event=REQUEST_CS lamport=757
2026/11/08 21:53:16.295397 [Node 1] event=SEND_REQUEST to=2 lamport=757
2026/11/08 21:53:16.295397 [Node 1] event=SEND_REQUEST to=3 lamport=757
2026/11/08 21:53:16.295925 [Node 1] event=RECEIVE_REPLY from=3 lamport=759 local=760 replies=1/2
2026/11/08 21:53:16.295925 [Node 1] event=RECEIVE_REPLY from=2 lamport=760 local=761 replies=2/2
2026/11/08 21:53:16.295925 [Node 1] event=ENTER_CS lamport=761
2026/11/08 21:53:17.296787 [Node 1] event=EXIT_CS lamport=761 deferred=0
2026/11/08 21:53:17.622912 [Node 1] event=RECEIVE_REQUEST from=2 lamport=761 local=762 state=RELEASED
2026/11/08 21:53:17.622912 [Node 1] event=SEND_REPLY to=2 lamport=763
2026/11/08 21:53:18.840093 [Node 1] event=REQUEST_CS lamport=764
2026/11/08 21:53:18.840093 [Node 1] event=SEND_REQUEST to=2 lamport=764
2026/11/08 21:53:18.840093 [Node 1] event=SEND_REQUEST to=3 lamport=764
2026/11/08 21:53:18.840723 [Node 1] event=RECEIVE_REPLY from=2 lamport=767 local=768 replies=1/2
2026/11/08 21:53:18.840723 [Node 1] event=RECEIVE_REPLY from=3 lamport=766 local=769 replies=2/2
2026/11/08 21:53:18.841228 [Node 1] event=ENTER_CS lamport=769
2026/11/08 21:53:18.872813 [Node 1] event=RECEIVE_REQUEST from=3 lamport=767 local=770 state=HELD
2026/11/08 21:53:18.872813 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:53:19.841708 [Node 1] event=EXIT_CS lamport=770 deferred=1
2026/11/08 21:53:19.841708 [Node 1] event=SEND_REPLY to=3 lamport=771
2026/11/08 21:53:21.981175 [Node 1] event=RECEIVE_REQUEST from=2 lamport=770 local=772 state=RELEASED
2026/11/08 21:53:21.981175 [Node 1] event=SEND_REPLY to=2 lamport=773
2026/11/08 21:53:23.042352 [Node 1] event=RECEIVE_REQUEST from=3 lamport=775 local=776 state=RELEASED
2026/11/08 21:53:23.042352 [Node 1] event=SEND_REPLY to=3 lamport=777
2026/11/08 21:53:23.249318 [Node 1] event=REQUEST_CS lamport=778
2026/11/08 21:53:23.249318 [Node 1] event=SEND_REQUEST to=3 lamport=778
2026/11/08 21:53:23.249318 [Node 1] event=SEND_REQUEST to=2 lamport=778
2026/11/08 21:53:23.249990 [Node 1] event=RECEIVE_REPLY from=2 lamport=780 local=781 replies=1/2
2026/11/08 21:53:24.043619 [Node 1] event=RECEIVE_REPLY from=3 lamport=782 local=783 replies=2/2
2026/11/08 21:53:24.043619 [Node 1] event=ENTER_CS lamport=783
2026/11/08 21:53:25.044150 [Node 1] event=EXIT_CS lamport=783 deferred=0
2026/11/08 21:53:26.216287 [Node 1] event=REQUEST_CS lamport=784
2026/11/08 21:53:26.216325 [Node 1] event=SEND_REQUEST to=2 lamport=784
2026/11/08 21:53:26.216325 [Node 1] event=SEND_REQUEST to=3 lamport=784
2026/11/08 21:53:26.216897 [Node 1] event=RECEIVE_REPLY from=2 lamport=786 local=787 replies=1/2
2026/11/08 21:53:26.217449 [Node 1] event=RECEIVE_REPLY from=3 lamport=786 local=788 replies=2/2
2026/11/08 21:53:26.217449 [Node 1] event=ENTER_CS lamport=788
2026/11/08 21:53:26.397761 [Node 1] event=RECEIVE_REQUEST from=2 lamport=787 local=789 state=HELD
2026/11/08 21:53:26.397761 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:53:26.778923 [Node 1] event=RECEIVE_REQUEST from=3 lamport=790 local=791 state=HELD
2026/11/08 21:53:26.778923 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:53:27.218582 [Node 1] event=EXIT_CS lamport=791 deferred=2
2026/11/08 21:53:27.218626 [Node 1] event=SEND_REPLY to=3 lamport=792
2026/11/08 21:53:27.218626 [Node 1] event=SEND_REPLY to=2 lamport=793
2026/11/08 21:53:30.136271 [Node 1] event=REQUEST_CS lamport=794
2026/11/08 21:53:30.136332 [Node 1] event=SEND_REQUEST to=2 lamport=794
2026/11/08 21:53:30.136332 [Node 1] event=SEND_REQUEST to=3 lamport=794
2026/11/08 21:53:30.136863 [Node 1] event=RECEIVE_REPLY from=3 lamport=798 local=799 replies=1/2
2026/11/08 21:53:30.136863 [Node 1] event=RECEIVE_REPLY from=2 lamport=797 local=800 replies=2/2
2026/11/08 21:53:30.136863 [Node 1] event=ENTER_CS lamport=800
2026/11/08 21:53:31.137394 [Node 1] event=EXIT_CS lamport=800 deferred=0
2026/11/08 21:53:31.461770 [Node 1] event=RECEIVE_REQUEST from=3 lamport=799 local=801 state=RELEASED
2026/11/08 21:53:31.461770 [Node 1] event=SEND_REPLY to=3 lamport=802
2026/11/08 21:53:31.714403 [Node 1] event=RECEIVE_REQUEST from=2 lamport=802 local=803 state=RELEASED
2026/11/08 21:53:31.714403 [Node 1] event=SEND_REPLY to=2 lamport=804
2026/11/08 21:53:33.627070 [Node 1] event=RECEIVE_REQUEST from=3 lamport=806 local=807 state=RELEASED
2026/11/08 21:53:33.627070 [Node 1] event=SEND_REPLY to=3 lamport=808
2026/11/08 21:53:34.468272 [Node 1] event=RECEIVE_REQUEST from=2 lamport=809 local=810 state=RELEASED
2026/11/08 21:53:34.468272 [Node 1] event=SEND_REPLY to=2 lamport=811
2026/11/08 21:53:34.639947 [Node 1] event=REQUEST_CS lamport=812
2026/11/08 21:53:34.639947 [Node 1] event=SEND_REQUEST to=2 lamport=812
2026/11/08 21:53:34.639947 [Node 1] event=SEND_REQUEST to=3 lamport=812
2026/11/08 21:53:34.640546 [Node 1] event=RECEIVE_REPLY from=3 lamport=814 local=815 replies=1/2
2026/11/08 21:53:35.628233 [Node 1] event=RECEIVE_REPLY from=2 lamport=815 local=816 replies=2/2
2026/11/08 21:53:35.628233 [Node 1] event=ENTER_CS lamport=816
2026/11/08 21:53:36.629232 [Node 1] event=EXIT_CS lamport=816 deferred=0
2026/11/08 21:53:37.050969 [Node 1] event=RECEIVE_REQUEST from=3 lamport=815 local=817 state=RELEASED
2026/11/08 21:53:37.050969 [Node 1] event=SEND_REPLY to=3 lamport=818
2026/11/08 21:53:39.039847 [Node 1] event=RECEIVE_REQUEST from=2 lamport=818 local=819 state=RELEASED
2026/11/08 21:53:39.039847 [Node 1] event=SEND_REPLY to=2 lamport=820
2026/11/08 21:53:39.296597 [Node 1] event=REQUEST_CS lamport=821
2026/11/08 21:53:39.296639 [Node 1] event=SEND_REQUEST to=2 lamport=821
2026/11/08 21:53:39.296639 [Node 1] event=SEND_REQUEST to=3 lamport=821
2026/11/08 21:53:39.297161 [Node 1] event=RECEIVE_REPLY from=3 lamport=823 local=824 replies=1/2
2026/11/08 21:53:40.040607 [Node 1] event=RECEIVE_REPLY from=2 lamport=825 local=826 replies=2/2
2026/11/08 21:53:40.041151 [Node 1] event=ENTER_CS lamport=826
2026/11/08 21:53:40.936519 [Node 1] event=RECEIVE_REQUEST from=3 lamport=824 local=827 state=HELD
2026/11/08 21:53:40.936519 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:53:41.042185 [Node 1] event=EXIT_CS lamport=827 deferred=1
2026/11/08 21:53:41.042232 [Node 1] event=SEND_REPLY to=3 lamport=828
2026/11/08 21:53:42.692904 [Node 1] event=REQUEST_CS lamport=829
2026/11/08 21:53:42.692966 [Node 1] event=SEND_REQUEST to=2 lamport=829
2026/11/08 21:53:42.692966 [Node 1] event=SEND_REQUEST to=3 lamport=829
2026/11/08 21:53:42.693500 [Node 1] event=RECEIVE_REPLY from=2 lamport=831 local=832 replies=1/2
2026/11/08 21:53:42.693500 [Node 1] event=RECEIVE_REPLY from=3 lamport=831 local=833 replies=2/2
2026/11/08 21:53:42.693500 [Node 1] event=ENTER_CS lamport=833
2026/11/08 21:53:43.366391 [Node 1] event=RECEIVE_REQUEST from=2 lamport=832 local=834 state=HELD
2026/11/08 21:53:43.366391 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:53:43.532502 [Node 1] event=RECEIVE_REQUEST from=3 lamport=835 local=836 state=HELD
2026/11/08 21:53:43.532502 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:53:43.693733 [Node 1] event=EXIT_CS lamport=836 deferred=2
2026/11/08 21:53:43.693733 [Node 1] event=SEND_REPLY to=3 lamport=837
2026/11/08 21:53:43.693733 [Node 1] event=SEND_REPLY to=2 lamport=838
2026/11/08 21:53:46.442171 [Node 1] event=RECEIVE_REQUEST from=2 lamport=841 local=842 state=RELEASED
2026/11/08 21:53:46.442171 [Node 1] event=SEND_REPLY to=2 lamport=843
2026/11/08 21:53:47.115312 [Node 1] event=REQUEST_CS lamport=844
2026/11/08 21:53:47.115312 [Node 1] event=SEND_REQUEST to=2 lamport=844
2026/11/08 21:53:47.115312 [Node 1] event=SEND_REQUEST to=3 lamport=844
2026/11/08 21:53:47.116257 [Node 1] event=RECEIVE_REPLY from=3 lamport=846 local=847 replies=1/2
2026/11/08 21:53:47.443422 [Node 1] event=RECEIVE_REPLY from=2 lamport=847 local=848 replies=2/2
2026/11/08 21:53:47.443422 [Node 1] event=ENTER_CS lamport=848
2026/11/08 21:53:47.555583 [Node 1] event=RECEIVE_REQUEST from=3 lamport=847 local=849 state=HELD
2026/11/08 21:53:47.555583 [Node 1] event=DEFER_REPLY to=3
2026/11/08 21:53:48.444398 [Node 1] event=EXIT_CS lamport=849 deferred=1
2026/11/08 21:53:48.444398 [Node 1] event=SEND_REPLY to=3 lamport=850
2026/11/08 21:53:48.917249 [Node 1] event=RECEIVE_REQUEST from=2 lamport=850 local=851 state=RELEASED
2026/11/08 21:53:48.917804 [Node 1] event=SEND_REPLY to=2 lamport=852
2026/11/08 21:53:51.800334 [Node 1] event=RECEIVE_REQUEST from=3 lamport=854 local=855 state=RELEASED
2026/11/08 21:53:51.800334 [Node 1] event=SEND_REPLY to=3 lamport=856
2026/11/08 21:53:52.246978 [Node 1] event=REQUEST_CS lamport=857
2026/11/08 21:53:52.246978 [Node 1] event=SEND_REQUEST to=2 lamport=857
2026/11/08 21:53:52.246978 [Node 1] event=SEND_REQUEST to=3 lamport=857
2026/11/08 21:53:52.247597 [Node 1] event=RECEIVE_REPLY from=2 lamport=859 local=860 replies=1/2
2026/11/08 21:53:52.801544 [Node 1] event=RECEIVE_REPLY from=3 lamport=860 local=861 replies=2/2
2026/11/08 21:53:52.801544 [Node 1] event=ENTER_CS lamport=861
2026/11/08 21:53:52.802711 [Node 1] event=RECEIVE_REQUEST from=2 lamport=860 local=862 state=HELD
2026/11/08 21:53:52.802711 [Node 1] event=DEFER_REPLY to=2
2026/11/08 21:53:53.801942 [Node 1] event=EXIT_CS lamport=862 deferred=1
2026/11/08 21:53:53.801942 [Node 1] event=SEND_REPLY to=2 lamport=863
2026/11/08 21:53:55.579993 [Node 1] event=RECEIVE_REQUEST from=3 lamport=863 local=864 state=RELEASED
2026/11/08 21:53:55.579993 [Node 1] event=SEND_REPLY to=3 lamport=865
2026/11/08 21:53:55.904672 [Node 1] event=RECEIVE_REQUEST from=2 lamport=867 local=868 state=RELEASED
2026/11/08 21:53:55.904672 [Node 1] event=SEND_REPLY to=2 lamport=869
2026/11/08 21:53:56.490202 [Node 1] event=REQUEST_CS lamport=870
2026/11/08 21:53:56.490202 [Node 1] event=SEND_REQUEST to=2 lamport=870
2026/11/08 21:53:56.490202 [Node 1] event=SEND_REQUEST to=3 lamport=870
2026/11/08 21:53:58.142969 [Node 1] event=RETRY to=2 error="rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:51285->127.0.0.1:5002: wsarecv: An existing connection was forcibly closed by the remote host."

```

### Node 2

```
2026/11/08 21:48:28.583766 [Node 2] event=STARTUP address=localhost:5002 peers=2
2026/11/08 21:48:29.792411 [Node 2] event=REQUEST_CS lamport=1
2026/11/08 21:48:29.792411 [Node 2] event=SEND_REQUEST to=1 lamport=1
2026/11/08 21:48:29.792411 [Node 2] event=SEND_REQUEST to=3 lamport=1
2026/11/08 21:48:29.800412 [Node 2] event=RECEIVE_REPLY from=1 lamport=3 local=4 replies=1/2
2026/11/08 21:48:29.896411 [Node 2] event=RECEIVE_REQUEST from=1 lamport=4 local=5 state=WANTED
2026/11/08 21:48:29.896411 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:48:31.800801 [Node 2] event=RECEIVE_REPLY from=3 lamport=3 local=6 replies=2/2
2026/11/08 21:48:31.800801 [Node 2] event=ENTER_CS lamport=6
2026/11/08 21:48:32.800921 [Node 2] event=EXIT_CS lamport=6 deferred=1
2026/11/08 21:48:32.800921 [Node 2] event=SEND_REPLY to=1 lamport=7
2026/11/08 21:48:33.485550 [Node 2] event=RECEIVE_REQUEST from=3 lamport=7 local=8 state=RELEASED
2026/11/08 21:48:33.485550 [Node 2] event=SEND_REPLY to=3 lamport=9
2026/11/08 21:48:33.946419 [Node 2] event=REQUEST_CS lamport=10
2026/11/08 21:48:33.946452 [Node 2] event=SEND_REQUEST to=1 lamport=10
2026/11/08 21:48:33.946452 [Node 2] event=SEND_REQUEST to=3 lamport=10
2026/11/08 21:48:33.946979 [Node 2] event=RECEIVE_REPLY from=1 lamport=12 local=13 replies=1/2
2026/11/08 21:48:34.804571 [Node 2] event=RECEIVE_REPLY from=3 lamport=13 local=14 replies=2/2
2026/11/08 21:48:34.804571 [Node 2] event=ENTER_CS lamport=14
2026/11/08 21:48:35.152282 [Node 2] event=RECEIVE_REQUEST from=1 lamport=13 local=15 state=HELD
2026/11/08 21:48:35.152835 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:48:35.804863 [Node 2] event=EXIT_CS lamport=15 deferred=1
2026/11/08 21:48:35.804863 [Node 2] event=SEND_REPLY to=1 lamport=16
2026/11/08 21:48:36.758110 [Node 2] event=RECEIVE_REQUEST from=3 lamport=16 local=17 state=RELEASED
2026/11/08 21:48:36.758632 [Node 2] event=SEND_REPLY to=3 lamport=18
2026/11/08 21:48:38.236784 [Node 2] event=REQUEST_CS lamport=19
2026/11/08 21:48:38.236784 [Node 2] event=SEND_REQUEST to=1 lamport=19
2026/11/08 21:48:38.236784 [Node 2] event=SEND_REQUEST to=3 lamport=19
2026/11/08 21:48:38.237402 [Node 2] event=RECEIVE_REPLY from=3 lamport=22 local=23 replies=1/2
2026/11/08 21:48:38.237402 [Node 2] event=RECEIVE_REPLY from=1 lamport=21 local=24 replies=2/2
2026/11/08 21:48:38.237402 [Node 2] event=ENTER_CS lamport=24
2026/11/08 21:48:39.237739 [Node 2] event=EXIT_CS lamport=24 deferred=0
2026/11/08 21:48:39.961400 [Node 2] event=RECEIVE_REQUEST from=1 lamport=22 local=25 state=RELEASED
2026/11/08 21:48:39.961400 [Node 2] event=SEND_REPLY to=1 lamport=26
2026/11/08 21:48:40.054925 [Node 2] event=RECEIVE_REQUEST from=3 lamport=25 local=27 state=RELEASED
2026/11/08 21:48:40.054925 [Node 2] event=SEND_REPLY to=3 lamport=28
2026/11/08 21:48:42.037145 [Node 2] event=REQUEST_CS lamport=29
2026/11/08 21:48:42.037145 [Node 2] event=SEND_REQUEST to=1 lamport=29
2026/11/08 21:48:42.037145 [Node 2] event=SEND_REQUEST to=3 lamport=29
2026/11/08 21:48:42.038147 [Node 2] event=RECEIVE_REPLY from=1 lamport=31 local=32 replies=1/2
2026/11/08 21:48:42.038147 [Node 2] event=RECEIVE_REPLY from=3 lamport=32 local=33 replies=2/2
2026/11/08 21:48:42.038147 [Node 2] event=ENTER_CS lamport=33
2026/11/08 21:48:43.038652 [Node 2] event=EXIT_CS lamport=33 deferred=0
2026/11/08 21:48:43.447654 [Node 2] event=RECEIVE_REQUEST from=3 lamport=33 local=34 state=RELEASED
2026/11/08 21:48:43.447654 [Node 2] event=SEND_REPLY to=3 lamport=35
2026/11/08 21:48:43.970199 [Node 2] event=RECEIVE_REQUEST from=1 lamport=36 local=37 state=RELEASED
2026/11/08 21:48:43.970199 [Node 2] event=SEND_REPLY to=1 lamport=38
2026/11/08 21:48:45.740213 [Node 2] event=REQUEST_CS lamport=39
2026/11/08 21:48:45.740213 [Node 2] event=SEND_REQUEST to=1 lamport=39
2026/11/08 21:48:45.740213 [Node 2] event=SEND_REQUEST to=3 lamport=39
2026/11/08 21:48:45.740213 [Node 2] event=RECEIVE_REPLY from=1 lamport=42 local=43 replies=1/2
2026/11/08 21:48:45.741213 [Node 2] event=RECEIVE_REPLY from=3 lamport=41 local=44 replies=2/2
2026/11/08 21:48:45.741213 [Node 2] event=ENTER_CS lamport=44
2026/11/08 21:48:45.859942 [Node 2] event=RECEIVE_REQUEST from=3 lamport=42 local=45 state=HELD
2026/11/08 21:48:45.859942 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:48:46.741925 [Node 2] event=EXIT_CS lamport=45 deferred=1
2026/11/08 21:48:46.741925 [Node 2] event=SEND_REPLY to=3 lamport=46
2026/11/08 21:48:47.095362 [Node 2] event=RECEIVE_REQUEST from=1 lamport=45 local=47 state=RELEASED
2026/11/08 21:48:47.095362 [Node 2] event=SEND_REPLY to=1 lamport=48
2026/11/08 21:48:48.213898 [Node 2] event=REQUEST_CS lamport=49
2026/11/08 21:48:48.213898 [Node 2] event=SEND_REQUEST to=1 lamport=49
2026/11/08 21:48:48.213898 [Node 2] event=SEND_REQUEST to=3 lamport=49
2026/11/08 21:48:48.213898 [Node 2] event=RECEIVE_REPLY from=3 lamport=51 local=52 replies=1/2
2026/11/08 21:48:48.744171 [Node 2] event=RECEIVE_REPLY from=1 lamport=52 local=53 replies=2/2
2026/11/08 21:48:48.744171 [Node 2] event=ENTER_CS lamport=53
2026/11/08 21:48:49.449148 [Node 2] event=RECEIVE_REQUEST from=3 lamport=52 local=54 state=HELD
2026/11/08 21:48:49.449148 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:48:49.745332 [Node 2] event=EXIT_CS lamport=54 deferred=1
2026/11/08 21:48:49.745332 [Node 2] event=SEND_REPLY to=3 lamport=55
2026/11/08 21:48:51.229880 [Node 2] event=RECEIVE_REQUEST from=1 lamport=55 local=56 state=RELEASED
2026/11/08 21:48:51.229880 [Node 2] event=SEND_REPLY to=1 lamport=57
2026/11/08 21:48:52.287606 [Node 2] event=RECEIVE_REQUEST from=3 lamport=59 local=60 state=RELEASED
2026/11/08 21:48:52.287606 [Node 2] event=SEND_REPLY to=3 lamport=61
2026/11/08 21:48:53.076607 [Node 2] event=REQUEST_CS lamport=62
2026/11/08 21:48:53.076607 [Node 2] event=SEND_REQUEST to=1 lamport=62
2026/11/08 21:48:53.076607 [Node 2] event=SEND_REQUEST to=3 lamport=62
2026/11/08 21:48:53.076607 [Node 2] event=RECEIVE_REPLY from=1 lamport=64 local=65 replies=1/2
2026/11/08 21:48:53.288674 [Node 2] event=RECEIVE_REPLY from=3 lamport=65 local=66 replies=2/2
2026/11/08 21:48:53.288674 [Node 2] event=ENTER_CS lamport=66
2026/11/08 21:48:53.866702 [Node 2] event=RECEIVE_REQUEST from=1 lamport=65 local=67 state=HELD
2026/11/08 21:48:53.866702 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:48:54.290144 [Node 2] event=EXIT_CS lamport=67 deferred=1
2026/11/08 21:48:54.290144 [Node 2] event=SEND_REPLY to=1 lamport=68
2026/11/08 21:48:55.402492 [Node 2] event=RECEIVE_REQUEST from=3 lamport=68 local=69 state=RELEASED
2026/11/08 21:48:55.402492 [Node 2] event=SEND_REPLY to=3 lamport=70
2026/11/08 21:48:57.588557 [Node 2] event=RECEIVE_REQUEST from=1 lamport=72 local=73 state=RELEASED
2026/11/08 21:48:57.588557 [Node 2] event=SEND_REPLY to=1 lamport=74
2026/11/08 21:48:58.099310 [Node 2] event=REQUEST_CS lamport=75
2026/11/08 21:48:58.099310 [Node 2] event=SEND_REQUEST to=3 lamport=75
2026/11/08 21:48:58.099310 [Node 2] event=SEND_REQUEST to=1 lamport=75
2026/11/08 21:48:58.099310 [Node 2] event=RECEIVE_REPLY from=3 lamport=77 local=78 replies=1/2
2026/11/08 21:48:58.589816 [Node 2] event=RECEIVE_REPLY from=1 lamport=79 local=80 replies=2/2
2026/11/08 21:48:58.589816 [Node 2] event=ENTER_CS lamport=80
2026/11/08 21:48:59.590894 [Node 2] event=EXIT_CS lamport=80 deferred=0
2026/11/08 21:48:59.650893 [Node 2] event=RECEIVE_REQUEST from=3 lamport=78 local=81 state=RELEASED
2026/11/08 21:48:59.650893 [Node 2] event=SEND_REPLY to=3 lamport=82
2026/11/08 21:48:59.682893 [Node 2] event=RECEIVE_REQUEST from=1 lamport=82 local=83 state=RELEASED
2026/11/08 21:48:59.682893 [Node 2] event=SEND_REPLY to=1 lamport=84
2026/11/08 21:49:01.614867 [Node 2] event=REQUEST_CS lamport=85
2026/11/08 21:49:01.614867 [Node 2] event=SEND_REQUEST to=3 lamport=85
2026/11/08 21:49:01.614867 [Node 2] event=SEND_REQUEST to=1 lamport=85
2026/11/08 21:49:01.615753 [Node 2] event=RECEIVE_REPLY from=3 lamport=88 local=89 replies=1/2
2026/11/08 21:49:01.652335 [Node 2] event=RECEIVE_REPLY from=1 lamport=89 local=90 replies=2/2
2026/11/08 21:49:01.652335 [Node 2] event=ENTER_CS lamport=90
2026/11/08 21:49:02.652978 [Node 2] event=EXIT_CS lamport=90 deferred=0
2026/11/08 21:49:04.278450 [Node 2] event=RECEIVE_REQUEST from=3 lamport=89 local=91 state=RELEASED
2026/11/08 21:49:04.278450 [Node 2] event=SEND_REPLY to=3 lamport=92
2026/11/08 21:49:04.588451 [Node 2] event=RECEIVE_REQUEST from=1 lamport=92 local=93 state=RELEASED
2026/11/08 21:49:04.588451 [Node 2] event=SEND_REPLY to=1 lamport=94
2026/11/08 21:49:05.786416 [Node 2] event=REQUEST_CS lamport=95
2026/11/08 21:49:05.786484 [Node 2] event=SEND_REQUEST to=3 lamport=95
2026/11/08 21:49:05.786484 [Node 2] event=SEND_REQUEST to=1 lamport=95
2026/11/08 21:49:05.787007 [Node 2] event=RECEIVE_REPLY from=3 lamport=97 local=98 replies=1/2
2026/11/08 21:49:06.280857 [Node 2] event=RECEIVE_REPLY from=1 lamport=98 local=99 replies=2/2
2026/11/08 21:49:06.280857 [Node 2] event=ENTER_CS lamport=99
2026/11/08 21:49:06.726118 [Node 2] event=RECEIVE_REQUEST from=3 lamport=98 local=100 state=HELD
2026/11/08 21:49:06.726118 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:49:07.281441 [Node 2] event=EXIT_CS lamport=100 deferred=1
2026/11/08 21:49:07.281441 [Node 2] event=SEND_REPLY to=3 lamport=101
2026/11/08 21:49:08.570309 [Node 2] event=RECEIVE_REQUEST from=1 lamport=101 local=102 state=RELEASED
2026/11/08 21:49:08.570835 [Node 2] event=SEND_REPLY to=1 lamport=103
2026/11/08 21:49:09.828145 [Node 2] event=RECEIVE_REQUEST from=3 lamport=105 local=106 state=RELEASED
2026/11/08 21:49:09.828743 [Node 2] event=SEND_REPLY to=3 lamport=107
2026/11/08 21:49:10.941783 [Node 2] event=REQUEST_CS lamport=108
2026/11/08 21:49:10.941783 [Node 2] event=SEND_REQUEST to=1 lamport=108
2026/11/08 21:49:10.941783 [Node 2] event=SEND_REQUEST to=3 lamport=108
2026/11/08 21:49:10.941783 [Node 2] event=RECEIVE_REPLY from=3 lamport=111 local=112 replies=1/2
2026/11/08 21:49:10.941783 [Node 2] event=RECEIVE_REPLY from=1 lamport=110 local=113 replies=2/2
2026/11/08 21:49:10.941783 [Node 2] event=ENTER_CS lamport=113
2026/11/08 21:49:11.700783 [Node 2] event=RECEIVE_REQUEST from=1 lamport=111 local=114 state=HELD
2026/11/08 21:49:11.700783 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:49:11.943782 [Node 2] event=EXIT_CS lamport=114 deferred=1
2026/11/08 21:49:11.943782 [Node 2] event=SEND_REPLY to=1 lamport=115
2026/11/08 21:49:13.886785 [Node 2] event=RECEIVE_REQUEST from=3 lamport=114 local=116 state=RELEASED
2026/11/08 21:49:13.886785 [Node 2] event=SEND_REPLY to=3 lamport=117
2026/11/08 21:49:14.228215 [Node 2] event=REQUEST_CS lamport=118
2026/11/08 21:49:14.228215 [Node 2] event=SEND_REQUEST to=1 lamport=118
2026/11/08 21:49:14.228215 [Node 2] event=SEND_REQUEST to=3 lamport=118
2026/11/08 21:49:14.228215 [Node 2] event=RECEIVE_REPLY from=1 lamport=120 local=121 replies=1/2
2026/11/08 21:49:14.473214 [Node 2] event=RECEIVE_REQUEST from=1 lamport=121 local=122 state=WANTED
2026/11/08 21:49:14.473214 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:49:14.889214 [Node 2] event=RECEIVE_REPLY from=3 lamport=124 local=125 replies=2/2
2026/11/08 21:49:14.889214 [Node 2] event=ENTER_CS lamport=125
2026/11/08 21:49:15.891001 [Node 2] event=EXIT_CS lamport=125 deferred=1
2026/11/08 21:49:15.891001 [Node 2] event=SEND_REPLY to=1 lamport=126
2026/11/08 21:49:17.160999 [Node 2] event=RECEIVE_REQUEST from=3 lamport=125 local=127 state=RELEASED
2026/11/08 21:49:17.161578 [Node 2] event=SEND_REPLY to=3 lamport=128
2026/11/08 21:49:19.579330 [Node 2] event=RECEIVE_REQUEST from=1 lamport=130 local=131 state=RELEASED
2026/11/08 21:49:19.579330 [Node 2] event=SEND_REPLY to=1 lamport=132
2026/11/08 21:49:19.772539 [Node 2] event=REQUEST_CS lamport=133
2026/11/08 21:49:19.772539 [Node 2] event=SEND_REQUEST to=3 lamport=133
2026/11/08 21:49:19.772539 [Node 2] event=SEND_REQUEST to=1 lamport=133
2026/11/08 21:49:19.772539 [Node 2] event=RECEIVE_REPLY from=3 lamport=135 local=136 replies=1/2
2026/11/08 21:49:20.580955 [Node 2] event=RECEIVE_REPLY from=1 lamport=136 local=137 replies=2/2
2026/11/08 21:49:20.580955 [Node 2] event=ENTER_CS lamport=137
2026/11/08 21:49:20.948917 [Node 2] event=RECEIVE_REQUEST from=3 lamport=136 local=138 state=HELD
2026/11/08 21:49:20.948917 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:49:21.582207 [Node 2] event=EXIT_CS lamport=138 deferred=1
2026/11/08 21:49:21.582207 [Node 2] event=SEND_REPLY to=3 lamport=139
2026/11/08 21:49:23.396838 [Node 2] event=RECEIVE_REQUEST from=1 lamport=139 local=140 state=RELEASED
2026/11/08 21:49:23.397442 [Node 2] event=SEND_REPLY to=1 lamport=141
2026/11/08 21:49:23.883087 [Node 2] event=RECEIVE_REQUEST from=3 lamport=143 local=144 state=RELEASED
2026/11/08 21:49:23.883087 [Node 2] event=SEND_REPLY to=3 lamport=145
2026/11/08 21:49:23.897182 [Node 2] event=REQUEST_CS lamport=146
2026/11/08 21:49:23.897211 [Node 2] event=SEND_REQUEST to=1 lamport=146
2026/11/08 21:49:23.897211 [Node 2] event=SEND_REQUEST to=3 lamport=146
2026/11/08 21:49:24.398741 [Node 2] event=RECEIVE_REPLY from=1 lamport=148 local=149 replies=1/2
2026/11/08 21:49:25.399080 [Node 2] event=RECEIVE_REPLY from=3 lamport=151 local=152 replies=2/2
2026/11/08 21:49:25.399080 [Node 2] event=ENTER_CS lamport=152
2026/11/08 21:49:26.400061 [Node 2] event=EXIT_CS lamport=152 deferred=0
2026/11/08 21:49:27.834612 [Node 2] event=RECEIVE_REQUEST from=3 lamport=152 local=153 state=RELEASED
2026/11/08 21:49:27.834612 [Node 2] event=SEND_REPLY to=3 lamport=154
2026/11/08 21:49:28.170386 [Node 2] event=RECEIVE_REQUEST from=1 lamport=155 local=156 state=RELEASED
2026/11/08 21:49:28.170386 [Node 2] event=SEND_REPLY to=1 lamport=157
2026/11/08 21:49:29.309702 [Node 2] event=REQUEST_CS lamport=158
2026/11/08 21:49:29.309702 [Node 2] event=SEND_REQUEST to=3 lamport=158
2026/11/08 21:49:29.309702 [Node 2] event=SEND_REQUEST to=1 lamport=158
2026/11/08 21:49:29.309702 [Node 2] event=RECEIVE_REPLY from=3 lamport=160 local=161 replies=1/2
2026/11/08 21:49:29.837702 [Node 2] event=RECEIVE_REPLY from=1 lamport=161 local=162 replies=2/2
2026/11/08 21:49:29.837702 [Node 2] event=ENTER_CS lamport=162
2026/11/08 21:49:30.839250 [Node 2] event=EXIT_CS lamport=162 deferred=0
2026/11/08 21:49:31.660819 [Node 2] event=RECEIVE_REQUEST from=3 lamport=161 local=163 state=RELEASED
2026/11/08 21:49:31.660819 [Node 2] event=SEND_REPLY to=3 lamport=164
2026/11/08 21:49:33.609940 [Node 2] event=RECEIVE_REQUEST from=1 lamport=164 local=165 state=RELEASED
2026/11/08 21:49:33.609940 [Node 2] event=SEND_REPLY to=1 lamport=166
2026/11/08 21:49:34.316255 [Node 2] event=REQUEST_CS lamport=167
2026/11/08 21:49:34.316255 [Node 2] event=SEND_REQUEST to=1 lamport=167
2026/11/08 21:49:34.316255 [Node 2] event=SEND_REQUEST to=3 lamport=167
2026/11/08 21:49:34.317256 [Node 2] event=RECEIVE_REPLY from=3 lamport=169 local=170 replies=1/2
2026/11/08 21:49:34.611253 [Node 2] event=RECEIVE_REPLY from=1 lamport=170 local=171 replies=2/2
2026/11/08 21:49:34.611253 [Node 2] event=ENTER_CS lamport=171
2026/11/08 21:49:35.612780 [Node 2] event=EXIT_CS lamport=171 deferred=0
2026/11/08 21:49:35.650781 [Node 2] event=RECEIVE_REQUEST from=3 lamport=170 local=172 state=RELEASED
2026/11/08 21:49:35.650781 [Node 2] event=SEND_REPLY to=3 lamport=173
2026/11/08 21:49:36.405970 [Node 2] event=RECEIVE_REQUEST from=1 lamport=173 local=174 state=RELEASED
2026/11/08 21:49:36.405970 [Node 2] event=SEND_REPLY to=1 lamport=175
2026/11/08 21:49:38.729037 [Node 2] event=RECEIVE_REQUEST from=1 lamport=179 local=180 state=RELEASED
2026/11/08 21:49:38.729037 [Node 2] event=SEND_REPLY to=1 lamport=181
2026/11/08 21:49:39.227905 [Node 2] event=REQUEST_CS lamport=182
2026/11/08 21:49:39.227905 [Node 2] event=SEND_REQUEST to=3 lamport=182
2026/11/08 21:49:39.227905 [Node 2] event=SEND_REQUEST to=1 lamport=182
2026/11/08 21:49:39.229031 [Node 2] event=RECEIVE_REPLY from=3 lamport=184 local=185 replies=1/2
2026/11/08 21:49:39.729545 [Node 2] event=RECEIVE_REPLY from=1 lamport=185 local=186 replies=2/2
2026/11/08 21:49:39.729545 [Node 2] event=ENTER_CS lamport=186
2026/11/08 21:49:39.780758 [Node 2] event=RECEIVE_REQUEST from=3 lamport=185 local=187 state=HELD
2026/11/08 21:49:39.780758 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:49:40.730288 [Node 2] event=EXIT_CS lamport=187 deferred=1
2026/11/08 21:49:40.730288 [Node 2] event=SEND_REPLY to=3 lamport=188
2026/11/08 21:49:42.308091 [Node 2] event=RECEIVE_REQUEST from=1 lamport=188 local=189 state=RELEASED
2026/11/08 21:49:42.308091 [Node 2] event=SEND_REPLY to=1 lamport=190
2026/11/08 21:49:43.605729 [Node 2] event=REQUEST_CS lamport=191
2026/11/08 21:49:43.605729 [Node 2] event=SEND_REQUEST to=1 lamport=191
2026/11/08 21:49:43.605729 [Node 2] event=SEND_REQUEST to=3 lamport=191
2026/11/08 21:49:43.606730 [Node 2] event=RECEIVE_REPLY from=1 lamport=194 local=195 replies=1/2
2026/11/08 21:49:43.606730 [Node 2] event=RECEIVE_REPLY from=3 lamport=193 local=196 replies=2/2
2026/11/08 21:49:43.606730 [Node 2] event=ENTER_CS lamport=196
2026/11/08 21:49:44.261728 [Node 2] event=RECEIVE_REQUEST from=3 lamport=194 local=197 state=HELD
2026/11/08 21:49:44.261728 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:49:44.607729 [Node 2] event=EXIT_CS lamport=197 deferred=1
2026/11/08 21:49:44.607729 [Node 2] event=SEND_REPLY to=3 lamport=198
2026/11/08 21:49:44.914259 [Node 2] event=RECEIVE_REQUEST from=1 lamport=197 local=199 state=RELEASED
2026/11/08 21:49:44.914259 [Node 2] event=SEND_REPLY to=1 lamport=200
2026/11/08 21:49:46.695382 [Node 2] event=RECEIVE_REQUEST from=3 lamport=202 local=203 state=RELEASED
2026/11/08 21:49:46.695382 [Node 2] event=SEND_REPLY to=3 lamport=204
2026/11/08 21:49:47.402941 [Node 2] event=REQUEST_CS lamport=205
2026/11/08 21:49:47.402941 [Node 2] event=SEND_REQUEST to=1 lamport=205
2026/11/08 21:49:47.402941 [Node 2] event=SEND_REQUEST to=3 lamport=205
2026/11/08 21:49:47.403553 [Node 2] event=RECEIVE_REPLY from=1 lamport=207 local=208 replies=1/2
2026/11/08 21:49:47.696634 [Node 2] event=RECEIVE_REPLY from=3 lamport=208 local=209 replies=2/2
2026/11/08 21:49:47.696634 [Node 2] event=ENTER_CS lamport=209
2026/11/08 21:49:47.912668 [Node 2] event=RECEIVE_REQUEST from=1 lamport=208 local=210 state=HELD
2026/11/08 21:49:47.912668 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:49:48.696934 [Node 2] event=EXIT_CS lamport=210 deferred=1
2026/11/08 21:49:48.696981 [Node 2] event=SEND_REPLY to=1 lamport=211
2026/11/08 21:49:50.583648 [Node 2] event=RECEIVE_REQUEST from=3 lamport=211 local=212 state=RELEASED
2026/11/08 21:49:50.583648 [Node 2] event=SEND_REPLY to=3 lamport=213
2026/11/08 21:49:51.438738 [Node 2] event=REQUEST_CS lamport=214
2026/11/08 21:49:51.438799 [Node 2] event=SEND_REQUEST to=1 lamport=214
2026/11/08 21:49:51.438799 [Node 2] event=SEND_REQUEST to=3 lamport=214
2026/11/08 21:49:51.439324 [Node 2] event=RECEIVE_REPLY from=1 lamport=216 local=217 replies=1/2
2026/11/08 21:49:51.585678 [Node 2] event=RECEIVE_REPLY from=3 lamport=218 local=219 replies=2/2
2026/11/08 21:49:51.585678 [Node 2] event=ENTER_CS lamport=219
2026/11/08 21:49:52.155451 [Node 2] event=RECEIVE_REQUEST from=1 lamport=217 local=220 state=HELD
2026/11/08 21:49:52.155451 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:49:52.586453 [Node 2] event=EXIT_CS lamport=220 deferred=1
2026/11/08 21:49:52.586453 [Node 2] event=SEND_REPLY to=1 lamport=221
2026/11/08 21:49:54.621989 [Node 2] event=REQUEST_CS lamport=222
2026/11/08 21:49:54.621989 [Node 2] event=SEND_REQUEST to=1 lamport=222
2026/11/08 21:49:54.621989 [Node 2] event=SEND_REQUEST to=3 lamport=222
2026/11/08 21:49:54.622990 [Node 2] event=RECEIVE_REPLY from=3 lamport=224 local=225 replies=1/2
2026/11/08 21:49:54.622990 [Node 2] event=RECEIVE_REPLY from=1 lamport=224 local=226 replies=2/2
2026/11/08 21:49:54.622990 [Node 2] event=ENTER_CS lamport=226
2026/11/08 21:49:54.935990 [Node 2] event=RECEIVE_REQUEST from=3 lamport=225 local=227 state=HELD
2026/11/08 21:49:54.935990 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:49:55.623576 [Node 2] event=EXIT_CS lamport=227 deferred=1
2026/11/08 21:49:55.623576 [Node 2] event=SEND_REPLY to=3 lamport=228
2026/11/08 21:49:57.227981 [Node 2] event=RECEIVE_REQUEST from=1 lamport=228 local=229 state=RELEASED
2026/11/08 21:49:57.227981 [Node 2] event=SEND_REPLY to=1 lamport=230
2026/11/08 21:49:58.135041 [Node 2] event=REQUEST_CS lamport=231
2026/11/08 21:49:58.135041 [Node 2] event=SEND_REQUEST to=1 lamport=231
2026/11/08 21:49:58.135041 [Node 2] event=SEND_REQUEST to=3 lamport=231
2026/11/08 21:49:58.135814 [Node 2] event=RECEIVE_REPLY from=3 lamport=233 local=234 replies=1/2
2026/11/08 21:49:58.229905 [Node 2] event=RECEIVE_REPLY from=1 lamport=235 local=236 replies=2/2
2026/11/08 21:49:58.229905 [Node 2] event=ENTER_CS lamport=236
2026/11/08 21:49:59.231551 [Node 2] event=EXIT_CS lamport=236 deferred=0
2026/11/08 21:49:59.739551 [Node 2] event=RECEIVE_REQUEST from=3 lamport=234 local=237 state=RELEASED
2026/11/08 21:49:59.739551 [Node 2] event=SEND_REPLY to=3 lamport=238
2026/11/08 21:50:01.475551 [Node 2] event=RECEIVE_REQUEST from=1 lamport=238 local=239 state=RELEASED
2026/11/08 21:50:01.475551 [Node 2] event=SEND_REPLY to=1 lamport=240
2026/11/08 21:50:02.384094 [Node 2] event=RECEIVE_REQUEST from=3 lamport=243 local=244 state=RELEASED
2026/11/08 21:50:02.384094 [Node 2] event=SEND_REPLY to=3 lamport=245
2026/11/08 21:50:02.756389 [Node 2] event=REQUEST_CS lamport=246
2026/11/08 21:50:02.756389 [Node 2] event=SEND_REQUEST to=1 lamport=246
2026/11/08 21:50:02.756389 [Node 2] event=SEND_REQUEST to=3 lamport=246
2026/11/08 21:50:02.757183 [Node 2] event=RECEIVE_REPLY from=1 lamport=248 local=249 replies=1/2
2026/11/08 21:50:03.478307 [Node 2] event=RECEIVE_REPLY from=3 lamport=249 local=250 replies=2/2
2026/11/08 21:50:03.478307 [Node 2] event=ENTER_CS lamport=250
2026/11/08 21:50:03.561384 [Node 2] event=RECEIVE_REQUEST from=1 lamport=249 local=251 state=HELD
2026/11/08 21:50:03.561384 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:04.478686 [Node 2] event=EXIT_CS lamport=251 deferred=1
2026/11/08 21:50:04.478686 [Node 2] event=SEND_REPLY to=1 lamport=252
2026/11/08 21:50:06.191544 [Node 2] event=RECEIVE_REQUEST from=3 lamport=252 local=253 state=RELEASED
2026/11/08 21:50:06.191544 [Node 2] event=SEND_REPLY to=3 lamport=254
2026/11/08 21:50:07.034858 [Node 2] event=REQUEST_CS lamport=255
2026/11/08 21:50:07.034858 [Node 2] event=SEND_REQUEST to=1 lamport=255
2026/11/08 21:50:07.034858 [Node 2] event=SEND_REQUEST to=3 lamport=255
2026/11/08 21:50:07.034858 [Node 2] event=RECEIVE_REPLY from=1 lamport=257 local=258 replies=1/2
2026/11/08 21:50:07.192859 [Node 2] event=RECEIVE_REPLY from=3 lamport=258 local=259 replies=2/2
2026/11/08 21:50:07.192859 [Node 2] event=ENTER_CS lamport=259
2026/11/08 21:50:08.193858 [Node 2] event=EXIT_CS lamport=259 deferred=0
2026/11/08 21:50:09.472858 [Node 2] event=RECEIVE_REQUEST from=1 lamport=258 local=260 state=RELEASED
2026/11/08 21:50:09.472858 [Node 2] event=SEND_REPLY to=1 lamport=261
2026/11/08 21:50:11.005396 [Node 2] event=RECEIVE_REQUEST from=3 lamport=261 local=262 state=RELEASED
2026/11/08 21:50:11.005396 [Node 2] event=SEND_REPLY to=3 lamport=263
2026/11/08 21:50:11.893394 [Node 2] event=REQUEST_CS lamport=264
2026/11/08 21:50:11.893394 [Node 2] event=SEND_REQUEST to=1 lamport=264
2026/11/08 21:50:11.893394 [Node 2] event=SEND_REQUEST to=3 lamport=264
2026/11/08 21:50:11.893394 [Node 2] event=RECEIVE_REPLY from=1 lamport=267 local=268 replies=1/2
2026/11/08 21:50:12.007941 [Node 2] event=RECEIVE_REPLY from=3 lamport=269 local=270 replies=2/2
2026/11/08 21:50:12.007941 [Node 2] event=ENTER_CS lamport=270
2026/11/08 21:50:13.008938 [Node 2] event=EXIT_CS lamport=270 deferred=0
2026/11/08 21:50:13.579938 [Node 2] event=RECEIVE_REQUEST from=3 lamport=270 local=271 state=RELEASED
2026/11/08 21:50:13.579938 [Node 2] event=SEND_REPLY to=3 lamport=272
2026/11/08 21:50:14.056938 [Node 2] event=RECEIVE_REQUEST from=1 lamport=273 local=274 state=RELEASED
2026/11/08 21:50:14.056938 [Node 2] event=SEND_REPLY to=1 lamport=275
2026/11/08 21:50:15.910069 [Node 2] event=REQUEST_CS lamport=276
2026/11/08 21:50:15.910069 [Node 2] event=SEND_REQUEST to=3 lamport=276
2026/11/08 21:50:15.910069 [Node 2] event=SEND_REQUEST to=1 lamport=276
2026/11/08 21:50:15.910069 [Node 2] event=RECEIVE_REPLY from=3 lamport=278 local=279 replies=1/2
2026/11/08 21:50:15.911070 [Node 2] event=RECEIVE_REPLY from=1 lamport=279 local=280 replies=2/2
2026/11/08 21:50:15.911070 [Node 2] event=ENTER_CS lamport=280
2026/11/08 21:50:16.652068 [Node 2] event=RECEIVE_REQUEST from=1 lamport=280 local=281 state=HELD
2026/11/08 21:50:16.652068 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:16.912069 [Node 2] event=EXIT_CS lamport=281 deferred=1
2026/11/08 21:50:16.912069 [Node 2] event=SEND_REPLY to=1 lamport=282
2026/11/08 21:50:18.016069 [Node 2] event=REQUEST_CS lamport=283
2026/11/08 21:50:18.016069 [Node 2] event=SEND_REQUEST to=1 lamport=283
2026/11/08 21:50:18.016069 [Node 2] event=SEND_REQUEST to=3 lamport=283
2026/11/08 21:50:18.016069 [Node 2] event=RECEIVE_REPLY from=1 lamport=286 local=287 replies=1/2
2026/11/08 21:50:18.016069 [Node 2] event=RECEIVE_REPLY from=3 lamport=285 local=288 replies=2/2
2026/11/08 21:50:18.016069 [Node 2] event=ENTER_CS lamport=288
2026/11/08 21:50:18.124070 [Node 2] event=RECEIVE_REQUEST from=3 lamport=286 local=289 state=HELD
2026/11/08 21:50:18.124070 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:50:19.018068 [Node 2] event=EXIT_CS lamport=289 deferred=1
2026/11/08 21:50:19.018068 [Node 2] event=SEND_REPLY to=3 lamport=290
2026/11/08 21:50:21.542646 [Node 2] event=RECEIVE_REQUEST from=1 lamport=289 local=291 state=RELEASED
2026/11/08 21:50:21.542646 [Node 2] event=SEND_REPLY to=1 lamport=292
2026/11/08 21:50:21.758495 [Node 2] event=RECEIVE_REQUEST from=3 lamport=294 local=295 state=RELEASED
2026/11/08 21:50:21.758495 [Node 2] event=SEND_REPLY to=3 lamport=296
2026/11/08 21:50:22.380813 [Node 2] event=REQUEST_CS lamport=297
2026/11/08 21:50:22.380813 [Node 2] event=SEND_REQUEST to=1 lamport=297
2026/11/08 21:50:22.380813 [Node 2] event=SEND_REQUEST to=3 lamport=297
2026/11/08 21:50:22.543813 [Node 2] event=RECEIVE_REPLY from=1 lamport=299 local=300 replies=1/2
2026/11/08 21:50:23.544858 [Node 2] event=RECEIVE_REPLY from=3 lamport=302 local=303 replies=2/2
2026/11/08 21:50:23.544858 [Node 2] event=ENTER_CS lamport=303
2026/11/08 21:50:23.562347 [Node 2] event=RECEIVE_REQUEST from=1 lamport=301 local=304 state=HELD
2026/11/08 21:50:23.562374 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:24.545549 [Node 2] event=EXIT_CS lamport=304 deferred=1
2026/11/08 21:50:24.545549 [Node 2] event=SEND_REPLY to=1 lamport=305
2026/11/08 21:50:26.408605 [Node 2] event=RECEIVE_REQUEST from=3 lamport=305 local=306 state=RELEASED
2026/11/08 21:50:26.408605 [Node 2] event=SEND_REPLY to=3 lamport=307
2026/11/08 21:50:26.701148 [Node 2] event=REQUEST_CS lamport=308
2026/11/08 21:50:26.701148 [Node 2] event=SEND_REQUEST to=1 lamport=308
2026/11/08 21:50:26.701148 [Node 2] event=SEND_REQUEST to=3 lamport=308
2026/11/08 21:50:26.701148 [Node 2] event=RECEIVE_REPLY from=1 lamport=310 local=311 replies=1/2
2026/11/08 21:50:27.409709 [Node 2] event=RECEIVE_REPLY from=3 lamport=311 local=312 replies=2/2
2026/11/08 21:50:27.409709 [Node 2] event=ENTER_CS lamport=312
2026/11/08 21:50:28.363709 [Node 2] event=RECEIVE_REQUEST from=1 lamport=311 local=313 state=HELD
2026/11/08 21:50:28.363709 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:28.410709 [Node 2] event=EXIT_CS lamport=313 deferred=1
2026/11/08 21:50:28.410709 [Node 2] event=SEND_REPLY to=1 lamport=314
2026/11/08 21:50:28.448709 [Node 2] event=RECEIVE_REQUEST from=3 lamport=314 local=315 state=RELEASED
2026/11/08 21:50:28.448709 [Node 2] event=SEND_REPLY to=3 lamport=316
2026/11/08 21:50:29.934347 [Node 2] event=REQUEST_CS lamport=317
2026/11/08 21:50:29.934347 [Node 2] event=SEND_REQUEST to=1 lamport=317
2026/11/08 21:50:29.934347 [Node 2] event=SEND_REQUEST to=3 lamport=317
2026/11/08 21:50:29.934347 [Node 2] event=RECEIVE_REPLY from=1 lamport=319 local=320 replies=1/2
2026/11/08 21:50:30.411525 [Node 2] event=RECEIVE_REPLY from=3 lamport=320 local=321 replies=2/2
2026/11/08 21:50:30.411525 [Node 2] event=ENTER_CS lamport=321
2026/11/08 21:50:30.714556 [Node 2] event=RECEIVE_REQUEST from=1 lamport=320 local=322 state=HELD
2026/11/08 21:50:30.714556 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:31.412735 [Node 2] event=EXIT_CS lamport=322 deferred=1
2026/11/08 21:50:31.412735 [Node 2] event=SEND_REPLY to=1 lamport=323
2026/11/08 21:50:32.501116 [Node 2] event=REQUEST_CS lamport=324
2026/11/08 21:50:32.501116 [Node 2] event=SEND_REQUEST to=1 lamport=324
2026/11/08 21:50:32.501116 [Node 2] event=SEND_REQUEST to=3 lamport=324
2026/11/08 21:50:32.501736 [Node 2] event=RECEIVE_REPLY from=3 lamport=326 local=327 replies=1/2
2026/11/08 21:50:32.501736 [Node 2] event=RECEIVE_REPLY from=1 lamport=326 local=328 replies=2/2
2026/11/08 21:50:32.501736 [Node 2] event=ENTER_CS lamport=328
2026/11/08 21:50:32.508253 [Node 2] event=RECEIVE_REQUEST from=3 lamport=327 local=329 state=HELD
2026/11/08 21:50:32.508253 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:50:33.502464 [Node 2] event=EXIT_CS lamport=329 deferred=1
2026/11/08 21:50:33.502464 [Node 2] event=SEND_REPLY to=3 lamport=330
2026/11/08 21:50:34.506439 [Node 2] event=RECEIVE_REQUEST from=1 lamport=330 local=331 state=RELEASED
2026/11/08 21:50:34.506439 [Node 2] event=SEND_REPLY to=1 lamport=332
2026/11/08 21:50:35.699017 [Node 2] event=RECEIVE_REQUEST from=3 lamport=334 local=335 state=RELEASED
2026/11/08 21:50:35.699593 [Node 2] event=SEND_REPLY to=3 lamport=336
2026/11/08 21:50:36.498857 [Node 2] event=REQUEST_CS lamport=337
2026/11/08 21:50:36.498921 [Node 2] event=SEND_REQUEST to=3 lamport=337
2026/11/08 21:50:36.498921 [Node 2] event=SEND_REQUEST to=1 lamport=337
2026/11/08 21:50:36.499443 [Node 2] event=RECEIVE_REPLY from=1 lamport=339 local=340 replies=1/2
2026/11/08 21:50:36.699697 [Node 2] event=RECEIVE_REPLY from=3 lamport=340 local=341 replies=2/2
2026/11/08 21:50:36.699697 [Node 2] event=ENTER_CS lamport=341
2026/11/08 21:50:36.931396 [Node 2] event=RECEIVE_REQUEST from=1 lamport=340 local=342 state=HELD
2026/11/08 21:50:36.931396 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:37.700572 [Node 2] event=EXIT_CS lamport=342 deferred=1
2026/11/08 21:50:37.700641 [Node 2] event=SEND_REPLY to=1 lamport=343
2026/11/08 21:50:39.572932 [Node 2] event=RECEIVE_REQUEST from=3 lamport=343 local=344 state=RELEASED
2026/11/08 21:50:39.573491 [Node 2] event=SEND_REPLY to=3 lamport=345
2026/11/08 21:50:39.882745 [Node 2] event=REQUEST_CS lamport=346
2026/11/08 21:50:39.882745 [Node 2] event=SEND_REQUEST to=1 lamport=346
2026/11/08 21:50:39.882745 [Node 2] event=SEND_REQUEST to=3 lamport=346
2026/11/08 21:50:39.883457 [Node 2] event=RECEIVE_REPLY from=1 lamport=348 local=349 replies=1/2
2026/11/08 21:50:40.574058 [Node 2] event=RECEIVE_REPLY from=3 lamport=349 local=350 replies=2/2
2026/11/08 21:50:40.574058 [Node 2] event=ENTER_CS lamport=350
2026/11/08 21:50:41.574279 [Node 2] event=EXIT_CS lamport=350 deferred=0
2026/11/08 21:50:41.881401 [Node 2] event=RECEIVE_REQUEST from=1 lamport=349 local=351 state=RELEASED
2026/11/08 21:50:41.881401 [Node 2] event=SEND_REPLY to=1 lamport=352
2026/11/08 21:50:43.829352 [Node 2] event=RECEIVE_REQUEST from=3 lamport=352 local=353 state=RELEASED
2026/11/08 21:50:43.829352 [Node 2] event=SEND_REPLY to=3 lamport=354
2026/11/08 21:50:44.588719 [Node 2] event=REQUEST_CS lamport=355
2026/11/08 21:50:44.588719 [Node 2] event=SEND_REQUEST to=1 lamport=355
2026/11/08 21:50:44.588719 [Node 2] event=SEND_REQUEST to=3 lamport=355
2026/11/08 21:50:44.589337 [Node 2] event=RECEIVE_REPLY from=1 lamport=358 local=359 replies=1/2
2026/11/08 21:50:44.830412 [Node 2] event=RECEIVE_REPLY from=3 lamport=359 local=360 replies=2/2
2026/11/08 21:50:44.830412 [Node 2] event=ENTER_CS lamport=360
2026/11/08 21:50:45.122762 [Node 2] event=RECEIVE_REQUEST from=1 lamport=359 local=361 state=HELD
2026/11/08 21:50:45.122762 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:45.831263 [Node 2] event=EXIT_CS lamport=361 deferred=1
2026/11/08 21:50:45.831263 [Node 2] event=SEND_REPLY to=1 lamport=362
2026/11/08 21:50:47.287417 [Node 2] event=RECEIVE_REQUEST from=3 lamport=362 local=363 state=RELEASED
2026/11/08 21:50:47.287988 [Node 2] event=SEND_REPLY to=3 lamport=364
2026/11/08 21:50:48.178033 [Node 2] event=REQUEST_CS lamport=365
2026/11/08 21:50:48.178033 [Node 2] event=SEND_REQUEST to=1 lamport=365
2026/11/08 21:50:48.178033 [Node 2] event=SEND_REQUEST to=3 lamport=365
2026/11/08 21:50:48.178574 [Node 2] event=RECEIVE_REPLY from=1 lamport=367 local=368 replies=1/2
2026/11/08 21:50:48.288714 [Node 2] event=RECEIVE_REPLY from=3 lamport=368 local=369 replies=2/2
2026/11/08 21:50:48.288714 [Node 2] event=ENTER_CS lamport=369
2026/11/08 21:50:49.050316 [Node 2] event=RECEIVE_REQUEST from=1 lamport=368 local=370 state=HELD
2026/11/08 21:50:49.050316 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:49.289439 [Node 2] event=EXIT_CS lamport=370 deferred=1
2026/11/08 21:50:49.289473 [Node 2] event=SEND_REPLY to=1 lamport=371
2026/11/08 21:50:49.849200 [Node 2] event=RECEIVE_REQUEST from=3 lamport=371 local=372 state=RELEASED
2026/11/08 21:50:49.849200 [Node 2] event=SEND_REPLY to=3 lamport=373
2026/11/08 21:50:50.317103 [Node 2] event=REQUEST_CS lamport=374
2026/11/08 21:50:50.317103 [Node 2] event=SEND_REQUEST to=1 lamport=374
2026/11/08 21:50:50.317103 [Node 2] event=SEND_REQUEST to=3 lamport=374
2026/11/08 21:50:50.317661 [Node 2] event=RECEIVE_REPLY from=1 lamport=376 local=377 replies=1/2
2026/11/08 21:50:51.289927 [Node 2] event=RECEIVE_REPLY from=3 lamport=377 local=378 replies=2/2
2026/11/08 21:50:51.289927 [Node 2] event=ENTER_CS lamport=378
2026/11/08 21:50:52.290754 [Node 2] event=EXIT_CS lamport=378 deferred=0
2026/11/08 21:50:52.644315 [Node 2] event=RECEIVE_REQUEST from=1 lamport=377 local=379 state=RELEASED
2026/11/08 21:50:52.644315 [Node 2] event=SEND_REPLY to=1 lamport=380
2026/11/08 21:50:53.374818 [Node 2] event=RECEIVE_REQUEST from=3 lamport=380 local=381 state=RELEASED
2026/11/08 21:50:53.374818 [Node 2] event=SEND_REPLY to=3 lamport=382
2026/11/08 21:50:53.621265 [Node 2] event=REQUEST_CS lamport=383
2026/11/08 21:50:53.621265 [Node 2] event=SEND_REQUEST to=1 lamport=383
2026/11/08 21:50:53.621265 [Node 2] event=SEND_REQUEST to=3 lamport=383
2026/11/08 21:50:53.645544 [Node 2] event=RECEIVE_REPLY from=1 lamport=385 local=386 replies=1/2
2026/11/08 21:50:54.646403 [Node 2] event=RECEIVE_REPLY from=3 lamport=388 local=389 replies=2/2
2026/11/08 21:50:54.646403 [Node 2] event=ENTER_CS lamport=389
2026/11/08 21:50:55.647587 [Node 2] event=EXIT_CS lamport=389 deferred=0
2026/11/08 21:50:55.950152 [Node 2] event=RECEIVE_REQUEST from=3 lamport=389 local=390 state=RELEASED
2026/11/08 21:50:55.950152 [Node 2] event=SEND_REPLY to=3 lamport=391
2026/11/08 21:50:57.272371 [Node 2] event=REQUEST_CS lamport=392
2026/11/08 21:50:57.272421 [Node 2] event=SEND_REQUEST to=1 lamport=392
2026/11/08 21:50:57.272421 [Node 2] event=SEND_REQUEST to=3 lamport=392
2026/11/08 21:50:57.272979 [Node 2] event=RECEIVE_REPLY from=1 lamport=394 local=395 replies=1/2
2026/11/08 21:50:57.273505 [Node 2] event=RECEIVE_REPLY from=3 lamport=395 local=396 replies=2/2
2026/11/08 21:50:57.273505 [Node 2] event=ENTER_CS lamport=396
2026/11/08 21:50:57.501653 [Node 2] event=RECEIVE_REQUEST from=1 lamport=395 local=397 state=HELD
2026/11/08 21:50:57.501653 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:50:58.274469 [Node 2] event=EXIT_CS lamport=397 deferred=1
2026/11/08 21:50:58.274523 [Node 2] event=SEND_REPLY to=1 lamport=398
2026/11/08 21:50:58.818627 [Node 2] event=RECEIVE_REQUEST from=3 lamport=398 local=399 state=RELEASED
2026/11/08 21:50:58.818627 [Node 2] event=SEND_REPLY to=3 lamport=400
2026/11/08 21:50:59.567455 [Node 2] event=REQUEST_CS lamport=401
2026/11/08 21:50:59.567455 [Node 2] event=SEND_REQUEST to=1 lamport=401
2026/11/08 21:50:59.567455 [Node 2] event=SEND_REQUEST to=3 lamport=401
2026/11/08 21:50:59.568023 [Node 2] event=RECEIVE_REPLY from=1 lamport=403 local=404 replies=1/2
2026/11/08 21:51:00.276933 [Node 2] event=RECEIVE_REPLY from=3 lamport=404 local=405 replies=2/2
2026/11/08 21:51:00.276933 [Node 2] event=ENTER_CS lamport=405
2026/11/08 21:51:01.278032 [Node 2] event=EXIT_CS lamport=405 deferred=0
2026/11/08 21:51:01.674312 [Node 2] event=RECEIVE_REQUEST from=3 lamport=405 local=406 state=RELEASED
2026/11/08 21:51:01.674312 [Node 2] event=SEND_REPLY to=3 lamport=407
2026/11/08 21:51:02.600009 [Node 2] event=RECEIVE_REQUEST from=1 lamport=408 local=409 state=RELEASED
2026/11/08 21:51:02.600009 [Node 2] event=SEND_REPLY to=1 lamport=410
2026/11/08 21:51:04.803489 [Node 2] event=RECEIVE_REQUEST from=3 lamport=412 local=413 state=RELEASED
2026/11/08 21:51:04.803489 [Node 2] event=SEND_REPLY to=3 lamport=414
2026/11/08 21:51:05.066814 [Node 2] event=REQUEST_CS lamport=415
2026/11/08 21:51:05.066814 [Node 2] event=SEND_REQUEST to=1 lamport=415
2026/11/08 21:51:05.066814 [Node 2] event=SEND_REQUEST to=3 lamport=415
2026/11/08 21:51:05.067411 [Node 2] event=RECEIVE_REPLY from=1 lamport=417 local=418 replies=1/2
2026/11/08 21:51:05.805104 [Node 2] event=RECEIVE_REPLY from=3 lamport=418 local=419 replies=2/2
2026/11/08 21:51:05.805104 [Node 2] event=ENTER_CS lamport=419
2026/11/08 21:51:06.806937 [Node 2] event=EXIT_CS lamport=419 deferred=0
2026/11/08 21:51:07.492757 [Node 2] event=RECEIVE_REQUEST from=1 lamport=418 local=420 state=RELEASED
2026/11/08 21:51:07.492757 [Node 2] event=SEND_REPLY to=1 lamport=421
2026/11/08 21:51:08.831672 [Node 2] event=RECEIVE_REQUEST from=3 lamport=421 local=422 state=RELEASED
2026/11/08 21:51:08.831672 [Node 2] event=SEND_REPLY to=3 lamport=423
2026/11/08 21:51:09.862417 [Node 2] event=REQUEST_CS lamport=424
2026/11/08 21:51:09.862417 [Node 2] event=SEND_REQUEST to=1 lamport=424
2026/11/08 21:51:09.862417 [Node 2] event=SEND_REQUEST to=3 lamport=424
2026/11/08 21:51:09.863016 [Node 2] event=RECEIVE_REPLY from=3 lamport=428 local=429 replies=1/2
2026/11/08 21:51:09.863016 [Node 2] event=RECEIVE_REPLY from=1 lamport=427 local=430 replies=2/2
2026/11/08 21:51:09.863016 [Node 2] event=ENTER_CS lamport=430
2026/11/08 21:51:10.864059 [Node 2] event=EXIT_CS lamport=430 deferred=0
2026/11/08 21:51:12.243857 [Node 2] event=RECEIVE_REQUEST from=1 lamport=428 local=431 state=RELEASED
2026/11/08 21:51:12.243877 [Node 2] event=SEND_REPLY to=1 lamport=432
2026/11/08 21:51:12.797544 [Node 2] event=RECEIVE_REQUEST from=3 lamport=431 local=433 state=RELEASED
2026/11/08 21:51:12.797544 [Node 2] event=SEND_REPLY to=3 lamport=434
2026/11/08 21:51:13.106345 [Node 2] event=REQUEST_CS lamport=435
2026/11/08 21:51:13.106345 [Node 2] event=SEND_REQUEST to=1 lamport=435
2026/11/08 21:51:13.106345 [Node 2] event=SEND_REQUEST to=3 lamport=435
2026/11/08 21:51:13.244530 [Node 2] event=RECEIVE_REPLY from=1 lamport=437 local=438 replies=1/2
2026/11/08 21:51:14.244921 [Node 2] event=RECEIVE_REPLY from=3 lamport=440 local=441 replies=2/2
2026/11/08 21:51:14.244921 [Node 2] event=ENTER_CS lamport=441
2026/11/08 21:51:15.246371 [Node 2] event=EXIT_CS lamport=441 deferred=0
2026/11/08 21:51:15.760122 [Node 2] event=RECEIVE_REQUEST from=1 lamport=439 local=442 state=RELEASED
2026/11/08 21:51:15.760727 [Node 2] event=SEND_REPLY to=1 lamport=443
2026/11/08 21:51:17.266719 [Node 2] event=REQUEST_CS lamport=444
2026/11/08 21:51:17.266719 [Node 2] event=SEND_REQUEST to=1 lamport=444
2026/11/08 21:51:17.266719 [Node 2] event=SEND_REQUEST to=3 lamport=444
2026/11/08 21:51:17.267352 [Node 2] event=RECEIVE_REPLY from=3 lamport=446 local=447 replies=1/2
2026/11/08 21:51:17.267352 [Node 2] event=RECEIVE_REPLY from=1 lamport=447 local=448 replies=2/2
2026/11/08 21:51:17.267898 [Node 2] event=ENTER_CS lamport=448
2026/11/08 21:51:17.930584 [Node 2] event=RECEIVE_REQUEST from=3 lamport=447 local=449 state=HELD
2026/11/08 21:51:17.930584 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:51:18.268470 [Node 2] event=EXIT_CS lamport=449 deferred=1
2026/11/08 21:51:18.268516 [Node 2] event=SEND_REPLY to=3 lamport=450
2026/11/08 21:51:18.685583 [Node 2] event=RECEIVE_REQUEST from=1 lamport=450 local=451 state=RELEASED
2026/11/08 21:51:18.685583 [Node 2] event=SEND_REPLY to=1 lamport=452
2026/11/08 21:51:19.717713 [Node 2] event=REQUEST_CS lamport=453
2026/11/08 21:51:19.717759 [Node 2] event=SEND_REQUEST to=1 lamport=453
2026/11/08 21:51:19.717759 [Node 2] event=SEND_REQUEST to=3 lamport=453
2026/11/08 21:51:19.718281 [Node 2] event=RECEIVE_REPLY from=3 lamport=455 local=456 replies=1/2
2026/11/08 21:51:20.270859 [Node 2] event=RECEIVE_REPLY from=1 lamport=456 local=457 replies=2/2
2026/11/08 21:51:20.270859 [Node 2] event=ENTER_CS lamport=457
2026/11/08 21:51:20.709220 [Node 2] event=RECEIVE_REQUEST from=3 lamport=456 local=458 state=HELD
2026/11/08 21:51:20.709748 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:51:21.271816 [Node 2] event=EXIT_CS lamport=458 deferred=1
2026/11/08 21:51:21.271816 [Node 2] event=SEND_REPLY to=3 lamport=459
2026/11/08 21:51:22.844979 [Node 2] event=REQUEST_CS lamport=460
2026/11/08 21:51:22.844979 [Node 2] event=SEND_REQUEST to=1 lamport=460
2026/11/08 21:51:22.844979 [Node 2] event=SEND_REQUEST to=3 lamport=460
2026/11/08 21:51:22.845731 [Node 2] event=RECEIVE_REPLY from=3 lamport=462 local=463 replies=1/2
2026/11/08 21:51:22.845731 [Node 2] event=RECEIVE_REPLY from=1 lamport=462 local=464 replies=2/2
2026/11/08 21:51:22.845731 [Node 2] event=ENTER_CS lamport=464
2026/11/08 21:51:23.731975 [Node 2] event=RECEIVE_REQUEST from=1 lamport=463 local=465 state=HELD
2026/11/08 21:51:23.732533 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:51:23.846688 [Node 2] event=EXIT_CS lamport=465 deferred=1
2026/11/08 21:51:23.846732 [Node 2] event=SEND_REPLY to=1 lamport=466
2026/11/08 21:51:25.016300 [Node 2] event=RECEIVE_REQUEST from=3 lamport=466 local=467 state=RELEASED
2026/11/08 21:51:25.016300 [Node 2] event=SEND_REPLY to=3 lamport=468
2026/11/08 21:51:27.570503 [Node 2] event=REQUEST_CS lamport=469
2026/11/08 21:51:27.570503 [Node 2] event=SEND_REQUEST to=1 lamport=469
2026/11/08 21:51:27.570503 [Node 2] event=SEND_REQUEST to=3 lamport=469
2026/11/08 21:51:27.571080 [Node 2] event=RECEIVE_REPLY from=3 lamport=472 local=473 replies=1/2
2026/11/08 21:51:27.571674 [Node 2] event=RECEIVE_REPLY from=1 lamport=471 local=474 replies=2/2
2026/11/08 21:51:27.571674 [Node 2] event=ENTER_CS lamport=474
2026/11/08 21:51:28.133799 [Node 2] event=RECEIVE_REQUEST from=1 lamport=472 local=475 state=HELD
2026/11/08 21:51:28.133799 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:51:28.233981 [Node 2] event=RECEIVE_REQUEST from=3 lamport=475 local=476 state=HELD
2026/11/08 21:51:28.234553 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:51:28.572335 [Node 2] event=EXIT_CS lamport=476 deferred=2
2026/11/08 21:51:28.572385 [Node 2] event=SEND_REPLY to=3 lamport=477
2026/11/08 21:51:28.572385 [Node 2] event=SEND_REPLY to=1 lamport=478
2026/11/08 21:51:32.454762 [Node 2] event=REQUEST_CS lamport=479
2026/11/08 21:51:32.454762 [Node 2] event=SEND_REQUEST to=1 lamport=479
2026/11/08 21:51:32.454762 [Node 2] event=SEND_REQUEST to=3 lamport=479
2026/11/08 21:51:32.455582 [Node 2] event=RECEIVE_REPLY from=3 lamport=483 local=484 replies=1/2
2026/11/08 21:51:32.455582 [Node 2] event=RECEIVE_REPLY from=1 lamport=482 local=485 replies=2/2
2026/11/08 21:51:32.455582 [Node 2] event=ENTER_CS lamport=485
2026/11/08 21:51:32.485501 [Node 2] event=RECEIVE_REQUEST from=3 lamport=484 local=486 state=HELD
2026/11/08 21:51:32.485501 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:51:32.567111 [Node 2] event=RECEIVE_REQUEST from=1 lamport=487 local=488 state=HELD
2026/11/08 21:51:32.567111 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:51:33.456325 [Node 2] event=EXIT_CS lamport=488 deferred=2
2026/11/08 21:51:33.456325 [Node 2] event=SEND_REPLY to=1 lamport=489
2026/11/08 21:51:33.456325 [Node 2] event=SEND_REPLY to=3 lamport=490
2026/11/08 21:51:34.644542 [Node 2] event=REQUEST_CS lamport=491
2026/11/08 21:51:34.644542 [Node 2] event=SEND_REQUEST to=1 lamport=491
2026/11/08 21:51:34.644542 [Node 2] event=SEND_REQUEST to=3 lamport=491
2026/11/08 21:51:34.645302 [Node 2] event=RECEIVE_REPLY from=3 lamport=494 local=495 replies=1/2
2026/11/08 21:51:35.458847 [Node 2] event=RECEIVE_REPLY from=1 lamport=495 local=496 replies=2/2
2026/11/08 21:51:35.458847 [Node 2] event=ENTER_CS lamport=496
2026/11/08 21:51:36.411855 [Node 2] event=RECEIVE_REQUEST from=3 lamport=495 local=497 state=HELD
2026/11/08 21:51:36.411855 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:51:36.459160 [Node 2] event=EXIT_CS lamport=497 deferred=1
2026/11/08 21:51:36.459160 [Node 2] event=SEND_REPLY to=3 lamport=498
2026/11/08 21:51:36.864310 [Node 2] event=RECEIVE_REQUEST from=1 lamport=498 local=499 state=RELEASED
2026/11/08 21:51:36.864310 [Node 2] event=SEND_REPLY to=1 lamport=500
2026/11/08 21:51:39.849876 [Node 2] event=REQUEST_CS lamport=501
2026/11/08 21:51:39.849876 [Node 2] event=SEND_REQUEST to=1 lamport=501
2026/11/08 21:51:39.849876 [Node 2] event=SEND_REQUEST to=3 lamport=501
2026/11/08 21:51:39.849876 [Node 2] event=RECEIVE_REPLY from=1 lamport=504 local=505 replies=1/2
2026/11/08 21:51:39.850878 [Node 2] event=RECEIVE_REPLY from=3 lamport=503 local=506 replies=2/2
2026/11/08 21:51:39.850878 [Node 2] event=ENTER_CS lamport=506
2026/11/08 21:51:40.152163 [Node 2] event=RECEIVE_REQUEST from=1 lamport=505 local=507 state=HELD
2026/11/08 21:51:40.152163 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:51:40.852067 [Node 2] event=EXIT_CS lamport=507 deferred=1
2026/11/08 21:51:40.852067 [Node 2] event=SEND_REPLY to=1 lamport=508
2026/11/08 21:51:41.034166 [Node 2] event=RECEIVE_REQUEST from=3 lamport=508 local=509 state=RELEASED
2026/11/08 21:51:41.034166 [Node 2] event=SEND_REPLY to=3 lamport=510
2026/11/08 21:51:43.045604 [Node 2] event=RECEIVE_REQUEST from=1 lamport=512 local=513 state=RELEASED
2026/11/08 21:51:43.045604 [Node 2] event=SEND_REPLY to=1 lamport=514
2026/11/08 21:51:43.783829 [Node 2] event=REQUEST_CS lamport=515
2026/11/08 21:51:43.783829 [Node 2] event=SEND_REQUEST to=1 lamport=515
2026/11/08 21:51:43.783829 [Node 2] event=SEND_REQUEST to=3 lamport=515
2026/11/08 21:51:43.784462 [Node 2] event=RECEIVE_REPLY from=3 lamport=517 local=518 replies=1/2
2026/11/08 21:51:44.046817 [Node 2] event=RECEIVE_REPLY from=1 lamport=518 local=519 replies=2/2
2026/11/08 21:51:44.046817 [Node 2] event=ENTER_CS lamport=519
2026/11/08 21:51:45.047654 [Node 2] event=EXIT_CS lamport=519 deferred=0
2026/11/08 21:51:46.320235 [Node 2] event=RECEIVE_REQUEST from=1 lamport=519 local=520 state=RELEASED
2026/11/08 21:51:46.320235 [Node 2] event=SEND_REPLY to=1 lamport=521
2026/11/08 21:51:46.596278 [Node 2] event=RECEIVE_REQUEST from=3 lamport=522 local=523 state=RELEASED
2026/11/08 21:51:46.596278 [Node 2] event=SEND_REPLY to=3 lamport=524
2026/11/08 21:51:47.594630 [Node 2] event=REQUEST_CS lamport=525
2026/11/08 21:51:47.594630 [Node 2] event=SEND_REQUEST to=3 lamport=525
2026/11/08 21:51:47.594630 [Node 2] event=SEND_REQUEST to=1 lamport=525
2026/11/08 21:51:47.595349 [Node 2] event=RECEIVE_REPLY from=1 lamport=527 local=528 replies=1/2
2026/11/08 21:51:48.321298 [Node 2] event=RECEIVE_REPLY from=3 lamport=528 local=529 replies=2/2
2026/11/08 21:51:48.321298 [Node 2] event=ENTER_CS lamport=529
2026/11/08 21:51:49.322029 [Node 2] event=EXIT_CS lamport=529 deferred=0
2026/11/08 21:51:50.668727 [Node 2] event=RECEIVE_REQUEST from=1 lamport=528 local=530 state=RELEASED
2026/11/08 21:51:50.668727 [Node 2] event=SEND_REPLY to=1 lamport=531
2026/11/08 21:51:51.744163 [Node 2] event=RECEIVE_REQUEST from=3 lamport=531 local=532 state=RELEASED
2026/11/08 21:51:51.744163 [Node 2] event=SEND_REPLY to=3 lamport=533
2026/11/08 21:51:52.798728 [Node 2] event=REQUEST_CS lamport=534
2026/11/08 21:51:52.798728 [Node 2] event=SEND_REQUEST to=1 lamport=534
2026/11/08 21:51:52.798728 [Node 2] event=SEND_REQUEST to=3 lamport=534
2026/11/08 21:51:52.799380 [Node 2] event=RECEIVE_REPLY from=1 lamport=537 local=538 replies=1/2
2026/11/08 21:51:52.799380 [Node 2] event=RECEIVE_REPLY from=3 lamport=539 local=540 replies=2/2
2026/11/08 21:51:52.799380 [Node 2] event=ENTER_CS lamport=540
2026/11/08 21:51:53.322103 [Node 2] event=RECEIVE_REQUEST from=1 lamport=538 local=541 state=HELD
2026/11/08 21:51:53.322103 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:51:53.799853 [Node 2] event=EXIT_CS lamport=541 deferred=1
2026/11/08 21:51:53.799853 [Node 2] event=SEND_REPLY to=1 lamport=542
2026/11/08 21:51:55.015374 [Node 2] event=RECEIVE_REQUEST from=3 lamport=542 local=543 state=RELEASED
2026/11/08 21:51:55.015374 [Node 2] event=SEND_REPLY to=3 lamport=544
2026/11/08 21:51:55.663282 [Node 2] event=REQUEST_CS lamport=545
2026/11/08 21:51:55.663282 [Node 2] event=SEND_REQUEST to=3 lamport=545
2026/11/08 21:51:55.663282 [Node 2] event=SEND_REQUEST to=1 lamport=545
2026/11/08 21:51:55.664436 [Node 2] event=RECEIVE_REPLY from=1 lamport=547 local=548 replies=1/2
2026/11/08 21:51:56.016085 [Node 2] event=RECEIVE_REPLY from=3 lamport=549 local=550 replies=2/2
2026/11/08 21:51:56.016085 [Node 2] event=ENTER_CS lamport=550
2026/11/08 21:51:57.017238 [Node 2] event=EXIT_CS lamport=550 deferred=0
2026/11/08 21:51:57.997063 [Node 2] event=RECEIVE_REQUEST from=1 lamport=548 local=551 state=RELEASED
2026/11/08 21:51:57.997063 [Node 2] event=SEND_REPLY to=1 lamport=552
2026/11/08 21:51:59.362568 [Node 2] event=REQUEST_CS lamport=553
2026/11/08 21:51:59.362568 [Node 2] event=SEND_REQUEST to=1 lamport=553
2026/11/08 21:51:59.362568 [Node 2] event=SEND_REQUEST to=3 lamport=553
2026/11/08 21:51:59.363523 [Node 2] event=RECEIVE_REPLY from=3 lamport=555 local=556 replies=1/2
2026/11/08 21:51:59.363523 [Node 2] event=RECEIVE_REPLY from=1 lamport=555 local=557 replies=2/2
2026/11/08 21:51:59.363523 [Node 2] event=ENTER_CS lamport=557
2026/11/08 21:51:59.921350 [Node 2] event=RECEIVE_REQUEST from=3 lamport=556 local=558 state=HELD
2026/11/08 21:51:59.921350 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:52:00.364440 [Node 2] event=EXIT_CS lamport=558 deferred=1
2026/11/08 21:52:00.364440 [Node 2] event=SEND_REPLY to=3 lamport=559
2026/11/08 21:52:00.527440 [Node 2] event=RECEIVE_REQUEST from=1 lamport=559 local=560 state=RELEASED
2026/11/08 21:52:00.527440 [Node 2] event=SEND_REPLY to=1 lamport=561
2026/11/08 21:52:03.382988 [Node 2] event=RECEIVE_REQUEST from=3 lamport=563 local=564 state=RELEASED
2026/11/08 21:52:03.382988 [Node 2] event=SEND_REPLY to=3 lamport=565
2026/11/08 21:52:03.769988 [Node 2] event=REQUEST_CS lamport=566
2026/11/08 21:52:03.769988 [Node 2] event=SEND_REQUEST to=1 lamport=566
2026/11/08 21:52:03.769988 [Node 2] event=SEND_REQUEST to=3 lamport=566
2026/11/08 21:52:03.769988 [Node 2] event=RECEIVE_REPLY from=1 lamport=568 local=569 replies=1/2
2026/11/08 21:52:04.384989 [Node 2] event=RECEIVE_REPLY from=3 lamport=569 local=570 replies=2/2
2026/11/08 21:52:04.384989 [Node 2] event=ENTER_CS lamport=570
2026/11/08 21:52:05.385416 [Node 2] event=EXIT_CS lamport=570 deferred=0
2026/11/08 21:52:05.629229 [Node 2] event=RECEIVE_REQUEST from=1 lamport=569 local=571 state=RELEASED
2026/11/08 21:52:05.629229 [Node 2] event=SEND_REPLY to=1 lamport=572
2026/11/08 21:52:07.949735 [Node 2] event=RECEIVE_REQUEST from=3 lamport=572 local=573 state=RELEASED
2026/11/08 21:52:07.949735 [Node 2] event=SEND_REPLY to=3 lamport=574
2026/11/08 21:52:08.475735 [Node 2] event=REQUEST_CS lamport=575
2026/11/08 21:52:08.475735 [Node 2] event=SEND_REQUEST to=1 lamport=575
2026/11/08 21:52:08.475735 [Node 2] event=SEND_REQUEST to=3 lamport=575
2026/11/08 21:52:08.476735 [Node 2] event=RECEIVE_REPLY from=1 lamport=578 local=579 replies=1/2
2026/11/08 21:52:08.530735 [Node 2] event=RECEIVE_REQUEST from=1 lamport=579 local=580 state=WANTED
2026/11/08 21:52:08.530735 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:52:08.951735 [Node 2] event=RECEIVE_REPLY from=3 lamport=582 local=583 replies=2/2
2026/11/08 21:52:08.951735 [Node 2] event=ENTER_CS lamport=583
2026/11/08 21:52:09.952272 [Node 2] event=EXIT_CS lamport=583 deferred=1
2026/11/08 21:52:09.952272 [Node 2] event=SEND_REPLY to=1 lamport=584
2026/11/08 21:52:11.223268 [Node 2] event=REQUEST_CS lamport=585
2026/11/08 21:52:11.223268 [Node 2] event=SEND_REQUEST to=1 lamport=585
2026/11/08 21:52:11.223268 [Node 2] event=SEND_REQUEST to=3 lamport=585
2026/11/08 21:52:11.223268 [Node 2] event=RECEIVE_REPLY from=1 lamport=587 local=588 replies=1/2
2026/11/08 21:52:11.224270 [Node 2] event=RECEIVE_REPLY from=3 lamport=587 local=589 replies=2/2
2026/11/08 21:52:11.224270 [Node 2] event=ENTER_CS lamport=589
2026/11/08 21:52:12.225820 [Node 2] event=EXIT_CS lamport=589 deferred=0
2026/11/08 21:52:12.342578 [Node 2] event=RECEIVE_REQUEST from=1 lamport=588 local=590 state=RELEASED
2026/11/08 21:52:12.342578 [Node 2] event=SEND_REPLY to=1 lamport=591
2026/11/08 21:52:12.927316 [Node 2] event=RECEIVE_REQUEST from=3 lamport=591 local=592 state=RELEASED
2026/11/08 21:52:12.927316 [Node 2] event=SEND_REPLY to=3 lamport=593
2026/11/08 21:52:15.474377 [Node 2] event=RECEIVE_REQUEST from=1 lamport=595 local=596 state=RELEASED
2026/11/08 21:52:15.474377 [Node 2] event=SEND_REPLY to=1 lamport=597
2026/11/08 21:52:15.783378 [Node 2] event=REQUEST_CS lamport=598
2026/11/08 21:52:15.783378 [Node 2] event=SEND_REQUEST to=1 lamport=598
2026/11/08 21:52:15.783378 [Node 2] event=SEND_REQUEST to=3 lamport=598
2026/11/08 21:52:15.783378 [Node 2] event=RECEIVE_REPLY from=3 lamport=600 local=601 replies=1/2
2026/11/08 21:52:16.241378 [Node 2] event=RECEIVE_REQUEST from=3 lamport=601 local=602 state=WANTED
2026/11/08 21:52:16.241378 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:52:16.476377 [Node 2] event=RECEIVE_REPLY from=1 lamport=604 local=605 replies=2/2
2026/11/08 21:52:16.476377 [Node 2] event=ENTER_CS lamport=605
2026/11/08 21:52:17.477377 [Node 2] event=EXIT_CS lamport=605 deferred=1
2026/11/08 21:52:17.477377 [Node 2] event=SEND_REPLY to=3 lamport=606
2026/11/08 21:52:18.596911 [Node 2] event=RECEIVE_REQUEST from=1 lamport=605 local=607 state=RELEASED
2026/11/08 21:52:18.596911 [Node 2] event=SEND_REPLY to=1 lamport=608
2026/11/08 21:52:19.843921 [Node 2] event=REQUEST_CS lamport=609
2026/11/08 21:52:19.843921 [Node 2] event=SEND_REQUEST to=1 lamport=609
2026/11/08 21:52:19.843921 [Node 2] event=SEND_REQUEST to=3 lamport=609
2026/11/08 21:52:19.844745 [Node 2] event=RECEIVE_REPLY from=1 lamport=613 local=614 replies=1/2
2026/11/08 21:52:19.844745 [Node 2] event=RECEIVE_REPLY from=3 lamport=611 local=615 replies=2/2
2026/11/08 21:52:19.844745 [Node 2] event=ENTER_CS lamport=615
2026/11/08 21:52:20.845806 [Node 2] event=EXIT_CS lamport=615 deferred=0
2026/11/08 21:52:21.510580 [Node 2] event=RECEIVE_REQUEST from=3 lamport=612 local=616 state=RELEASED
2026/11/08 21:52:21.510580 [Node 2] event=SEND_REPLY to=3 lamport=617
2026/11/08 21:52:22.226676 [Node 2] event=RECEIVE_REQUEST from=1 lamport=616 local=618 state=RELEASED
2026/11/08 21:52:22.226676 [Node 2] event=SEND_REPLY to=1 lamport=619
2026/11/08 21:52:23.303831 [Node 2] event=REQUEST_CS lamport=620
2026/11/08 21:52:23.303831 [Node 2] event=SEND_REQUEST to=1 lamport=620
2026/11/08 21:52:23.303831 [Node 2] event=SEND_REQUEST to=3 lamport=620
2026/11/08 21:52:23.303831 [Node 2] event=RECEIVE_REPLY from=3 lamport=622 local=623 replies=1/2
2026/11/08 21:52:23.511831 [Node 2] event=RECEIVE_REPLY from=1 lamport=623 local=624 replies=2/2
2026/11/08 21:52:23.511831 [Node 2] event=ENTER_CS lamport=624
2026/11/08 21:52:24.512831 [Node 2] event=EXIT_CS lamport=624 deferred=0
2026/11/08 21:52:25.432832 [Node 2] event=RECEIVE_REQUEST from=1 lamport=624 local=625 state=RELEASED
2026/11/08 21:52:25.432832 [Node 2] event=SEND_REPLY to=1 lamport=626
2026/11/08 21:52:26.044832 [Node 2] event=RECEIVE_REQUEST from=3 lamport=627 local=628 state=RELEASED
2026/11/08 21:52:26.044832 [Node 2] event=SEND_REPLY to=3 lamport=629
2026/11/08 21:52:28.171637 [Node 2] event=REQUEST_CS lamport=630
2026/11/08 21:52:28.171637 [Node 2] event=SEND_REQUEST to=1 lamport=630
2026/11/08 21:52:28.171637 [Node 2] event=SEND_REQUEST to=3 lamport=630
2026/11/08 21:52:28.172444 [Node 2] event=RECEIVE_REPLY from=1 lamport=632 local=633 replies=1/2
2026/11/08 21:52:28.172444 [Node 2] event=RECEIVE_REPLY from=3 lamport=633 local=634 replies=2/2
2026/11/08 21:52:28.172444 [Node 2] event=ENTER_CS lamport=634
2026/11/08 21:52:29.173110 [Node 2] event=EXIT_CS lamport=634 deferred=0
2026/11/08 21:52:30.430493 [Node 2] event=RECEIVE_REQUEST from=1 lamport=633 local=635 state=RELEASED
2026/11/08 21:52:30.430493 [Node 2] event=SEND_REPLY to=1 lamport=636
2026/11/08 21:52:30.678492 [Node 2] event=RECEIVE_REQUEST from=3 lamport=636 local=637 state=RELEASED
2026/11/08 21:52:30.678492 [Node 2] event=SEND_REPLY to=3 lamport=638
2026/11/08 21:52:32.514293 [Node 2] event=REQUEST_CS lamport=639
2026/11/08 21:52:32.514293 [Node 2] event=SEND_REQUEST to=1 lamport=639
2026/11/08 21:52:32.514293 [Node 2] event=SEND_REQUEST to=3 lamport=639
2026/11/08 21:52:32.515248 [Node 2] event=RECEIVE_REPLY from=3 lamport=642 local=643 replies=1/2
2026/11/08 21:52:32.515248 [Node 2] event=RECEIVE_REPLY from=1 lamport=641 local=644 replies=2/2
2026/11/08 21:52:32.515248 [Node 2] event=ENTER_CS lamport=644
2026/11/08 21:52:33.019374 [Node 2] event=RECEIVE_REQUEST from=1 lamport=642 local=645 state=HELD
2026/11/08 21:52:33.019374 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:52:33.515802 [Node 2] event=EXIT_CS lamport=645 deferred=1
2026/11/08 21:52:33.515859 [Node 2] event=SEND_REPLY to=1 lamport=646
2026/11/08 21:52:33.893090 [Node 2] event=RECEIVE_REQUEST from=3 lamport=645 local=647 state=RELEASED
2026/11/08 21:52:33.893113 [Node 2] event=SEND_REPLY to=3 lamport=648
2026/11/08 21:52:35.552890 [Node 2] event=REQUEST_CS lamport=649
2026/11/08 21:52:35.552890 [Node 2] event=SEND_REQUEST to=1 lamport=649
2026/11/08 21:52:35.552890 [Node 2] event=SEND_REQUEST to=3 lamport=649
2026/11/08 21:52:35.553670 [Node 2] event=RECEIVE_REPLY from=3 lamport=652 local=653 replies=1/2
2026/11/08 21:52:35.553670 [Node 2] event=RECEIVE_REPLY from=1 lamport=651 local=654 replies=2/2
2026/11/08 21:52:35.553670 [Node 2] event=ENTER_CS lamport=654
2026/11/08 21:52:36.323273 [Node 2] event=RECEIVE_REQUEST from=1 lamport=652 local=655 state=HELD
2026/11/08 21:52:36.323292 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:52:36.554402 [Node 2] event=EXIT_CS lamport=655 deferred=1
2026/11/08 21:52:36.554402 [Node 2] event=SEND_REPLY to=1 lamport=656
2026/11/08 21:52:38.828896 [Node 2] event=RECEIVE_REQUEST from=3 lamport=655 local=657 state=RELEASED
2026/11/08 21:52:38.828896 [Node 2] event=SEND_REPLY to=3 lamport=658
2026/11/08 21:52:39.296317 [Node 2] event=RECEIVE_REQUEST from=1 lamport=660 local=661 state=RELEASED
2026/11/08 21:52:39.296913 [Node 2] event=SEND_REPLY to=1 lamport=662
2026/11/08 21:52:40.133315 [Node 2] event=REQUEST_CS lamport=663
2026/11/08 21:52:40.133315 [Node 2] event=SEND_REQUEST to=1 lamport=663
2026/11/08 21:52:40.133315 [Node 2] event=SEND_REQUEST to=3 lamport=663
2026/11/08 21:52:40.134104 [Node 2] event=RECEIVE_REPLY from=3 lamport=665 local=666 replies=1/2
2026/11/08 21:52:40.831666 [Node 2] event=RECEIVE_REPLY from=1 lamport=666 local=667 replies=2/2
2026/11/08 21:52:40.831666 [Node 2] event=ENTER_CS lamport=667
2026/11/08 21:52:41.832014 [Node 2] event=EXIT_CS lamport=667 deferred=0
2026/11/08 21:52:43.336287 [Node 2] event=REQUEST_CS lamport=668
2026/11/08 21:52:43.336352 [Node 2] event=SEND_REQUEST to=1 lamport=668
2026/11/08 21:52:43.336352 [Node 2] event=SEND_REQUEST to=3 lamport=668
2026/11/08 21:52:43.336876 [Node 2] event=RECEIVE_REPLY from=1 lamport=670 local=671 replies=1/2
2026/11/08 21:52:43.337927 [Node 2] event=RECEIVE_REPLY from=3 lamport=670 local=672 replies=2/2
2026/11/08 21:52:43.337927 [Node 2] event=ENTER_CS lamport=672
2026/11/08 21:52:43.813437 [Node 2] event=RECEIVE_REQUEST from=3 lamport=671 local=673 state=HELD
2026/11/08 21:52:43.813437 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:52:43.825579 [Node 2] event=RECEIVE_REQUEST from=1 lamport=674 local=675 state=HELD
2026/11/08 21:52:43.825579 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:52:44.338659 [Node 2] event=EXIT_CS lamport=675 deferred=2
2026/11/08 21:52:44.338659 [Node 2] event=SEND_REPLY to=1 lamport=676
2026/11/08 21:52:44.338659 [Node 2] event=SEND_REPLY to=3 lamport=677
2026/11/08 21:52:47.784360 [Node 2] event=REQUEST_CS lamport=678
2026/11/08 21:52:47.784360 [Node 2] event=SEND_REQUEST to=1 lamport=678
2026/11/08 21:52:47.784360 [Node 2] event=SEND_REQUEST to=3 lamport=678
2026/11/08 21:52:47.785245 [Node 2] event=RECEIVE_REPLY from=3 lamport=681 local=682 replies=1/2
2026/11/08 21:52:47.785270 [Node 2] event=RECEIVE_REPLY from=1 lamport=682 local=683 replies=2/2
2026/11/08 21:52:47.785270 [Node 2] event=ENTER_CS lamport=683
2026/11/08 21:52:48.388782 [Node 2] event=RECEIVE_REQUEST from=1 lamport=683 local=684 state=HELD
2026/11/08 21:52:48.389424 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:52:48.402490 [Node 2] event=RECEIVE_REQUEST from=3 lamport=686 local=687 state=HELD
2026/11/08 21:52:48.402490 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:52:48.785579 [Node 2] event=EXIT_CS lamport=687 deferred=2
2026/11/08 21:52:48.785579 [Node 2] event=SEND_REPLY to=3 lamport=688
2026/11/08 21:52:48.785579 [Node 2] event=SEND_REPLY to=1 lamport=689
2026/11/08 21:52:51.133723 [Node 2] event=REQUEST_CS lamport=690
2026/11/08 21:52:51.133723 [Node 2] event=SEND_REQUEST to=1 lamport=690
2026/11/08 21:52:51.133723 [Node 2] event=SEND_REQUEST to=3 lamport=690
2026/11/08 21:52:51.134718 [Node 2] event=RECEIVE_REPLY from=1 lamport=693 local=694 replies=1/2
2026/11/08 21:52:51.134718 [Node 2] event=RECEIVE_REPLY from=3 lamport=694 local=695 replies=2/2
2026/11/08 21:52:51.135245 [Node 2] event=ENTER_CS lamport=695
2026/11/08 21:52:52.135755 [Node 2] event=EXIT_CS lamport=695 deferred=0
2026/11/08 21:52:52.198385 [Node 2] event=RECEIVE_REQUEST from=3 lamport=695 local=696 state=RELEASED
2026/11/08 21:52:52.198917 [Node 2] event=SEND_REPLY to=3 lamport=697
2026/11/08 21:52:52.881008 [Node 2] event=RECEIVE_REQUEST from=1 lamport=698 local=699 state=RELEASED
2026/11/08 21:52:52.881008 [Node 2] event=SEND_REPLY to=1 lamport=700
2026/11/08 21:52:55.465959 [Node 2] event=REQUEST_CS lamport=701
2026/11/08 21:52:55.465959 [Node 2] event=SEND_REQUEST to=3 lamport=701
2026/11/08 21:52:55.465959 [Node 2] event=SEND_REQUEST to=1 lamport=701
2026/11/08 21:52:55.466536 [Node 2] event=RECEIVE_REPLY from=3 lamport=703 local=704 replies=1/2
2026/11/08 21:52:55.466536 [Node 2] event=RECEIVE_REPLY from=1 lamport=704 local=705 replies=2/2
2026/11/08 21:52:55.466536 [Node 2] event=ENTER_CS lamport=705
2026/11/08 21:52:56.467052 [Node 2] event=EXIT_CS lamport=705 deferred=0
2026/11/08 21:52:56.993302 [Node 2] event=RECEIVE_REQUEST from=3 lamport=704 local=706 state=RELEASED
2026/11/08 21:52:56.993302 [Node 2] event=SEND_REPLY to=3 lamport=707
2026/11/08 21:52:57.829182 [Node 2] event=RECEIVE_REQUEST from=1 lamport=707 local=708 state=RELEASED
2026/11/08 21:52:57.829182 [Node 2] event=SEND_REPLY to=1 lamport=709
2026/11/08 21:53:00.099526 [Node 2] event=RECEIVE_REQUEST from=1 lamport=713 local=714 state=RELEASED
2026/11/08 21:53:00.099526 [Node 2] event=SEND_REPLY to=1 lamport=715
2026/11/08 21:53:00.170818 [Node 2] event=RECEIVE_REQUEST from=3 lamport=716 local=717 state=RELEASED
2026/11/08 21:53:00.170818 [Node 2] event=SEND_REPLY to=3 lamport=718
2026/11/08 21:53:00.358885 [Node 2] event=REQUEST_CS lamport=719
2026/11/08 21:53:00.358929 [Node 2] event=SEND_REQUEST to=1 lamport=719
2026/11/08 21:53:00.358929 [Node 2] event=SEND_REQUEST to=3 lamport=719
2026/11/08 21:53:01.100212 [Node 2] event=RECEIVE_REPLY from=1 lamport=721 local=722 replies=1/2
2026/11/08 21:53:02.101392 [Node 2] event=RECEIVE_REPLY from=3 lamport=724 local=725 replies=2/2
2026/11/08 21:53:02.101392 [Node 2] event=ENTER_CS lamport=725
2026/11/08 21:53:03.102121 [Node 2] event=EXIT_CS lamport=725 deferred=0
2026/11/08 21:53:03.404129 [Node 2] event=RECEIVE_REQUEST from=1 lamport=723 local=726 state=RELEASED
2026/11/08 21:53:03.404129 [Node 2] event=SEND_REPLY to=1 lamport=727
2026/11/08 21:53:05.356513 [Node 2] event=RECEIVE_REQUEST from=3 lamport=727 local=728 state=RELEASED
2026/11/08 21:53:05.356513 [Node 2] event=SEND_REPLY to=3 lamport=729
2026/11/08 21:53:05.943536 [Node 2] event=REQUEST_CS lamport=730
2026/11/08 21:53:05.943536 [Node 2] event=SEND_REQUEST to=1 lamport=730
2026/11/08 21:53:05.943536 [Node 2] event=SEND_REQUEST to=3 lamport=730
2026/11/08 21:53:05.944099 [Node 2] event=RECEIVE_REPLY from=1 lamport=732 local=733 replies=1/2
2026/11/08 21:53:06.357310 [Node 2] event=RECEIVE_REPLY from=3 lamport=733 local=734 replies=2/2
2026/11/08 21:53:06.357310 [Node 2] event=ENTER_CS lamport=734
2026/11/08 21:53:07.358017 [Node 2] event=EXIT_CS lamport=734 deferred=0
2026/11/08 21:53:07.934340 [Node 2] event=RECEIVE_REQUEST from=1 lamport=733 local=735 state=RELEASED
2026/11/08 21:53:07.934340 [Node 2] event=SEND_REPLY to=1 lamport=736
2026/11/08 21:53:09.378339 [Node 2] event=RECEIVE_REQUEST from=3 lamport=736 local=737 state=RELEASED
2026/11/08 21:53:09.378339 [Node 2] event=SEND_REPLY to=3 lamport=738
2026/11/08 21:53:09.676010 [Node 2] event=REQUEST_CS lamport=739
2026/11/08 21:53:09.676010 [Node 2] event=SEND_REQUEST to=1 lamport=739
2026/11/08 21:53:09.676010 [Node 2] event=SEND_REQUEST to=3 lamport=739
2026/11/08 21:53:09.676562 [Node 2] event=RECEIVE_REPLY from=1 lamport=741 local=742 replies=1/2
2026/11/08 21:53:10.379711 [Node 2] event=RECEIVE_REPLY from=3 lamport=743 local=744 replies=2/2
2026/11/08 21:53:10.379711 [Node 2] event=ENTER_CS lamport=744
2026/11/08 21:53:11.379939 [Node 2] event=EXIT_CS lamport=744 deferred=0
2026/11/08 21:53:11.570234 [Node 2] event=RECEIVE_REQUEST from=3 lamport=744 local=745 state=RELEASED
2026/11/08 21:53:11.570234 [Node 2] event=SEND_REPLY to=3 lamport=746
2026/11/08 21:53:12.308074 [Node 2] event=RECEIVE_REQUEST from=1 lamport=747 local=748 state=RELEASED
2026/11/08 21:53:12.308074 [Node 2] event=SEND_REPLY to=1 lamport=749
2026/11/08 21:53:14.241356 [Node 2] event=RECEIVE_REQUEST from=3 lamport=751 local=752 state=RELEASED
2026/11/08 21:53:14.241356 [Node 2] event=SEND_REPLY to=3 lamport=753
2026/11/08 21:53:14.648071 [Node 2] event=REQUEST_CS lamport=754
2026/11/08 21:53:14.648115 [Node 2] event=SEND_REQUEST to=1 lamport=754
2026/11/08 21:53:14.648115 [Node 2] event=SEND_REQUEST to=3 lamport=754
2026/11/08 21:53:14.648644 [Node 2] event=RECEIVE_REPLY from=1 lamport=756 local=757 replies=1/2
2026/11/08 21:53:15.242475 [Node 2] event=RECEIVE_REPLY from=3 lamport=757 local=758 replies=2/2
2026/11/08 21:53:15.242475 [Node 2] event=ENTER_CS lamport=758
2026/11/08 21:53:16.242787 [Node 2] event=EXIT_CS lamport=758 deferred=0
2026/11/08 21:53:16.295397 [Node 2] event=RECEIVE_REQUEST from=1 lamport=757 local=759 state=RELEASED
2026/11/08 21:53:16.295397 [Node 2] event=SEND_REPLY to=1 lamport=760
2026/11/08 21:53:17.622350 [Node 2] event=REQUEST_CS lamport=761
2026/11/08 21:53:17.622350 [Node 2] event=SEND_REQUEST to=3 lamport=761
2026/11/08 21:53:17.622350 [Node 2] event=SEND_REQUEST to=1 lamport=761
2026/11/08 21:53:17.622912 [Node 2] event=RECEIVE_REPLY from=1 lamport=763 local=764 replies=1/2
2026/11/08 21:53:17.622912 [Node 2] event=RECEIVE_REPLY from=3 lamport=763 local=765 replies=2/2
2026/11/08 21:53:17.622912 [Node 2] event=ENTER_CS lamport=765
2026/11/08 21:53:18.624114 [Node 2] event=EXIT_CS lamport=765 deferred=0
2026/11/08 21:53:18.840093 [Node 2] event=RECEIVE_REQUEST from=1 lamport=764 local=766 state=RELEASED
2026/11/08 21:53:18.840723 [Node 2] event=SEND_REPLY to=1 lamport=767
2026/11/08 21:53:18.872813 [Node 2] event=RECEIVE_REQUEST from=3 lamport=767 local=768 state=RELEASED
2026/11/08 21:53:18.872813 [Node 2] event=SEND_REPLY to=3 lamport=769
2026/11/08 21:53:21.980584 [Node 2] event=REQUEST_CS lamport=770
2026/11/08 21:53:21.980654 [Node 2] event=SEND_REQUEST to=1 lamport=770
2026/11/08 21:53:21.980654 [Node 2] event=SEND_REQUEST to=3 lamport=770
2026/11/08 21:53:21.981175 [Node 2] event=RECEIVE_REPLY from=3 lamport=774 local=775 replies=1/2
2026/11/08 21:53:21.981175 [Node 2] event=RECEIVE_REPLY from=1 lamport=773 local=776 replies=2/2
2026/11/08 21:53:21.981175 [Node 2] event=ENTER_CS lamport=776
2026/11/08 21:53:22.981795 [Node 2] event=EXIT_CS lamport=776 deferred=0
2026/11/08 21:53:23.042352 [Node 2] event=RECEIVE_REQUEST from=3 lamport=775 local=777 state=RELEASED
2026/11/08 21:53:23.042352 [Node 2] event=SEND_REPLY to=3 lamport=778
2026/11/08 21:53:23.249318 [Node 2] event=RECEIVE_REQUEST from=1 lamport=778 local=779 state=RELEASED
2026/11/08 21:53:23.249318 [Node 2] event=SEND_REPLY to=1 lamport=780
2026/11/08 21:53:26.216897 [Node 2] event=RECEIVE_REQUEST from=1 lamport=784 local=785 state=RELEASED
2026/11/08 21:53:26.216897 [Node 2] event=SEND_REPLY to=1 lamport=786
2026/11/08 21:53:26.397720 [Node 2] event=REQUEST_CS lamport=787
2026/11/08 21:53:26.397761 [Node 2] event=SEND_REQUEST to=3 lamport=787
2026/11/08 21:53:26.397761 [Node 2] event=SEND_REQUEST to=1 lamport=787
2026/11/08 21:53:26.398360 [Node 2] event=RECEIVE_REPLY from=3 lamport=789 local=790 replies=1/2
2026/11/08 21:53:26.778923 [Node 2] event=RECEIVE_REQUEST from=3 lamport=790 local=791 state=WANTED
2026/11/08 21:53:26.778923 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:53:27.218626 [Node 2] event=RECEIVE_REPLY from=1 lamport=793 local=794 replies=2/2
2026/11/08 21:53:27.218626 [Node 2] event=ENTER_CS lamport=794
2026/11/08 21:53:28.219146 [Node 2] event=EXIT_CS lamport=794 deferred=1
2026/11/08 21:53:28.219146 [Node 2] event=SEND_REPLY to=3 lamport=795
2026/11/08 21:53:30.136332 [Node 2] event=RECEIVE_REQUEST from=1 lamport=794 local=796 state=RELEASED
2026/11/08 21:53:30.136839 [Node 2] event=SEND_REPLY to=1 lamport=797
2026/11/08 21:53:31.461770 [Node 2] event=RECEIVE_REQUEST from=3 lamport=799 local=800 state=RELEASED
2026/11/08 21:53:31.461770 [Node 2] event=SEND_REPLY to=3 lamport=801
2026/11/08 21:53:31.714403 [Node 2] event=REQUEST_CS lamport=802
2026/11/08 21:53:31.714403 [Node 2] event=SEND_REQUEST to=1 lamport=802
2026/11/08 21:53:31.714403 [Node 2] event=SEND_REQUEST to=3 lamport=802
2026/11/08 21:53:31.714403 [Node 2] event=RECEIVE_REPLY from=1 lamport=804 local=805 replies=1/2
2026/11/08 21:53:32.462733 [Node 2] event=RECEIVE_REPLY from=3 lamport=805 local=806 replies=2/2
2026/11/08 21:53:32.462733 [Node 2] event=ENTER_CS lamport=806
2026/11/08 21:53:33.463540 [Node 2] event=EXIT_CS lamport=806 deferred=0
2026/11/08 21:53:33.627070 [Node 2] event=RECEIVE_REQUEST from=3 lamport=806 local=807 state=RELEASED
2026/11/08 21:53:33.627070 [Node 2] event=SEND_REPLY to=3 lamport=808
2026/11/08 21:53:34.468230 [Node 2] event=REQUEST_CS lamport=809
2026/11/08 21:53:34.468272 [Node 2] event=SEND_REQUEST to=1 lamport=809
2026/11/08 21:53:34.468272 [Node 2] event=SEND_REQUEST to=3 lamport=809
2026/11/08 21:53:34.468793 [Node 2] event=RECEIVE_REPLY from=1 lamport=811 local=812 replies=1/2
2026/11/08 21:53:34.627891 [Node 2] event=RECEIVE_REPLY from=3 lamport=812 local=813 replies=2/2
2026/11/08 21:53:34.627891 [Node 2] event=ENTER_CS lamport=813
2026/11/08 21:53:34.639947 [Node 2] event=RECEIVE_REQUEST from=1 lamport=812 local=814 state=HELD
2026/11/08 21:53:34.639947 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:53:35.628180 [Node 2] event=EXIT_CS lamport=814 deferred=1
2026/11/08 21:53:35.628233 [Node 2] event=SEND_REPLY to=1 lamport=815
2026/11/08 21:53:37.050969 [Node 2] event=RECEIVE_REQUEST from=3 lamport=815 local=816 state=RELEASED
2026/11/08 21:53:37.050969 [Node 2] event=SEND_REPLY to=3 lamport=817
2026/11/08 21:53:39.039847 [Node 2] event=REQUEST_CS lamport=818
2026/11/08 21:53:39.039847 [Node 2] event=SEND_REQUEST to=1 lamport=818
2026/11/08 21:53:39.039847 [Node 2] event=SEND_REQUEST to=3 lamport=818
2026/11/08 21:53:39.040447 [Node 2] event=RECEIVE_REPLY from=3 lamport=821 local=822 replies=1/2
2026/11/08 21:53:39.040447 [Node 2] event=RECEIVE_REPLY from=1 lamport=820 local=823 replies=2/2
2026/11/08 21:53:39.040447 [Node 2] event=ENTER_CS lamport=823
2026/11/08 21:53:39.296639 [Node 2] event=RECEIVE_REQUEST from=1 lamport=821 local=824 state=HELD
2026/11/08 21:53:39.296639 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:53:40.040607 [Node 2] event=EXIT_CS lamport=824 deferred=1
2026/11/08 21:53:40.040607 [Node 2] event=SEND_REPLY to=1 lamport=825
2026/11/08 21:53:40.936519 [Node 2] event=RECEIVE_REQUEST from=3 lamport=824 local=826 state=RELEASED
2026/11/08 21:53:40.936519 [Node 2] event=SEND_REPLY to=3 lamport=827
2026/11/08 21:53:42.693477 [Node 2] event=RECEIVE_REQUEST from=1 lamport=829 local=830 state=RELEASED
2026/11/08 21:53:42.693500 [Node 2] event=SEND_REPLY to=1 lamport=831
2026/11/08 21:53:43.366391 [Node 2] event=REQUEST_CS lamport=832
2026/11/08 21:53:43.366391 [Node 2] event=SEND_REQUEST to=1 lamport=832
2026/11/08 21:53:43.366391 [Node 2] event=SEND_REQUEST to=3 lamport=832
2026/11/08 21:53:43.366916 [Node 2] event=RECEIVE_REPLY from=3 lamport=834 local=835 replies=1/2
2026/11/08 21:53:43.532502 [Node 2] event=RECEIVE_REQUEST from=3 lamport=835 local=836 state=WANTED
2026/11/08 21:53:43.532502 [Node 2] event=DEFER_REPLY to=3
2026/11/08 21:53:43.693733 [Node 2] event=RECEIVE_REPLY from=1 lamport=838 local=839 replies=2/2
2026/11/08 21:53:43.693733 [Node 2] event=ENTER_CS lamport=839
2026/11/08 21:53:44.694563 [Node 2] event=EXIT_CS lamport=839 deferred=1
2026/11/08 21:53:44.694563 [Node 2] event=SEND_REPLY to=3 lamport=840
2026/11/08 21:53:46.442112 [Node 2] event=REQUEST_CS lamport=841
2026/11/08 21:53:46.442171 [Node 2] event=SEND_REQUEST to=3 lamport=841
2026/11/08 21:53:46.442171 [Node 2] event=SEND_REQUEST to=1 lamport=841
2026/11/08 21:53:46.442700 [Node 2] event=RECEIVE_REPLY from=3 lamport=843 local=844 replies=1/2
2026/11/08 21:53:46.442700 [Node 2] event=RECEIVE_REPLY from=1 lamport=843 local=845 replies=2/2
2026/11/08 21:53:46.442700 [Node 2] event=ENTER_CS lamport=845
2026/11/08 21:53:47.115312 [Node 2] event=RECEIVE_REQUEST from=1 lamport=844 local=846 state=HELD
2026/11/08 21:53:47.115312 [Node 2] event=DEFER_REPLY to=1
2026/11/08 21:53:47.443371 [Node 2] event=EXIT_CS lamport=846 deferred=1
2026/11/08 21:53:47.443422 [Node 2] event=SEND_REPLY to=1 lamport=847
2026/11/08 21:53:47.555583 [Node 2] event=RECEIVE_REQUEST from=3 lamport=847 local=848 state=RELEASED
2026/11/08 21:53:47.555583 [Node 2] event=SEND_REPLY to=3 lamport=849
2026/11/08 21:53:48.917249 [Node 2] event=REQUEST_CS lamport=850
2026/11/08 21:53:48.917249 [Node 2] event=SEND_REQUEST to=3 lamport=850
2026/11/08 21:53:48.917249 [Node 2] event=SEND_REQUEST to=1 lamport=850
2026/11/08 21:53:48.917829 [Node 2] event=RECEIVE_REPLY from=1 lamport=852 local=853 replies=1/2
2026/11/08 21:53:49.445605 [Node 2] event=RECEIVE_REPLY from=3 lamport=853 local=854 replies=2/2
2026/11/08 21:53:49.445605 [Node 2] event=ENTER_CS lamport=854
2026/11/08 21:53:50.446876 [Node 2] event=EXIT_CS lamport=854 deferred=0
2026/11/08 21:53:51.800334 [Node 2] event=RECEIVE_REQUEST from=3 lamport=854 local=855 state=RELEASED
2026/11/08 21:53:51.800334 [Node 2] event=SEND_REPLY to=3 lamport=856
2026/11/08 21:53:52.246978 [Node 2] event=RECEIVE_REQUEST from=1 lamport=857 local=858 state=RELEASED
2026/11/08 21:53:52.246978 [Node 2] event=SEND_REPLY to=1 lamport=859
2026/11/08 21:53:52.802681 [Node 2] event=REQUEST_CS lamport=860
2026/11/08 21:53:52.802711 [Node 2] event=SEND_REQUEST to=1 lamport=860
2026/11/08 21:53:52.802711 [Node 2] event=SEND_REQUEST to=3 lamport=860
2026/11/08 21:53:52.803230 [Node 2] event=RECEIVE_REPLY from=3 lamport=862 local=863 replies=1/2
2026/11/08 21:53:53.801942 [Node 2] event=RECEIVE_REPLY from=1 lamport=863 local=864 replies=2/2
2026/11/08 21:53:53.801942 [Node 2] event=ENTER_CS lamport=864
2026/11/08 21:53:54.802231 [Node 2] event=EXIT_CS lamport=864 deferred=0
2026/11/08 21:53:55.579993 [Node 2] event=RECEIVE_REQUEST from=3 lamport=863 local=865 state=RELEASED
2026/11/08 21:53:55.579993 [Node 2] event=SEND_REPLY to=3 lamport=866
2026/11/08 21:53:55.904672 [Node 2] event=REQUEST_CS lamport=867
2026/11/08 21:53:55.904672 [Node 2] event=SEND_REQUEST to=1 lamport=867
2026/11/08 21:53:55.904672 [Node 2] event=SEND_REQUEST to=3 lamport=867
2026/11/08 21:53:55.905556 [Node 2] event=RECEIVE_REPLY from=1 lamport=869 local=870 replies=1/2
2026/11/08 21:53:56.307967 [Node 2] event=RETRY to=3 error="rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:51292->127.0.0.1:5003: wsarecv: An existing connection was forcibly closed by the remote host."
2026/11/08 21:53:56.490202 [Node 2] event=RECEIVE_REQUEST from=1 lamport=870 local=871 state=WANTED
2026/11/08 21:53:56.490202 [Node 2] event=DEFER_REPLY to=1

```

### Node 3

```
2026/11/08 21:48:31.339760 [Node 3] event=STARTUP address=localhost:5003 peers=2
2026/11/08 21:48:31.800801 [Node 3] event=RECEIVE_REQUEST from=2 lamport=1 local=2 state=RELEASED
2026/11/08 21:48:31.800801 [Node 3] event=SEND_REPLY to=2 lamport=3
2026/11/08 21:48:31.897149 [Node 3] event=RECEIVE_REQUEST from=1 lamport=4 local=5 state=RELEASED
2026/11/08 21:48:31.898150 [Node 3] event=SEND_REPLY to=1 lamport=6
2026/11/08 21:48:33.484288 [Node 3] event=REQUEST_CS lamport=7
2026/11/08 21:48:33.484288 [Node 3] event=SEND_REQUEST to=1 lamport=7
2026/11/08 21:48:33.484288 [Node 3] event=SEND_REQUEST to=2 lamport=7
2026/11/08 21:48:33.486073 [Node 3] event=RECEIVE_REPLY from=2 lamport=9 local=10 replies=1/2
2026/11/08 21:48:33.802864 [Node 3] event=RECEIVE_REPLY from=1 lamport=10 local=11 replies=2/2
2026/11/08 21:48:33.802864 [Node 3] event=ENTER_CS lamport=11
2026/11/08 21:48:33.946452 [Node 3] event=RECEIVE_REQUEST from=2 lamport=10 local=12 state=HELD
2026/11/08 21:48:33.946452 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:48:34.804571 [Node 3] event=EXIT_CS lamport=12 deferred=1
2026/11/08 21:48:34.804571 [Node 3] event=SEND_REPLY to=2 lamport=13
2026/11/08 21:48:35.152282 [Node 3] event=RECEIVE_REQUEST from=1 lamport=13 local=14 state=RELEASED
2026/11/08 21:48:35.152835 [Node 3] event=SEND_REPLY to=1 lamport=15
2026/11/08 21:48:36.758053 [Node 3] event=REQUEST_CS lamport=16
2026/11/08 21:48:36.758110 [Node 3] event=SEND_REQUEST to=1 lamport=16
2026/11/08 21:48:36.758110 [Node 3] event=SEND_REQUEST to=2 lamport=16
2026/11/08 21:48:36.758632 [Node 3] event=RECEIVE_REPLY from=2 lamport=18 local=19 replies=1/2
2026/11/08 21:48:36.805787 [Node 3] event=RECEIVE_REPLY from=1 lamport=19 local=20 replies=2/2
2026/11/08 21:48:36.805787 [Node 3] event=ENTER_CS lamport=20
2026/11/08 21:48:37.806091 [Node 3] event=EXIT_CS lamport=20 deferred=0
2026/11/08 21:48:38.236784 [Node 3] event=RECEIVE_REQUEST from=2 lamport=19 local=21 state=RELEASED
2026/11/08 21:48:38.236784 [Node 3] event=SEND_REPLY to=2 lamport=22
2026/11/08 21:48:39.961400 [Node 3] event=RECEIVE_REQUEST from=1 lamport=22 local=23 state=RELEASED
2026/11/08 21:48:39.961400 [Node 3] event=SEND_REPLY to=1 lamport=24
2026/11/08 21:48:40.054925 [Node 3] event=REQUEST_CS lamport=25
2026/11/08 21:48:40.054925 [Node 3] event=SEND_REQUEST to=1 lamport=25
2026/11/08 21:48:40.054925 [Node 3] event=SEND_REQUEST to=2 lamport=25
2026/11/08 21:48:40.055855 [Node 3] event=RECEIVE_REPLY from=2 lamport=28 local=29 replies=1/2
2026/11/08 21:48:40.962924 [Node 3] event=RECEIVE_REPLY from=1 lamport=29 local=30 replies=2/2
2026/11/08 21:48:40.962924 [Node 3] event=ENTER_CS lamport=30
2026/11/08 21:48:41.963094 [Node 3] event=EXIT_CS lamport=30 deferred=0
2026/11/08 21:48:42.037145 [Node 3] event=RECEIVE_REQUEST from=2 lamport=29 local=31 state=RELEASED
2026/11/08 21:48:42.037145 [Node 3] event=SEND_REPLY to=2 lamport=32
2026/11/08 21:48:43.447654 [Node 3] event=REQUEST_CS lamport=33
2026/11/08 21:48:43.447654 [Node 3] event=SEND_REQUEST to=1 lamport=33
2026/11/08 21:48:43.447654 [Node 3] event=SEND_REQUEST to=2 lamport=33
2026/11/08 21:48:43.448653 [Node 3] event=RECEIVE_REPLY from=2 lamport=35 local=36 replies=1/2
2026/11/08 21:48:43.448653 [Node 3] event=RECEIVE_REPLY from=1 lamport=35 local=37 replies=2/2
2026/11/08 21:48:43.448653 [Node 3] event=ENTER_CS lamport=37
2026/11/08 21:48:43.970199 [Node 3] event=RECEIVE_REQUEST from=1 lamport=36 local=38 state=HELD
2026/11/08 21:48:43.970199 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:48:44.449919 [Node 3] event=EXIT_CS lamport=38 deferred=1
2026/11/08 21:48:44.449919 [Node 3] event=SEND_REPLY to=1 lamport=39
2026/11/08 21:48:45.740213 [Node 3] event=RECEIVE_REQUEST from=2 lamport=39 local=40 state=RELEASED
2026/11/08 21:48:45.740213 [Node 3] event=SEND_REPLY to=2 lamport=41
2026/11/08 21:48:45.859283 [Node 3] event=REQUEST_CS lamport=42
2026/11/08 21:48:45.859283 [Node 3] event=SEND_REQUEST to=1 lamport=42
2026/11/08 21:48:45.859283 [Node 3] event=SEND_REQUEST to=2 lamport=42
2026/11/08 21:48:45.859942 [Node 3] event=RECEIVE_REPLY from=1 lamport=44 local=45 replies=1/2
2026/11/08 21:48:46.741925 [Node 3] event=RECEIVE_REPLY from=2 lamport=46 local=47 replies=2/2
2026/11/08 21:48:46.741925 [Node 3] event=ENTER_CS lamport=47
2026/11/08 21:48:47.095362 [Node 3] event=RECEIVE_REQUEST from=1 lamport=45 local=48 state=HELD
2026/11/08 21:48:47.095362 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:48:47.742785 [Node 3] event=EXIT_CS lamport=48 deferred=1
2026/11/08 21:48:47.742785 [Node 3] event=SEND_REPLY to=1 lamport=49
2026/11/08 21:48:48.213898 [Node 3] event=RECEIVE_REQUEST from=2 lamport=49 local=50 state=RELEASED
2026/11/08 21:48:48.213898 [Node 3] event=SEND_REPLY to=2 lamport=51
2026/11/08 21:48:49.449148 [Node 3] event=REQUEST_CS lamport=52
2026/11/08 21:48:49.449148 [Node 3] event=SEND_REQUEST to=1 lamport=52
2026/11/08 21:48:49.449148 [Node 3] event=SEND_REQUEST to=2 lamport=52
2026/11/08 21:48:49.449148 [Node 3] event=RECEIVE_REPLY from=1 lamport=54 local=55 replies=1/2
2026/11/08 21:48:49.745332 [Node 3] event=RECEIVE_REPLY from=2 lamport=55 local=56 replies=2/2
2026/11/08 21:48:49.745332 [Node 3] event=ENTER_CS lamport=56
2026/11/08 21:48:50.746880 [Node 3] event=EXIT_CS lamport=56 deferred=0
2026/11/08 21:48:51.229880 [Node 3] event=RECEIVE_REQUEST from=1 lamport=55 local=57 state=RELEASED
2026/11/08 21:48:51.229880 [Node 3] event=SEND_REPLY to=1 lamport=58
2026/11/08 21:48:52.287606 [Node 3] event=REQUEST_CS lamport=59
2026/11/08 21:48:52.287606 [Node 3] event=SEND_REQUEST to=1 lamport=59
2026/11/08 21:48:52.287606 [Node 3] event=SEND_REQUEST to=2 lamport=59
2026/11/08 21:48:52.287606 [Node 3] event=RECEIVE_REPLY from=1 lamport=61 local=62 replies=1/2
2026/11/08 21:48:52.287606 [Node 3] event=RECEIVE_REPLY from=2 lamport=61 local=63 replies=2/2
2026/11/08 21:48:52.287606 [Node 3] event=ENTER_CS lamport=63
2026/11/08 21:48:53.076607 [Node 3] event=RECEIVE_REQUEST from=2 lamport=62 local=64 state=HELD
2026/11/08 21:48:53.076607 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:48:53.288674 [Node 3] event=EXIT_CS lamport=64 deferred=1
2026/11/08 21:48:53.288674 [Node 3] event=SEND_REPLY to=2 lamport=65
2026/11/08 21:48:53.866702 [Node 3] event=RECEIVE_REQUEST from=1 lamport=65 local=66 state=RELEASED
2026/11/08 21:48:53.866702 [Node 3] event=SEND_REPLY to=1 lamport=67
2026/11/08 21:48:55.402492 [Node 3] event=REQUEST_CS lamport=68
2026/11/08 21:48:55.402492 [Node 3] event=SEND_REQUEST to=1 lamport=68
2026/11/08 21:48:55.402492 [Node 3] event=SEND_REQUEST to=2 lamport=68
2026/11/08 21:48:55.402492 [Node 3] event=RECEIVE_REPLY from=1 lamport=71 local=72 replies=1/2
2026/11/08 21:48:55.402492 [Node 3] event=RECEIVE_REPLY from=2 lamport=70 local=73 replies=2/2
2026/11/08 21:48:55.402492 [Node 3] event=ENTER_CS lamport=73
2026/11/08 21:48:56.403339 [Node 3] event=EXIT_CS lamport=73 deferred=0
2026/11/08 21:48:57.588557 [Node 3] event=RECEIVE_REQUEST from=1 lamport=72 local=74 state=RELEASED
2026/11/08 21:48:57.588557 [Node 3] event=SEND_REPLY to=1 lamport=75
2026/11/08 21:48:58.099310 [Node 3] event=RECEIVE_REQUEST from=2 lamport=75 local=76 state=RELEASED
2026/11/08 21:48:58.099310 [Node 3] event=SEND_REPLY to=2 lamport=77
2026/11/08 21:48:59.650893 [Node 3] event=REQUEST_CS lamport=78
2026/11/08 21:48:59.650893 [Node 3] event=SEND_REQUEST to=1 lamport=78
2026/11/08 21:48:59.650893 [Node 3] event=SEND_REQUEST to=2 lamport=78
2026/11/08 21:48:59.650893 [Node 3] event=RECEIVE_REPLY from=2 lamport=82 local=83 replies=1/2
2026/11/08 21:48:59.650893 [Node 3] event=RECEIVE_REPLY from=1 lamport=81 local=84 replies=2/2
2026/11/08 21:48:59.650893 [Node 3] event=ENTER_CS lamport=84
2026/11/08 21:48:59.682893 [Node 3] event=RECEIVE_REQUEST from=1 lamport=82 local=85 state=HELD
2026/11/08 21:48:59.682893 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:49:00.651894 [Node 3] event=EXIT_CS lamport=85 deferred=1
2026/11/08 21:49:00.651894 [Node 3] event=SEND_REPLY to=1 lamport=86
2026/11/08 21:49:01.614867 [Node 3] event=RECEIVE_REQUEST from=2 lamport=85 local=87 state=RELEASED
2026/11/08 21:49:01.614867 [Node 3] event=SEND_REPLY to=2 lamport=88
2026/11/08 21:49:04.278450 [Node 3] event=REQUEST_CS lamport=89
2026/11/08 21:49:04.278450 [Node 3] event=SEND_REQUEST to=1 lamport=89
2026/11/08 21:49:04.278450 [Node 3] event=SEND_REQUEST to=2 lamport=89
2026/11/08 21:49:04.278450 [Node 3] event=RECEIVE_REPLY from=1 lamport=91 local=92 replies=1/2
2026/11/08 21:49:04.279452 [Node 3] event=RECEIVE_REPLY from=2 lamport=92 local=93 replies=2/2
2026/11/08 21:49:04.279452 [Node 3] event=ENTER_CS lamport=93
2026/11/08 21:49:04.588451 [Node 3] event=RECEIVE_REQUEST from=1 lamport=92 local=94 state=HELD
2026/11/08 21:49:04.588451 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:49:05.280415 [Node 3] event=EXIT_CS lamport=94 deferred=1
2026/11/08 21:49:05.280415 [Node 3] event=SEND_REPLY to=1 lamport=95
2026/11/08 21:49:05.786484 [Node 3] event=RECEIVE_REQUEST from=2 lamport=95 local=96 state=RELEASED
2026/11/08 21:49:05.786988 [Node 3] event=SEND_REPLY to=2 lamport=97
2026/11/08 21:49:06.726069 [Node 3] event=REQUEST_CS lamport=98
2026/11/08 21:49:06.726118 [Node 3] event=SEND_REQUEST to=1 lamport=98
2026/11/08 21:49:06.726118 [Node 3] event=SEND_REQUEST to=2 lamport=98
2026/11/08 21:49:06.726650 [Node 3] event=RECEIVE_REPLY from=1 lamport=100 local=101 replies=1/2
2026/11/08 21:49:07.282112 [Node 3] event=RECEIVE_REPLY from=2 lamport=101 local=102 replies=2/2
2026/11/08 21:49:07.282112 [Node 3] event=ENTER_CS lamport=102
2026/11/08 21:49:08.282626 [Node 3] event=EXIT_CS lamport=102 deferred=0
2026/11/08 21:49:08.570309 [Node 3] event=RECEIVE_REQUEST from=1 lamport=101 local=103 state=RELEASED
2026/11/08 21:49:08.570835 [Node 3] event=SEND_REPLY to=1 lamport=104
2026/11/08 21:49:09.828145 [Node 3] event=REQUEST_CS lamport=105
2026/11/08 21:49:09.828145 [Node 3] event=SEND_REQUEST to=1 lamport=105
2026/11/08 21:49:09.828145 [Node 3] event=SEND_REQUEST to=2 lamport=105
2026/11/08 21:49:09.828743 [Node 3] event=RECEIVE_REPLY from=2 lamport=107 local=108 replies=1/2
2026/11/08 21:49:09.828743 [Node 3] event=RECEIVE_REPLY from=1 lamport=108 local=109 replies=2/2
2026/11/08 21:49:09.828743 [Node 3] event=ENTER_CS lamport=109
2026/11/08 21:49:10.828784 [Node 3] event=EXIT_CS lamport=109 deferred=0
2026/11/08 21:49:10.941783 [Node 3] event=RECEIVE_REQUEST from=2 lamport=108 local=110 state=RELEASED
2026/11/08 21:49:10.941783 [Node 3] event=SEND_REPLY to=2 lamport=111
2026/11/08 21:49:11.700783 [Node 3] event=RECEIVE_REQUEST from=1 lamport=111 local=112 state=RELEASED
2026/11/08 21:49:11.700783 [Node 3] event=SEND_REPLY to=1 lamport=113
2026/11/08 21:49:13.886785 [Node 3] event=REQUEST_CS lamport=114
2026/11/08 21:49:13.886785 [Node 3] event=SEND_REQUEST to=1 lamport=114
2026/11/08 21:49:13.886785 [Node 3] event=SEND_REQUEST to=2 lamport=114
2026/11/08 21:49:13.886785 [Node 3] event=RECEIVE_REPLY from=1 lamport=118 local=119 replies=1/2
2026/11/08 21:49:13.887784 [Node 3] event=RECEIVE_REPLY from=2 lamport=117 local=120 replies=2/2
2026/11/08 21:49:13.887784 [Node 3] event=ENTER_CS lamport=120
2026/11/08 21:49:14.228215 [Node 3] event=RECEIVE_REQUEST from=2 lamport=118 local=121 state=HELD
2026/11/08 21:49:14.228215 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:49:14.473214 [Node 3] event=RECEIVE_REQUEST from=1 lamport=121 local=122 state=HELD
2026/11/08 21:49:14.473214 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:49:14.889214 [Node 3] event=EXIT_CS lamport=122 deferred=2
2026/11/08 21:49:14.889214 [Node 3] event=SEND_REPLY to=1 lamport=123
2026/11/08 21:49:14.889214 [Node 3] event=SEND_REPLY to=2 lamport=124
2026/11/08 21:49:17.160999 [Node 3] event=REQUEST_CS lamport=125
2026/11/08 21:49:17.160999 [Node 3] event=SEND_REQUEST to=1 lamport=125
2026/11/08 21:49:17.160999 [Node 3] event=SEND_REQUEST to=2 lamport=125
2026/11/08 21:49:17.161578 [Node 3] event=RECEIVE_REPLY from=2 lamport=128 local=129 replies=1/2
2026/11/08 21:49:17.161578 [Node 3] event=RECEIVE_REPLY from=1 lamport=129 local=130 replies=2/2
2026/11/08 21:49:17.161578 [Node 3] event=ENTER_CS lamport=130
2026/11/08 21:49:18.162291 [Node 3] event=EXIT_CS lamport=130 deferred=0
2026/11/08 21:49:19.579330 [Node 3] event=RECEIVE_REQUEST from=1 lamport=130 local=131 state=RELEASED
2026/11/08 21:49:19.579330 [Node 3] event=SEND_REPLY to=1 lamport=132
2026/11/08 21:49:19.772539 [Node 3] event=RECEIVE_REQUEST from=2 lamport=133 local=134 state=RELEASED
2026/11/08 21:49:19.772539 [Node 3] event=SEND_REPLY to=2 lamport=135
2026/11/08 21:49:20.948917 [Node 3] event=REQUEST_CS lamport=136
2026/11/08 21:49:20.948917 [Node 3] event=SEND_REQUEST to=1 lamport=136
2026/11/08 21:49:20.948917 [Node 3] event=SEND_REQUEST to=2 lamport=136
2026/11/08 21:49:20.948917 [Node 3] event=RECEIVE_REPLY from=1 lamport=138 local=139 replies=1/2
2026/11/08 21:49:21.582207 [Node 3] event=RECEIVE_REPLY from=2 lamport=139 local=140 replies=2/2
2026/11/08 21:49:21.582207 [Node 3] event=ENTER_CS lamport=140
2026/11/08 21:49:22.583837 [Node 3] event=EXIT_CS lamport=140 deferred=0
2026/11/08 21:49:23.396838 [Node 3] event=RECEIVE_REQUEST from=1 lamport=139 local=141 state=RELEASED
2026/11/08 21:49:23.397455 [Node 3] event=SEND_REPLY to=1 lamport=142
2026/11/08 21:49:23.882497 [Node 3] event=REQUEST_CS lamport=143
2026/11/08 21:49:23.882563 [Node 3] event=SEND_REQUEST to=1 lamport=143
2026/11/08 21:49:23.882563 [Node 3] event=SEND_REQUEST to=2 lamport=143
2026/11/08 21:49:23.883087 [Node 3] event=RECEIVE_REPLY from=2 lamport=145 local=146 replies=1/2
2026/11/08 21:49:23.897211 [Node 3] event=RECEIVE_REQUEST from=2 lamport=146 local=147 state=WANTED
2026/11/08 21:49:23.897211 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:49:24.398741 [Node 3] event=RECEIVE_REPLY from=1 lamport=149 local=150 replies=2/2
2026/11/08 21:49:24.398741 [Node 3] event=ENTER_CS lamport=150
2026/11/08 21:49:25.399080 [Node 3] event=EXIT_CS lamport=150 deferred=1
2026/11/08 21:49:25.399080 [Node 3] event=SEND_REPLY to=2 lamport=151
2026/11/08 21:49:27.834612 [Node 3] event=REQUEST_CS lamport=152
2026/11/08 21:49:27.834612 [Node 3] event=SEND_REQUEST to=1 lamport=152
2026/11/08 21:49:27.834612 [Node 3] event=SEND_REQUEST to=2 lamport=152
2026/11/08 21:49:27.835329 [Node 3] event=RECEIVE_REPLY from=2 lamport=154 local=155 replies=1/2
2026/11/08 21:49:27.835329 [Node 3] event=RECEIVE_REPLY from=1 lamport=154 local=156 replies=2/2
2026/11/08 21:49:27.835329 [Node 3] event=ENTER_CS lamport=156
2026/11/08 21:49:28.170386 [Node 3] event=RECEIVE_REQUEST from=1 lamport=155 local=157 state=HELD
2026/11/08 21:49:28.170386 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:49:28.836313 [Node 3] event=EXIT_CS lamport=157 deferred=1
2026/11/08 21:49:28.836313 [Node 3] event=SEND_REPLY to=1 lamport=158
2026/11/08 21:49:29.309702 [Node 3] event=RECEIVE_REQUEST from=2 lamport=158 local=159 state=RELEASED
2026/11/08 21:49:29.309702 [Node 3] event=SEND_REPLY to=2 lamport=160
2026/11/08 21:49:31.660819 [Node 3] event=REQUEST_CS lamport=161
2026/11/08 21:49:31.660819 [Node 3] event=SEND_REQUEST to=2 lamport=161
2026/11/08 21:49:31.660819 [Node 3] event=SEND_REQUEST to=1 lamport=161
2026/11/08 21:49:31.661682 [Node 3] event=RECEIVE_REPLY from=1 lamport=163 local=164 replies=1/2
2026/11/08 21:49:31.661682 [Node 3] event=RECEIVE_REPLY from=2 lamport=164 local=165 replies=2/2
2026/11/08 21:49:31.661682 [Node 3] event=ENTER_CS lamport=165
2026/11/08 21:49:32.662365 [Node 3] event=EXIT_CS lamport=165 deferred=0
2026/11/08 21:49:33.609940 [Node 3] event=RECEIVE_REQUEST from=1 lamport=164 local=166 state=RELEASED
2026/11/08 21:49:33.609940 [Node 3] event=SEND_REPLY to=1 lamport=167
2026/11/08 21:49:34.316255 [Node 3] event=RECEIVE_REQUEST from=2 lamport=167 local=168 state=RELEASED
2026/11/08 21:49:34.316255 [Node 3] event=SEND_REPLY to=2 lamport=169
2026/11/08 21:49:35.650781 [Node 3] event=REQUEST_CS lamport=170
2026/11/08 21:49:35.650781 [Node 3] event=SEND_REQUEST to=1 lamport=170
2026/11/08 21:49:35.650781 [Node 3] event=SEND_REQUEST to=2 lamport=170
2026/11/08 21:49:35.650781 [Node 3] event=RECEIVE_REPLY from=2 lamport=173 local=174 replies=1/2
2026/11/08 21:49:35.650781 [Node 3] event=RECEIVE_REPLY from=1 lamport=172 local=175 replies=2/2
2026/11/08 21:49:35.650781 [Node 3] event=ENTER_CS lamport=175
2026/11/08 21:49:36.405970 [Node 3] event=RECEIVE_REQUEST from=1 lamport=173 local=176 state=HELD
2026/11/08 21:49:36.405970 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:49:36.652969 [Node 3] event=EXIT_CS lamport=176 deferred=1
2026/11/08 21:49:36.652969 [Node 3] event=SEND_REPLY to=1 lamport=177
2026/11/08 21:49:38.729037 [Node 3] event=RECEIVE_REQUEST from=1 lamport=179 local=180 state=RELEASED
2026/11/08 21:49:38.729037 [Node 3] event=SEND_REPLY to=1 lamport=181
2026/11/08 21:49:39.228496 [Node 3] event=RECEIVE_REQUEST from=2 lamport=182 local=183 state=RELEASED
2026/11/08 21:49:39.228496 [Node 3] event=SEND_REPLY to=2 lamport=184
2026/11/08 21:49:39.780758 [Node 3] event=REQUEST_CS lamport=185
2026/11/08 21:49:39.780758 [Node 3] event=SEND_REQUEST to=1 lamport=185
2026/11/08 21:49:39.780758 [Node 3] event=SEND_REQUEST to=2 lamport=185
2026/11/08 21:49:39.780758 [Node 3] event=RECEIVE_REPLY from=1 lamport=187 local=188 replies=1/2
2026/11/08 21:49:40.730288 [Node 3] event=RECEIVE_REPLY from=2 lamport=188 local=189 replies=2/2
2026/11/08 21:49:40.730288 [Node 3] event=ENTER_CS lamport=189
2026/11/08 21:49:41.730548 [Node 3] event=EXIT_CS lamport=189 deferred=0
2026/11/08 21:49:42.308091 [Node 3] event=RECEIVE_REQUEST from=1 lamport=188 local=190 state=RELEASED
2026/11/08 21:49:42.308091 [Node 3] event=SEND_REPLY to=1 lamport=191
2026/11/08 21:49:43.605729 [Node 3] event=RECEIVE_REQUEST from=2 lamport=191 local=192 state=RELEASED
2026/11/08 21:49:43.605729 [Node 3] event=SEND_REPLY to=2 lamport=193
2026/11/08 21:49:44.261728 [Node 3] event=REQUEST_CS lamport=194
2026/11/08 21:49:44.261728 [Node 3] event=SEND_REQUEST to=1 lamport=194
2026/11/08 21:49:44.261728 [Node 3] event=SEND_REQUEST to=2 lamport=194
2026/11/08 21:49:44.261728 [Node 3] event=RECEIVE_REPLY from=1 lamport=196 local=197 replies=1/2
2026/11/08 21:49:44.607729 [Node 3] event=RECEIVE_REPLY from=2 lamport=198 local=199 replies=2/2
2026/11/08 21:49:44.607729 [Node 3] event=ENTER_CS lamport=199
2026/11/08 21:49:44.914259 [Node 3] event=RECEIVE_REQUEST from=1 lamport=197 local=200 state=HELD
2026/11/08 21:49:44.914259 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:49:45.608261 [Node 3] event=EXIT_CS lamport=200 deferred=1
2026/11/08 21:49:45.608261 [Node 3] event=SEND_REPLY to=1 lamport=201
2026/11/08 21:49:46.695382 [Node 3] event=REQUEST_CS lamport=202
2026/11/08 21:49:46.695382 [Node 3] event=SEND_REQUEST to=2 lamport=202
2026/11/08 21:49:46.695382 [Node 3] event=SEND_REQUEST to=1 lamport=202
2026/11/08 21:49:46.695382 [Node 3] event=RECEIVE_REPLY from=2 lamport=204 local=205 replies=1/2
2026/11/08 21:49:46.695382 [Node 3] event=RECEIVE_REPLY from=1 lamport=204 local=206 replies=2/2
2026/11/08 21:49:46.696395 [Node 3] event=ENTER_CS lamport=206
2026/11/08 21:49:47.403553 [Node 3] event=RECEIVE_REQUEST from=2 lamport=205 local=207 state=HELD
2026/11/08 21:49:47.403553 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:49:47.696634 [Node 3] event=EXIT_CS lamport=207 deferred=1
2026/11/08 21:49:47.696634 [Node 3] event=SEND_REPLY to=2 lamport=208
2026/11/08 21:49:47.912668 [Node 3] event=RECEIVE_REQUEST from=1 lamport=208 local=209 state=RELEASED
2026/11/08 21:49:47.912668 [Node 3] event=SEND_REPLY to=1 lamport=210
2026/11/08 21:49:50.583648 [Node 3] event=REQUEST_CS lamport=211
2026/11/08 21:49:50.583648 [Node 3] event=SEND_REQUEST to=1 lamport=211
2026/11/08 21:49:50.583648 [Node 3] event=SEND_REQUEST to=2 lamport=211
2026/11/08 21:49:50.583648 [Node 3] event=RECEIVE_REPLY from=1 lamport=214 local=215 replies=1/2
2026/11/08 21:49:50.583648 [Node 3] event=RECEIVE_REPLY from=2 lamport=213 local=216 replies=2/2
2026/11/08 21:49:50.584650 [Node 3] event=ENTER_CS lamport=216
2026/11/08 21:49:51.439324 [Node 3] event=RECEIVE_REQUEST from=2 lamport=214 local=217 state=HELD
2026/11/08 21:49:51.439324 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:49:51.585626 [Node 3] event=EXIT_CS lamport=217 deferred=1
2026/11/08 21:49:51.585678 [Node 3] event=SEND_REPLY to=2 lamport=218
2026/11/08 21:49:52.155451 [Node 3] event=RECEIVE_REQUEST from=1 lamport=217 local=219 state=RELEASED
2026/11/08 21:49:52.155451 [Node 3] event=SEND_REPLY to=1 lamport=220
2026/11/08 21:49:54.621989 [Node 3] event=RECEIVE_REQUEST from=2 lamport=222 local=223 state=RELEASED
2026/11/08 21:49:54.621989 [Node 3] event=SEND_REPLY to=2 lamport=224
2026/11/08 21:49:54.935990 [Node 3] event=REQUEST_CS lamport=225
2026/11/08 21:49:54.935990 [Node 3] event=SEND_REQUEST to=1 lamport=225
2026/11/08 21:49:54.935990 [Node 3] event=SEND_REQUEST to=2 lamport=225
2026/11/08 21:49:54.937060 [Node 3] event=RECEIVE_REPLY from=1 lamport=227 local=228 replies=1/2
2026/11/08 21:49:55.623576 [Node 3] event=RECEIVE_REPLY from=2 lamport=228 local=229 replies=2/2
2026/11/08 21:49:55.623576 [Node 3] event=ENTER_CS lamport=229
2026/11/08 21:49:56.623969 [Node 3] event=EXIT_CS lamport=229 deferred=0
2026/11/08 21:49:57.227981 [Node 3] event=RECEIVE_REQUEST from=1 lamport=228 local=230 state=RELEASED
2026/11/08 21:49:57.227981 [Node 3] event=SEND_REPLY to=1 lamport=231
2026/11/08 21:49:58.135041 [Node 3] event=RECEIVE_REQUEST from=2 lamport=231 local=232 state=RELEASED
2026/11/08 21:49:58.135041 [Node 3] event=SEND_REPLY to=2 lamport=233
2026/11/08 21:49:59.739551 [Node 3] event=REQUEST_CS lamport=234
2026/11/08 21:49:59.739551 [Node 3] event=SEND_REQUEST to=1 lamport=234
2026/11/08 21:49:59.739551 [Node 3] event=SEND_REQUEST to=2 lamport=234
2026/11/08 21:49:59.739551 [Node 3] event=RECEIVE_REPLY from=2 lamport=238 local=239 replies=1/2
2026/11/08 21:49:59.740552 [Node 3] event=RECEIVE_REPLY from=1 lamport=237 local=240 replies=2/2
2026/11/08 21:49:59.740552 [Node 3] event=ENTER_CS lamport=240
2026/11/08 21:50:00.741551 [Node 3] event=EXIT_CS lamport=240 deferred=0
2026/11/08 21:50:01.475551 [Node 3] event=RECEIVE_REQUEST from=1 lamport=238 local=241 state=RELEASED
2026/11/08 21:50:01.475551 [Node 3] event=SEND_REPLY to=1 lamport=242
2026/11/08 21:50:02.384094 [Node 3] event=REQUEST_CS lamport=243
2026/11/08 21:50:02.384094 [Node 3] event=SEND_REQUEST to=1 lamport=243
2026/11/08 21:50:02.384094 [Node 3] event=SEND_REQUEST to=2 lamport=243
2026/11/08 21:50:02.385034 [Node 3] event=RECEIVE_REPLY from=2 lamport=245 local=246 replies=1/2
2026/11/08 21:50:02.477253 [Node 3] event=RECEIVE_REPLY from=1 lamport=245 local=247 replies=2/2
2026/11/08 21:50:02.477253 [Node 3] event=ENTER_CS lamport=247
2026/11/08 21:50:02.756389 [Node 3] event=RECEIVE_REQUEST from=2 lamport=246 local=248 state=HELD
2026/11/08 21:50:02.756389 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:03.478307 [Node 3] event=EXIT_CS lamport=248 deferred=1
2026/11/08 21:50:03.478307 [Node 3] event=SEND_REPLY to=2 lamport=249
2026/11/08 21:50:03.561384 [Node 3] event=RECEIVE_REQUEST from=1 lamport=249 local=250 state=RELEASED
2026/11/08 21:50:03.561384 [Node 3] event=SEND_REPLY to=1 lamport=251
2026/11/08 21:50:06.191544 [Node 3] event=REQUEST_CS lamport=252
2026/11/08 21:50:06.191544 [Node 3] event=SEND_REQUEST to=1 lamport=252
2026/11/08 21:50:06.191544 [Node 3] event=SEND_REQUEST to=2 lamport=252
2026/11/08 21:50:06.192546 [Node 3] event=RECEIVE_REPLY from=2 lamport=254 local=255 replies=1/2
2026/11/08 21:50:06.192546 [Node 3] event=RECEIVE_REPLY from=1 lamport=255 local=256 replies=2/2
2026/11/08 21:50:06.192546 [Node 3] event=ENTER_CS lamport=256
2026/11/08 21:50:07.034858 [Node 3] event=RECEIVE_REQUEST from=2 lamport=255 local=257 state=HELD
2026/11/08 21:50:07.034858 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:07.192859 [Node 3] event=EXIT_CS lamport=257 deferred=1
2026/11/08 21:50:07.192859 [Node 3] event=SEND_REPLY to=2 lamport=258
2026/11/08 21:50:09.472858 [Node 3] event=RECEIVE_REQUEST from=1 lamport=258 local=259 state=RELEASED
2026/11/08 21:50:09.472858 [Node 3] event=SEND_REPLY to=1 lamport=260
2026/11/08 21:50:11.005396 [Node 3] event=REQUEST_CS lamport=261
2026/11/08 21:50:11.005396 [Node 3] event=SEND_REQUEST to=1 lamport=261
2026/11/08 21:50:11.005396 [Node 3] event=SEND_REQUEST to=2 lamport=261
2026/11/08 21:50:11.006396 [Node 3] event=RECEIVE_REPLY from=1 lamport=265 local=266 replies=1/2
2026/11/08 21:50:11.006396 [Node 3] event=RECEIVE_REPLY from=2 lamport=263 local=267 replies=2/2
2026/11/08 21:50:11.006396 [Node 3] event=ENTER_CS lamport=267
2026/11/08 21:50:11.893394 [Node 3] event=RECEIVE_REQUEST from=2 lamport=264 local=268 state=HELD
2026/11/08 21:50:11.893394 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:12.007941 [Node 3] event=EXIT_CS lamport=268 deferred=1
2026/11/08 21:50:12.007941 [Node 3] event=SEND_REPLY to=2 lamport=269
2026/11/08 21:50:13.579938 [Node 3] event=REQUEST_CS lamport=270
2026/11/08 21:50:13.579938 [Node 3] event=SEND_REQUEST to=1 lamport=270
2026/11/08 21:50:13.579938 [Node 3] event=SEND_REQUEST to=2 lamport=270
2026/11/08 21:50:13.580941 [Node 3] event=RECEIVE_REPLY from=2 lamport=272 local=273 replies=1/2
2026/11/08 21:50:13.580941 [Node 3] event=RECEIVE_REPLY from=1 lamport=272 local=274 replies=2/2
2026/11/08 21:50:13.580941 [Node 3] event=ENTER_CS lamport=274
2026/11/08 21:50:14.056938 [Node 3] event=RECEIVE_REQUEST from=1 lamport=273 local=275 state=HELD
2026/11/08 21:50:14.056938 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:50:14.581496 [Node 3] event=EXIT_CS lamport=275 deferred=1
2026/11/08 21:50:14.581496 [Node 3] event=SEND_REPLY to=1 lamport=276
2026/11/08 21:50:15.910069 [Node 3] event=RECEIVE_REQUEST from=2 lamport=276 local=277 state=RELEASED
2026/11/08 21:50:15.910069 [Node 3] event=SEND_REPLY to=2 lamport=278
2026/11/08 21:50:16.652068 [Node 3] event=RECEIVE_REQUEST from=1 lamport=280 local=281 state=RELEASED
2026/11/08 21:50:16.652068 [Node 3] event=SEND_REPLY to=1 lamport=282
2026/11/08 21:50:18.016069 [Node 3] event=RECEIVE_REQUEST from=2 lamport=283 local=284 state=RELEASED
2026/11/08 21:50:18.016069 [Node 3] event=SEND_REPLY to=2 lamport=285
2026/11/08 21:50:18.124070 [Node 3] event=REQUEST_CS lamport=286
2026/11/08 21:50:18.124070 [Node 3] event=SEND_REQUEST to=1 lamport=286
2026/11/08 21:50:18.124070 [Node 3] event=SEND_REQUEST to=2 lamport=286
2026/11/08 21:50:18.124070 [Node 3] event=RECEIVE_REPLY from=1 lamport=288 local=289 replies=1/2
2026/11/08 21:50:19.018068 [Node 3] event=RECEIVE_REPLY from=2 lamport=290 local=291 replies=2/2
2026/11/08 21:50:19.018068 [Node 3] event=ENTER_CS lamport=291
2026/11/08 21:50:20.019051 [Node 3] event=EXIT_CS lamport=291 deferred=0
2026/11/08 21:50:21.542646 [Node 3] event=RECEIVE_REQUEST from=1 lamport=289 local=292 state=RELEASED
2026/11/08 21:50:21.542646 [Node 3] event=SEND_REPLY to=1 lamport=293
2026/11/08 21:50:21.758495 [Node 3] event=REQUEST_CS lamport=294
2026/11/08 21:50:21.758495 [Node 3] event=SEND_REQUEST to=1 lamport=294
2026/11/08 21:50:21.758495 [Node 3] event=SEND_REQUEST to=2 lamport=294
2026/11/08 21:50:21.759497 [Node 3] event=RECEIVE_REPLY from=2 lamport=296 local=297 replies=1/2
2026/11/08 21:50:22.380813 [Node 3] event=RECEIVE_REQUEST from=2 lamport=297 local=298 state=WANTED
2026/11/08 21:50:22.380813 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:22.543813 [Node 3] event=RECEIVE_REPLY from=1 lamport=300 local=301 replies=2/2
2026/11/08 21:50:22.543813 [Node 3] event=ENTER_CS lamport=301
2026/11/08 21:50:23.544858 [Node 3] event=EXIT_CS lamport=301 deferred=1
2026/11/08 21:50:23.544858 [Node 3] event=SEND_REPLY to=2 lamport=302
2026/11/08 21:50:23.561760 [Node 3] event=RECEIVE_REQUEST from=1 lamport=301 local=303 state=RELEASED
2026/11/08 21:50:23.561760 [Node 3] event=SEND_REPLY to=1 lamport=304
2026/11/08 21:50:26.408605 [Node 3] event=REQUEST_CS lamport=305
2026/11/08 21:50:26.408605 [Node 3] event=SEND_REQUEST to=1 lamport=305
2026/11/08 21:50:26.408605 [Node 3] event=SEND_REQUEST to=2 lamport=305
2026/11/08 21:50:26.408605 [Node 3] event=RECEIVE_REPLY from=2 lamport=307 local=308 replies=1/2
2026/11/08 21:50:26.408605 [Node 3] event=RECEIVE_REPLY from=1 lamport=308 local=309 replies=2/2
2026/11/08 21:50:26.408605 [Node 3] event=ENTER_CS lamport=309
2026/11/08 21:50:26.701148 [Node 3] event=RECEIVE_REQUEST from=2 lamport=308 local=310 state=HELD
2026/11/08 21:50:26.701148 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:27.409709 [Node 3] event=EXIT_CS lamport=310 deferred=1
2026/11/08 21:50:27.409709 [Node 3] event=SEND_REPLY to=2 lamport=311
2026/11/08 21:50:28.363709 [Node 3] event=RECEIVE_REQUEST from=1 lamport=311 local=312 state=RELEASED
2026/11/08 21:50:28.363709 [Node 3] event=SEND_REPLY to=1 lamport=313
2026/11/08 21:50:28.448709 [Node 3] event=REQUEST_CS lamport=314
2026/11/08 21:50:28.448709 [Node 3] event=SEND_REQUEST to=1 lamport=314
2026/11/08 21:50:28.448709 [Node 3] event=SEND_REQUEST to=2 lamport=314
2026/11/08 21:50:28.448709 [Node 3] event=RECEIVE_REPLY from=2 lamport=316 local=317 replies=1/2
2026/11/08 21:50:29.411345 [Node 3] event=RECEIVE_REPLY from=1 lamport=317 local=318 replies=2/2
2026/11/08 21:50:29.411345 [Node 3] event=ENTER_CS lamport=318
2026/11/08 21:50:29.934347 [Node 3] event=RECEIVE_REQUEST from=2 lamport=317 local=319 state=HELD
2026/11/08 21:50:29.934347 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:30.411490 [Node 3] event=EXIT_CS lamport=319 deferred=1
2026/11/08 21:50:30.411525 [Node 3] event=SEND_REPLY to=2 lamport=320
2026/11/08 21:50:30.714556 [Node 3] event=RECEIVE_REQUEST from=1 lamport=320 local=321 state=RELEASED
2026/11/08 21:50:30.714556 [Node 3] event=SEND_REPLY to=1 lamport=322
2026/11/08 21:50:32.501116 [Node 3] event=RECEIVE_REQUEST from=2 lamport=324 local=325 state=RELEASED
2026/11/08 21:50:32.501116 [Node 3] event=SEND_REPLY to=2 lamport=326
2026/11/08 21:50:32.508253 [Node 3] event=REQUEST_CS lamport=327
2026/11/08 21:50:32.508253 [Node 3] event=SEND_REQUEST to=1 lamport=327
2026/11/08 21:50:32.508253 [Node 3] event=SEND_REQUEST to=2 lamport=327
2026/11/08 21:50:32.508813 [Node 3] event=RECEIVE_REPLY from=1 lamport=329 local=330 replies=1/2
2026/11/08 21:50:33.502464 [Node 3] event=RECEIVE_REPLY from=2 lamport=330 local=331 replies=2/2
2026/11/08 21:50:33.502464 [Node 3] event=ENTER_CS lamport=331
2026/11/08 21:50:34.502864 [Node 3] event=EXIT_CS lamport=331 deferred=0
2026/11/08 21:50:34.506439 [Node 3] event=RECEIVE_REQUEST from=1 lamport=330 local=332 state=RELEASED
2026/11/08 21:50:34.506439 [Node 3] event=SEND_REPLY to=1 lamport=333
2026/11/08 21:50:35.699017 [Node 3] event=REQUEST_CS lamport=334
2026/11/08 21:50:35.699017 [Node 3] event=SEND_REQUEST to=1 lamport=334
2026/11/08 21:50:35.699017 [Node 3] event=SEND_REQUEST to=2 lamport=334
2026/11/08 21:50:35.699662 [Node 3] event=RECEIVE_REPLY from=2 lamport=336 local=337 replies=1/2
2026/11/08 21:50:35.699662 [Node 3] event=RECEIVE_REPLY from=1 lamport=336 local=338 replies=2/2
2026/11/08 21:50:35.699662 [Node 3] event=ENTER_CS lamport=338
2026/11/08 21:50:36.498921 [Node 3] event=RECEIVE_REQUEST from=2 lamport=337 local=339 state=HELD
2026/11/08 21:50:36.499426 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:36.699697 [Node 3] event=EXIT_CS lamport=339 deferred=1
2026/11/08 21:50:36.699697 [Node 3] event=SEND_REPLY to=2 lamport=340
2026/11/08 21:50:36.931396 [Node 3] event=RECEIVE_REQUEST from=1 lamport=340 local=341 state=RELEASED
2026/11/08 21:50:36.931396 [Node 3] event=SEND_REPLY to=1 lamport=342
2026/11/08 21:50:39.572932 [Node 3] event=REQUEST_CS lamport=343
2026/11/08 21:50:39.572932 [Node 3] event=SEND_REQUEST to=2 lamport=343
2026/11/08 21:50:39.572932 [Node 3] event=SEND_REQUEST to=1 lamport=343
2026/11/08 21:50:39.573491 [Node 3] event=RECEIVE_REPLY from=2 lamport=345 local=346 replies=1/2
2026/11/08 21:50:39.573491 [Node 3] event=RECEIVE_REPLY from=1 lamport=346 local=347 replies=2/2
2026/11/08 21:50:39.573491 [Node 3] event=ENTER_CS lamport=347
2026/11/08 21:50:39.882745 [Node 3] event=RECEIVE_REQUEST from=2 lamport=346 local=348 state=HELD
2026/11/08 21:50:39.882745 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:40.574058 [Node 3] event=EXIT_CS lamport=348 deferred=1
2026/11/08 21:50:40.574058 [Node 3] event=SEND_REPLY to=2 lamport=349
2026/11/08 21:50:41.881401 [Node 3] event=RECEIVE_REQUEST from=1 lamport=349 local=350 state=RELEASED
2026/11/08 21:50:41.881401 [Node 3] event=SEND_REPLY to=1 lamport=351
2026/11/08 21:50:43.829352 [Node 3] event=REQUEST_CS lamport=352
2026/11/08 21:50:43.829352 [Node 3] event=SEND_REQUEST to=1 lamport=352
2026/11/08 21:50:43.829352 [Node 3] event=SEND_REQUEST to=2 lamport=352
2026/11/08 21:50:43.829960 [Node 3] event=RECEIVE_REPLY from=2 lamport=354 local=355 replies=1/2
2026/11/08 21:50:43.829960 [Node 3] event=RECEIVE_REPLY from=1 lamport=356 local=357 replies=2/2
2026/11/08 21:50:43.829960 [Node 3] event=ENTER_CS lamport=357
2026/11/08 21:50:44.588719 [Node 3] event=RECEIVE_REQUEST from=2 lamport=355 local=358 state=HELD
2026/11/08 21:50:44.588719 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:44.830412 [Node 3] event=EXIT_CS lamport=358 deferred=1
2026/11/08 21:50:44.830412 [Node 3] event=SEND_REPLY to=2 lamport=359
2026/11/08 21:50:45.122762 [Node 3] event=RECEIVE_REQUEST from=1 lamport=359 local=360 state=RELEASED
2026/11/08 21:50:45.122762 [Node 3] event=SEND_REPLY to=1 lamport=361
2026/11/08 21:50:47.287417 [Node 3] event=REQUEST_CS lamport=362
2026/11/08 21:50:47.287417 [Node 3] event=SEND_REQUEST to=1 lamport=362
2026/11/08 21:50:47.287417 [Node 3] event=SEND_REQUEST to=2 lamport=362
2026/11/08 21:50:47.287988 [Node 3] event=RECEIVE_REPLY from=2 lamport=364 local=365 replies=1/2
2026/11/08 21:50:47.287988 [Node 3] event=RECEIVE_REPLY from=1 lamport=365 local=366 replies=2/2
2026/11/08 21:50:47.287988 [Node 3] event=ENTER_CS lamport=366
2026/11/08 21:50:48.178033 [Node 3] event=RECEIVE_REQUEST from=2 lamport=365 local=367 state=HELD
2026/11/08 21:50:48.178033 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:48.288714 [Node 3] event=EXIT_CS lamport=367 deferred=1
2026/11/08 21:50:48.288714 [Node 3] event=SEND_REPLY to=2 lamport=368
2026/11/08 21:50:49.050316 [Node 3] event=RECEIVE_REQUEST from=1 lamport=368 local=369 state=RELEASED
2026/11/08 21:50:49.050316 [Node 3] event=SEND_REPLY to=1 lamport=370
2026/11/08 21:50:49.849200 [Node 3] event=REQUEST_CS lamport=371
2026/11/08 21:50:49.849200 [Node 3] event=SEND_REQUEST to=1 lamport=371
2026/11/08 21:50:49.849200 [Node 3] event=SEND_REQUEST to=2 lamport=371
2026/11/08 21:50:49.849200 [Node 3] event=RECEIVE_REPLY from=2 lamport=373 local=374 replies=1/2
2026/11/08 21:50:50.289501 [Node 3] event=RECEIVE_REPLY from=1 lamport=374 local=375 replies=2/2
2026/11/08 21:50:50.289501 [Node 3] event=ENTER_CS lamport=375
2026/11/08 21:50:50.317103 [Node 3] event=RECEIVE_REQUEST from=2 lamport=374 local=376 state=HELD
2026/11/08 21:50:50.317103 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:51.289863 [Node 3] event=EXIT_CS lamport=376 deferred=1
2026/11/08 21:50:51.289927 [Node 3] event=SEND_REPLY to=2 lamport=377
2026/11/08 21:50:52.644315 [Node 3] event=RECEIVE_REQUEST from=1 lamport=377 local=378 state=RELEASED
2026/11/08 21:50:52.644315 [Node 3] event=SEND_REPLY to=1 lamport=379
2026/11/08 21:50:53.374818 [Node 3] event=REQUEST_CS lamport=380
2026/11/08 21:50:53.374818 [Node 3] event=SEND_REQUEST to=2 lamport=380
2026/11/08 21:50:53.374818 [Node 3] event=SEND_REQUEST to=1 lamport=380
2026/11/08 21:50:53.375731 [Node 3] event=RECEIVE_REPLY from=2 lamport=382 local=383 replies=1/2
2026/11/08 21:50:53.622034 [Node 3] event=RECEIVE_REQUEST from=2 lamport=383 local=384 state=WANTED
2026/11/08 21:50:53.622034 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:50:53.645544 [Node 3] event=RECEIVE_REPLY from=1 lamport=386 local=387 replies=2/2
2026/11/08 21:50:53.645544 [Node 3] event=ENTER_CS lamport=387
2026/11/08 21:50:54.646334 [Node 3] event=EXIT_CS lamport=387 deferred=1
2026/11/08 21:50:54.646403 [Node 3] event=SEND_REPLY to=2 lamport=388
2026/11/08 21:50:55.950152 [Node 3] event=REQUEST_CS lamport=389
2026/11/08 21:50:55.950152 [Node 3] event=SEND_REQUEST to=1 lamport=389
2026/11/08 21:50:55.950152 [Node 3] event=SEND_REQUEST to=2 lamport=389
2026/11/08 21:50:55.950730 [Node 3] event=RECEIVE_REPLY from=2 lamport=391 local=392 replies=1/2
2026/11/08 21:50:55.950730 [Node 3] event=RECEIVE_REPLY from=1 lamport=391 local=393 replies=2/2
2026/11/08 21:50:55.950730 [Node 3] event=ENTER_CS lamport=393
2026/11/08 21:50:56.951804 [Node 3] event=EXIT_CS lamport=393 deferred=0
2026/11/08 21:50:57.272979 [Node 3] event=RECEIVE_REQUEST from=2 lamport=392 local=394 state=RELEASED
2026/11/08 21:50:57.272979 [Node 3] event=SEND_REPLY to=2 lamport=395
2026/11/08 21:50:57.501653 [Node 3] event=RECEIVE_REQUEST from=1 lamport=395 local=396 state=RELEASED
2026/11/08 21:50:57.501653 [Node 3] event=SEND_REPLY to=1 lamport=397
2026/11/08 21:50:58.818627 [Node 3] event=REQUEST_CS lamport=398
2026/11/08 21:50:58.818627 [Node 3] event=SEND_REQUEST to=1 lamport=398
2026/11/08 21:50:58.818627 [Node 3] event=SEND_REQUEST to=2 lamport=398
2026/11/08 21:50:58.819631 [Node 3] event=RECEIVE_REPLY from=2 lamport=400 local=401 replies=1/2
2026/11/08 21:50:59.276286 [Node 3] event=RECEIVE_REPLY from=1 lamport=401 local=402 replies=2/2
2026/11/08 21:50:59.276286 [Node 3] event=ENTER_CS lamport=402
2026/11/08 21:50:59.567455 [Node 3] event=RECEIVE_REQUEST from=2 lamport=401 local=403 state=HELD
2026/11/08 21:50:59.568023 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:51:00.276933 [Node 3] event=EXIT_CS lamport=403 deferred=1
2026/11/08 21:51:00.276933 [Node 3] event=SEND_REPLY to=2 lamport=404
2026/11/08 21:51:01.673731 [Node 3] event=REQUEST_CS lamport=405
2026/11/08 21:51:01.673731 [Node 3] event=SEND_REQUEST to=1 lamport=405
2026/11/08 21:51:01.673731 [Node 3] event=SEND_REQUEST to=2 lamport=405
2026/11/08 21:51:01.674312 [Node 3] event=RECEIVE_REPLY from=2 lamport=407 local=408 replies=1/2
2026/11/08 21:51:01.674858 [Node 3] event=RECEIVE_REPLY from=1 lamport=407 local=409 replies=2/2
2026/11/08 21:51:01.674858 [Node 3] event=ENTER_CS lamport=409
2026/11/08 21:51:02.600009 [Node 3] event=RECEIVE_REQUEST from=1 lamport=408 local=410 state=HELD
2026/11/08 21:51:02.600009 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:51:02.675822 [Node 3] event=EXIT_CS lamport=410 deferred=1
2026/11/08 21:51:02.675822 [Node 3] event=SEND_REPLY to=1 lamport=411
2026/11/08 21:51:04.803489 [Node 3] event=REQUEST_CS lamport=412
2026/11/08 21:51:04.803489 [Node 3] event=SEND_REQUEST to=1 lamport=412
2026/11/08 21:51:04.803489 [Node 3] event=SEND_REQUEST to=2 lamport=412
2026/11/08 21:51:04.803489 [Node 3] event=RECEIVE_REPLY from=2 lamport=414 local=415 replies=1/2
2026/11/08 21:51:04.804406 [Node 3] event=RECEIVE_REPLY from=1 lamport=414 local=416 replies=2/2
2026/11/08 21:51:04.804406 [Node 3] event=ENTER_CS lamport=416
2026/11/08 21:51:05.067411 [Node 3] event=RECEIVE_REQUEST from=2 lamport=415 local=417 state=HELD
2026/11/08 21:51:05.067411 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:51:05.805104 [Node 3] event=EXIT_CS lamport=417 deferred=1
2026/11/08 21:51:05.805104 [Node 3] event=SEND_REPLY to=2 lamport=418
2026/11/08 21:51:07.492757 [Node 3] event=RECEIVE_REQUEST from=1 lamport=418 local=419 state=RELEASED
2026/11/08 21:51:07.492757 [Node 3] event=SEND_REPLY to=1 lamport=420
2026/11/08 21:51:08.831113 [Node 3] event=REQUEST_CS lamport=421
2026/11/08 21:51:08.831113 [Node 3] event=SEND_REQUEST to=1 lamport=421
2026/11/08 21:51:08.831113 [Node 3] event=SEND_REQUEST to=2 lamport=421
2026/11/08 21:51:08.831672 [Node 3] event=RECEIVE_REPLY from=2 lamport=423 local=424 replies=1/2
2026/11/08 21:51:08.831672 [Node 3] event=RECEIVE_REPLY from=1 lamport=425 local=426 replies=2/2
2026/11/08 21:51:08.832205 [Node 3] event=ENTER_CS lamport=426
2026/11/08 21:51:09.832838 [Node 3] event=EXIT_CS lamport=426 deferred=0
2026/11/08 21:51:09.862988 [Node 3] event=RECEIVE_REQUEST from=2 lamport=424 local=427 state=RELEASED
2026/11/08 21:51:09.863016 [Node 3] event=SEND_REPLY to=2 lamport=428
2026/11/08 21:51:12.243857 [Node 3] event=RECEIVE_REQUEST from=1 lamport=428 local=429 state=RELEASED
2026/11/08 21:51:12.243877 [Node 3] event=SEND_REPLY to=1 lamport=430
2026/11/08 21:51:12.797544 [Node 3] event=REQUEST_CS lamport=431
2026/11/08 21:51:12.797544 [Node 3] event=SEND_REQUEST to=1 lamport=431
2026/11/08 21:51:12.797544 [Node 3] event=SEND_REQUEST to=2 lamport=431
2026/11/08 21:51:12.798135 [Node 3] event=RECEIVE_REPLY from=2 lamport=434 local=435 replies=1/2
2026/11/08 21:51:13.106345 [Node 3] event=RECEIVE_REQUEST from=2 lamport=435 local=436 state=WANTED
2026/11/08 21:51:13.106913 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:51:13.244530 [Node 3] event=RECEIVE_REPLY from=1 lamport=438 local=439 replies=2/2
2026/11/08 21:51:13.244530 [Node 3] event=ENTER_CS lamport=439
2026/11/08 21:51:14.244921 [Node 3] event=EXIT_CS lamport=439 deferred=1
2026/11/08 21:51:14.244921 [Node 3] event=SEND_REPLY to=2 lamport=440
2026/11/08 21:51:15.760122 [Node 3] event=RECEIVE_REQUEST from=1 lamport=439 local=441 state=RELEASED
2026/11/08 21:51:15.760727 [Node 3] event=SEND_REPLY to=1 lamport=442
2026/11/08 21:51:17.267352 [Node 3] event=RECEIVE_REQUEST from=2 lamport=444 local=445 state=RELEASED
2026/11/08 21:51:17.267352 [Node 3] event=SEND_REPLY to=2 lamport=446
2026/11/08 21:51:17.930584 [Node 3] event=REQUEST_CS lamport=447
2026/11/08 21:51:17.930584 [Node 3] event=SEND_REQUEST to=1 lamport=447
2026/11/08 21:51:17.930584 [Node 3] event=SEND_REQUEST to=2 lamport=447
2026/11/08 21:51:17.931175 [Node 3] event=RECEIVE_REPLY from=1 lamport=449 local=450 replies=1/2
2026/11/08 21:51:18.268516 [Node 3] event=RECEIVE_REPLY from=2 lamport=450 local=451 replies=2/2
2026/11/08 21:51:18.268516 [Node 3] event=ENTER_CS lamport=451
2026/11/08 21:51:18.685583 [Node 3] event=RECEIVE_REQUEST from=1 lamport=450 local=452 state=HELD
2026/11/08 21:51:18.685583 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:51:19.269385 [Node 3] event=EXIT_CS lamport=452 deferred=1
2026/11/08 21:51:19.269432 [Node 3] event=SEND_REPLY to=1 lamport=453
2026/11/08 21:51:19.718281 [Node 3] event=RECEIVE_REQUEST from=2 lamport=453 local=454 state=RELEASED
2026/11/08 21:51:19.718281 [Node 3] event=SEND_REPLY to=2 lamport=455
2026/11/08 21:51:20.709172 [Node 3] event=REQUEST_CS lamport=456
2026/11/08 21:51:20.709220 [Node 3] event=SEND_REQUEST to=1 lamport=456
2026/11/08 21:51:20.709220 [Node 3] event=SEND_REQUEST to=2 lamport=456
2026/11/08 21:51:20.709748 [Node 3] event=RECEIVE_REPLY from=1 lamport=458 local=459 replies=1/2
2026/11/08 21:51:21.271816 [Node 3] event=RECEIVE_REPLY from=2 lamport=459 local=460 replies=2/2
2026/11/08 21:51:21.271816 [Node 3] event=ENTER_CS lamport=460
2026/11/08 21:51:22.272338 [Node 3] event=EXIT_CS lamport=460 deferred=0
2026/11/08 21:51:22.844979 [Node 3] event=RECEIVE_REQUEST from=2 lamport=460 local=461 state=RELEASED
2026/11/08 21:51:22.844979 [Node 3] event=SEND_REPLY to=2 lamport=462
2026/11/08 21:51:23.731975 [Node 3] event=RECEIVE_REQUEST from=1 lamport=463 local=464 state=RELEASED
2026/11/08 21:51:23.732533 [Node 3] event=SEND_REPLY to=1 lamport=465
2026/11/08 21:51:25.015732 [Node 3] event=REQUEST_CS lamport=466
2026/11/08 21:51:25.015785 [Node 3] event=SEND_REQUEST to=1 lamport=466
2026/11/08 21:51:25.015785 [Node 3] event=SEND_REQUEST to=2 lamport=466
2026/11/08 21:51:25.016300 [Node 3] event=RECEIVE_REPLY from=2 lamport=468 local=469 replies=1/2
2026/11/08 21:51:25.016300 [Node 3] event=RECEIVE_REPLY from=1 lamport=469 local=470 replies=2/2
2026/11/08 21:51:25.016832 [Node 3] event=ENTER_CS lamport=470
2026/11/08 21:51:26.017783 [Node 3] event=EXIT_CS lamport=470 deferred=0
2026/11/08 21:51:27.571080 [Node 3] event=RECEIVE_REQUEST from=2 lamport=469 local=471 state=RELEASED
2026/11/08 21:51:27.571080 [Node 3] event=SEND_REPLY to=2 lamport=472
2026/11/08 21:51:28.133799 [Node 3] event=RECEIVE_REQUEST from=1 lamport=472 local=473 state=RELEASED
2026/11/08 21:51:28.133799 [Node 3] event=SEND_REPLY to=1 lamport=474
2026/11/08 21:51:28.233981 [Node 3] event=REQUEST_CS lamport=475
2026/11/08 21:51:28.233981 [Node 3] event=SEND_REQUEST to=1 lamport=475
2026/11/08 21:51:28.233981 [Node 3] event=SEND_REQUEST to=2 lamport=475
2026/11/08 21:51:28.572385 [Node 3] event=RECEIVE_REPLY from=2 lamport=477 local=478 replies=1/2
2026/11/08 21:51:29.573913 [Node 3] event=RECEIVE_REPLY from=1 lamport=480 local=481 replies=2/2
2026/11/08 21:51:29.573913 [Node 3] event=ENTER_CS lamport=481
2026/11/08 21:51:30.574392 [Node 3] event=EXIT_CS lamport=481 deferred=0
2026/11/08 21:51:32.454762 [Node 3] event=RECEIVE_REQUEST from=2 lamport=479 local=482 state=RELEASED
2026/11/08 21:51:32.454762 [Node 3] event=SEND_REPLY to=2 lamport=483
2026/11/08 21:51:32.485501 [Node 3] event=REQUEST_CS lamport=484
2026/11/08 21:51:32.485501 [Node 3] event=SEND_REQUEST to=2 lamport=484
2026/11/08 21:51:32.485501 [Node 3] event=SEND_REQUEST to=1 lamport=484
2026/11/08 21:51:32.486072 [Node 3] event=RECEIVE_REPLY from=1 lamport=486 local=487 replies=1/2
2026/11/08 21:51:32.567111 [Node 3] event=RECEIVE_REQUEST from=1 lamport=487 local=488 state=WANTED
2026/11/08 21:51:32.567111 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:51:33.456325 [Node 3] event=RECEIVE_REPLY from=2 lamport=490 local=491 replies=2/2
2026/11/08 21:51:33.456325 [Node 3] event=ENTER_CS lamport=491
2026/11/08 21:51:34.457351 [Node 3] event=EXIT_CS lamport=491 deferred=1
2026/11/08 21:51:34.457351 [Node 3] event=SEND_REPLY to=1 lamport=492
2026/11/08 21:51:34.644542 [Node 3] event=RECEIVE_REQUEST from=2 lamport=491 local=493 state=RELEASED
2026/11/08 21:51:34.644542 [Node 3] event=SEND_REPLY to=2 lamport=494
2026/11/08 21:51:36.411282 [Node 3] event=REQUEST_CS lamport=495
2026/11/08 21:51:36.411282 [Node 3] event=SEND_REQUEST to=1 lamport=495
2026/11/08 21:51:36.411282 [Node 3] event=SEND_REQUEST to=2 lamport=495
2026/11/08 21:51:36.411855 [Node 3] event=RECEIVE_REPLY from=1 lamport=497 local=498 replies=1/2
2026/11/08 21:51:36.459160 [Node 3] event=RECEIVE_REPLY from=2 lamport=498 local=499 replies=2/2
2026/11/08 21:51:36.459690 [Node 3] event=ENTER_CS lamport=499
2026/11/08 21:51:36.864310 [Node 3] event=RECEIVE_REQUEST from=1 lamport=498 local=500 state=HELD
2026/11/08 21:51:36.864310 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:51:37.460180 [Node 3] event=EXIT_CS lamport=500 deferred=1
2026/11/08 21:51:37.460228 [Node 3] event=SEND_REPLY to=1 lamport=501
2026/11/08 21:51:39.849876 [Node 3] event=RECEIVE_REQUEST from=2 lamport=501 local=502 state=RELEASED
2026/11/08 21:51:39.849876 [Node 3] event=SEND_REPLY to=2 lamport=503
2026/11/08 21:51:40.152163 [Node 3] event=RECEIVE_REQUEST from=1 lamport=505 local=506 state=RELEASED
2026/11/08 21:51:40.152163 [Node 3] event=SEND_REPLY to=1 lamport=507
2026/11/08 21:51:41.034166 [Node 3] event=REQUEST_CS lamport=508
2026/11/08 21:51:41.034166 [Node 3] event=SEND_REQUEST to=1 lamport=508
2026/11/08 21:51:41.034166 [Node 3] event=SEND_REQUEST to=2 lamport=508
2026/11/08 21:51:41.034925 [Node 3] event=RECEIVE_REPLY from=2 lamport=510 local=511 replies=1/2
2026/11/08 21:51:41.853036 [Node 3] event=RECEIVE_REPLY from=1 lamport=511 local=512 replies=2/2
2026/11/08 21:51:41.853619 [Node 3] event=ENTER_CS lamport=512
2026/11/08 21:51:42.854044 [Node 3] event=EXIT_CS lamport=512 deferred=0
2026/11/08 21:51:43.045604 [Node 3] event=RECEIVE_REQUEST from=1 lamport=512 local=513 state=RELEASED
2026/11/08 21:51:43.045604 [Node 3] event=SEND_REPLY to=1 lamport=514
2026/11/08 21:51:43.783829 [Node 3] event=RECEIVE_REQUEST from=2 lamport=515 local=516 state=RELEASED
2026/11/08 21:51:43.783829 [Node 3] event=SEND_REPLY to=2 lamport=517
2026/11/08 21:51:46.320235 [Node 3] event=RECEIVE_REQUEST from=1 lamport=519 local=520 state=RELEASED
2026/11/08 21:51:46.320235 [Node 3] event=SEND_REPLY to=1 lamport=521
2026/11/08 21:51:46.596278 [Node 3] event=REQUEST_CS lamport=522
2026/11/08 21:51:46.596278 [Node 3] event=SEND_REQUEST to=1 lamport=522
2026/11/08 21:51:46.596278 [Node 3] event=SEND_REQUEST to=2 lamport=522
2026/11/08 21:51:46.596865 [Node 3] event=RECEIVE_REPLY from=2 lamport=524 local=525 replies=1/2
2026/11/08 21:51:47.320465 [Node 3] event=RECEIVE_REPLY from=1 lamport=525 local=526 replies=2/2
2026/11/08 21:51:47.320465 [Node 3] event=ENTER_CS lamport=526
2026/11/08 21:51:47.594630 [Node 3] event=RECEIVE_REQUEST from=2 lamport=525 local=527 state=HELD
2026/11/08 21:51:47.594630 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:51:48.321298 [Node 3] event=EXIT_CS lamport=527 deferred=1
2026/11/08 21:51:48.321298 [Node 3] event=SEND_REPLY to=2 lamport=528
2026/11/08 21:51:50.668727 [Node 3] event=RECEIVE_REQUEST from=1 lamport=528 local=529 state=RELEASED
2026/11/08 21:51:50.668727 [Node 3] event=SEND_REPLY to=1 lamport=530
2026/11/08 21:51:51.743573 [Node 3] event=REQUEST_CS lamport=531
2026/11/08 21:51:51.743636 [Node 3] event=SEND_REQUEST to=1 lamport=531
2026/11/08 21:51:51.743636 [Node 3] event=SEND_REQUEST to=2 lamport=531
2026/11/08 21:51:51.744163 [Node 3] event=RECEIVE_REPLY from=1 lamport=535 local=536 replies=1/2
2026/11/08 21:51:51.744163 [Node 3] event=RECEIVE_REPLY from=2 lamport=533 local=537 replies=2/2
2026/11/08 21:51:51.744163 [Node 3] event=ENTER_CS lamport=537
2026/11/08 21:51:52.745039 [Node 3] event=EXIT_CS lamport=537 deferred=0
2026/11/08 21:51:52.798728 [Node 3] event=RECEIVE_REQUEST from=2 lamport=534 local=538 state=RELEASED
2026/11/08 21:51:52.799321 [Node 3] event=SEND_REPLY to=2 lamport=539
2026/11/08 21:51:53.321526 [Node 3] event=RECEIVE_REQUEST from=1 lamport=538 local=540 state=RELEASED
2026/11/08 21:51:53.321526 [Node 3] event=SEND_REPLY to=1 lamport=541
2026/11/08 21:51:55.014655 [Node 3] event=REQUEST_CS lamport=542
2026/11/08 21:51:55.014655 [Node 3] event=SEND_REQUEST to=2 lamport=542
2026/11/08 21:51:55.014655 [Node 3] event=SEND_REQUEST to=1 lamport=542
2026/11/08 21:51:55.015374 [Node 3] event=RECEIVE_REPLY from=1 lamport=545 local=546 replies=1/2
2026/11/08 21:51:55.015900 [Node 3] event=RECEIVE_REPLY from=2 lamport=544 local=547 replies=2/2
2026/11/08 21:51:55.015900 [Node 3] event=ENTER_CS lamport=547
2026/11/08 21:51:55.663282 [Node 3] event=RECEIVE_REQUEST from=2 lamport=545 local=548 state=HELD
2026/11/08 21:51:55.663282 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:51:56.016085 [Node 3] event=EXIT_CS lamport=548 deferred=1
2026/11/08 21:51:56.016085 [Node 3] event=SEND_REPLY to=2 lamport=549
2026/11/08 21:51:57.997063 [Node 3] event=RECEIVE_REQUEST from=1 lamport=548 local=550 state=RELEASED
2026/11/08 21:51:57.997063 [Node 3] event=SEND_REPLY to=1 lamport=551
2026/11/08 21:51:59.362568 [Node 3] event=RECEIVE_REQUEST from=2 lamport=553 local=554 state=RELEASED
2026/11/08 21:51:59.362568 [Node 3] event=SEND_REPLY to=2 lamport=555
2026/11/08 21:51:59.920763 [Node 3] event=REQUEST_CS lamport=556
2026/11/08 21:51:59.920763 [Node 3] event=SEND_REQUEST to=1 lamport=556
2026/11/08 21:51:59.920763 [Node 3] event=SEND_REQUEST to=2 lamport=556
2026/11/08 21:51:59.921927 [Node 3] event=RECEIVE_REPLY from=1 lamport=558 local=559 replies=1/2
2026/11/08 21:52:00.364440 [Node 3] event=RECEIVE_REPLY from=2 lamport=559 local=560 replies=2/2
2026/11/08 21:52:00.364440 [Node 3] event=ENTER_CS lamport=560
2026/11/08 21:52:00.527440 [Node 3] event=RECEIVE_REQUEST from=1 lamport=559 local=561 state=HELD
2026/11/08 21:52:00.527440 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:52:01.365439 [Node 3] event=EXIT_CS lamport=561 deferred=1
2026/11/08 21:52:01.365439 [Node 3] event=SEND_REPLY to=1 lamport=562
2026/11/08 21:52:03.382988 [Node 3] event=REQUEST_CS lamport=563
2026/11/08 21:52:03.382988 [Node 3] event=SEND_REQUEST to=1 lamport=563
2026/11/08 21:52:03.382988 [Node 3] event=SEND_REQUEST to=2 lamport=563
2026/11/08 21:52:03.383989 [Node 3] event=RECEIVE_REPLY from=2 lamport=565 local=566 replies=1/2
2026/11/08 21:52:03.383989 [Node 3] event=RECEIVE_REPLY from=1 lamport=565 local=567 replies=2/2
2026/11/08 21:52:03.383989 [Node 3] event=ENTER_CS lamport=567
2026/11/08 21:52:03.769988 [Node 3] event=RECEIVE_REQUEST from=2 lamport=566 local=568 state=HELD
2026/11/08 21:52:03.769988 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:52:04.384989 [Node 3] event=EXIT_CS lamport=568 deferred=1
2026/11/08 21:52:04.384989 [Node 3] event=SEND_REPLY to=2 lamport=569
2026/11/08 21:52:05.629229 [Node 3] event=RECEIVE_REQUEST from=1 lamport=569 local=570 state=RELEASED
2026/11/08 21:52:05.629229 [Node 3] event=SEND_REPLY to=1 lamport=571
2026/11/08 21:52:07.949735 [Node 3] event=REQUEST_CS lamport=572
2026/11/08 21:52:07.949735 [Node 3] event=SEND_REQUEST to=1 lamport=572
2026/11/08 21:52:07.949735 [Node 3] event=SEND_REQUEST to=2 lamport=572
2026/11/08 21:52:07.950738 [Node 3] event=RECEIVE_REPLY from=1 lamport=576 local=577 replies=1/2
2026/11/08 21:52:07.950738 [Node 3] event=RECEIVE_REPLY from=2 lamport=574 local=578 replies=2/2
2026/11/08 21:52:07.950738 [Node 3] event=ENTER_CS lamport=578
2026/11/08 21:52:08.476735 [Node 3] event=RECEIVE_REQUEST from=2 lamport=575 local=579 state=HELD
2026/11/08 21:52:08.476735 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:52:08.530735 [Node 3] event=RECEIVE_REQUEST from=1 lamport=579 local=580 state=HELD
2026/11/08 21:52:08.530735 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:52:08.951735 [Node 3] event=EXIT_CS lamport=580 deferred=2
2026/11/08 21:52:08.951735 [Node 3] event=SEND_REPLY to=1 lamport=581
2026/11/08 21:52:08.951735 [Node 3] event=SEND_REPLY to=2 lamport=582
2026/11/08 21:52:11.223268 [Node 3] event=RECEIVE_REQUEST from=2 lamport=585 local=586 state=RELEASED
2026/11/08 21:52:11.223268 [Node 3] event=SEND_REPLY to=2 lamport=587
2026/11/08 21:52:12.342578 [Node 3] event=RECEIVE_REQUEST from=1 lamport=588 local=589 state=RELEASED
2026/11/08 21:52:12.342578 [Node 3] event=SEND_REPLY to=1 lamport=590
2026/11/08 21:52:12.927316 [Node 3] event=REQUEST_CS lamport=591
2026/11/08 21:52:12.927316 [Node 3] event=SEND_REQUEST to=1 lamport=591
2026/11/08 21:52:12.927316 [Node 3] event=SEND_REQUEST to=2 lamport=591
2026/11/08 21:52:12.928192 [Node 3] event=RECEIVE_REPLY from=2 lamport=593 local=594 replies=1/2
2026/11/08 21:52:13.344254 [Node 3] event=RECEIVE_REPLY from=1 lamport=594 local=595 replies=2/2
2026/11/08 21:52:13.344254 [Node 3] event=ENTER_CS lamport=595
2026/11/08 21:52:14.344616 [Node 3] event=EXIT_CS lamport=595 deferred=0
2026/11/08 21:52:15.474377 [Node 3] event=RECEIVE_REQUEST from=1 lamport=595 local=596 state=RELEASED
2026/11/08 21:52:15.474377 [Node 3] event=SEND_REPLY to=1 lamport=597
2026/11/08 21:52:15.783378 [Node 3] event=RECEIVE_REQUEST from=2 lamport=598 local=599 state=RELEASED
2026/11/08 21:52:15.783378 [Node 3] event=SEND_REPLY to=2 lamport=600
2026/11/08 21:52:16.241378 [Node 3] event=REQUEST_CS lamport=601
2026/11/08 21:52:16.241378 [Node 3] event=SEND_REQUEST to=1 lamport=601
2026/11/08 21:52:16.241378 [Node 3] event=SEND_REQUEST to=2 lamport=601
2026/11/08 21:52:16.476377 [Node 3] event=RECEIVE_REPLY from=1 lamport=603 local=604 replies=1/2
2026/11/08 21:52:17.477377 [Node 3] event=RECEIVE_REPLY from=2 lamport=606 local=607 replies=2/2
2026/11/08 21:52:17.477377 [Node 3] event=ENTER_CS lamport=607
2026/11/08 21:52:18.477911 [Node 3] event=EXIT_CS lamport=607 deferred=0
2026/11/08 21:52:18.596911 [Node 3] event=RECEIVE_REQUEST from=1 lamport=605 local=608 state=RELEASED
2026/11/08 21:52:18.596911 [Node 3] event=SEND_REPLY to=1 lamport=609
2026/11/08 21:52:19.843921 [Node 3] event=RECEIVE_REQUEST from=2 lamport=609 local=610 state=RELEASED
2026/11/08 21:52:19.843921 [Node 3] event=SEND_REPLY to=2 lamport=611
2026/11/08 21:52:21.510580 [Node 3] event=REQUEST_CS lamport=612
2026/11/08 21:52:21.510580 [Node 3] event=SEND_REQUEST to=1 lamport=612
2026/11/08 21:52:21.510580 [Node 3] event=SEND_REQUEST to=2 lamport=612
2026/11/08 21:52:21.511394 [Node 3] event=RECEIVE_REPLY from=1 lamport=615 local=616 replies=1/2
2026/11/08 21:52:21.511394 [Node 3] event=RECEIVE_REPLY from=2 lamport=617 local=618 replies=2/2
2026/11/08 21:52:21.511394 [Node 3] event=ENTER_CS lamport=618
2026/11/08 21:52:22.226676 [Node 3] event=RECEIVE_REQUEST from=1 lamport=616 local=619 state=HELD
2026/11/08 21:52:22.226676 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:52:22.511464 [Node 3] event=EXIT_CS lamport=619 deferred=1
2026/11/08 21:52:22.511464 [Node 3] event=SEND_REPLY to=1 lamport=620
2026/11/08 21:52:23.303831 [Node 3] event=RECEIVE_REQUEST from=2 lamport=620 local=621 state=RELEASED
2026/11/08 21:52:23.303831 [Node 3] event=SEND_REPLY to=2 lamport=622
2026/11/08 21:52:25.432832 [Node 3] event=RECEIVE_REQUEST from=1 lamport=624 local=625 state=RELEASED
2026/11/08 21:52:25.432832 [Node 3] event=SEND_REPLY to=1 lamport=626
2026/11/08 21:52:26.044832 [Node 3] event=REQUEST_CS lamport=627
2026/11/08 21:52:26.044832 [Node 3] event=SEND_REQUEST to=1 lamport=627
2026/11/08 21:52:26.044832 [Node 3] event=SEND_REQUEST to=2 lamport=627
2026/11/08 21:52:26.044832 [Node 3] event=RECEIVE_REPLY from=2 lamport=629 local=630 replies=1/2
2026/11/08 21:52:26.434831 [Node 3] event=RECEIVE_REPLY from=1 lamport=630 local=631 replies=2/2
2026/11/08 21:52:26.434831 [Node 3] event=ENTER_CS lamport=631
2026/11/08 21:52:27.435388 [Node 3] event=EXIT_CS lamport=631 deferred=0
2026/11/08 21:52:28.171637 [Node 3] event=RECEIVE_REQUEST from=2 lamport=630 local=632 state=RELEASED
2026/11/08 21:52:28.171637 [Node 3] event=SEND_REPLY to=2 lamport=633
2026/11/08 21:52:30.430493 [Node 3] event=RECEIVE_REQUEST from=1 lamport=633 local=634 state=RELEASED
2026/11/08 21:52:30.430493 [Node 3] event=SEND_REPLY to=1 lamport=635
2026/11/08 21:52:30.678492 [Node 3] event=REQUEST_CS lamport=636
2026/11/08 21:52:30.678492 [Node 3] event=SEND_REQUEST to=1 lamport=636
2026/11/08 21:52:30.678492 [Node 3] event=SEND_REQUEST to=2 lamport=636
2026/11/08 21:52:30.679496 [Node 3] event=RECEIVE_REPLY from=2 lamport=638 local=639 replies=1/2
2026/11/08 21:52:31.432133 [Node 3] event=RECEIVE_REPLY from=1 lamport=639 local=640 replies=2/2
2026/11/08 21:52:31.432133 [Node 3] event=ENTER_CS lamport=640
2026/11/08 21:52:32.433695 [Node 3] event=EXIT_CS lamport=640 deferred=0
2026/11/08 21:52:32.514293 [Node 3] event=RECEIVE_REQUEST from=2 lamport=639 local=641 state=RELEASED
2026/11/08 21:52:32.514293 [Node 3] event=SEND_REPLY to=2 lamport=642
2026/11/08 21:52:33.019374 [Node 3] event=RECEIVE_REQUEST from=1 lamport=642 local=643 state=RELEASED
2026/11/08 21:52:33.019374 [Node 3] event=SEND_REPLY to=1 lamport=644
2026/11/08 21:52:33.892506 [Node 3] event=REQUEST_CS lamport=645
2026/11/08 21:52:33.892506 [Node 3] event=SEND_REQUEST to=1 lamport=645
2026/11/08 21:52:33.892506 [Node 3] event=SEND_REQUEST to=2 lamport=645
2026/11/08 21:52:33.893622 [Node 3] event=RECEIVE_REPLY from=2 lamport=648 local=649 replies=1/2
2026/11/08 21:52:34.517144 [Node 3] event=RECEIVE_REPLY from=1 lamport=649 local=650 replies=2/2
2026/11/08 21:52:34.517144 [Node 3] event=ENTER_CS lamport=650
2026/11/08 21:52:35.518291 [Node 3] event=EXIT_CS lamport=650 deferred=0
2026/11/08 21:52:35.552890 [Node 3] event=RECEIVE_REQUEST from=2 lamport=649 local=651 state=RELEASED
2026/11/08 21:52:35.552890 [Node 3] event=SEND_REPLY to=2 lamport=652
2026/11/08 21:52:36.323273 [Node 3] event=RECEIVE_REQUEST from=1 lamport=652 local=653 state=RELEASED
2026/11/08 21:52:36.323292 [Node 3] event=SEND_REPLY to=1 lamport=654
2026/11/08 21:52:38.828896 [Node 3] event=REQUEST_CS lamport=655
2026/11/08 21:52:38.828896 [Node 3] event=SEND_REQUEST to=1 lamport=655
2026/11/08 21:52:38.828896 [Node 3] event=SEND_REQUEST to=2 lamport=655
2026/11/08 21:52:38.829692 [Node 3] event=RECEIVE_REPLY from=2 lamport=658 local=659 replies=1/2
2026/11/08 21:52:38.829692 [Node 3] event=RECEIVE_REPLY from=1 lamport=659 local=660 replies=2/2
2026/11/08 21:52:38.829692 [Node 3] event=ENTER_CS lamport=660
2026/11/08 21:52:39.296317 [Node 3] event=RECEIVE_REQUEST from=1 lamport=660 local=661 state=HELD
2026/11/08 21:52:39.296913 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:52:39.830238 [Node 3] event=EXIT_CS lamport=661 deferred=1
2026/11/08 21:52:39.830238 [Node 3] event=SEND_REPLY to=1 lamport=662
2026/11/08 21:52:40.133315 [Node 3] event=RECEIVE_REQUEST from=2 lamport=663 local=664 state=RELEASED
2026/11/08 21:52:40.133315 [Node 3] event=SEND_REPLY to=2 lamport=665
2026/11/08 21:52:43.336352 [Node 3] event=RECEIVE_REQUEST from=2 lamport=668 local=669 state=RELEASED
2026/11/08 21:52:43.336859 [Node 3] event=SEND_REPLY to=2 lamport=670
2026/11/08 21:52:43.813437 [Node 3] event=REQUEST_CS lamport=671
2026/11/08 21:52:43.813437 [Node 3] event=SEND_REQUEST to=1 lamport=671
2026/11/08 21:52:43.813437 [Node 3] event=SEND_REQUEST to=2 lamport=671
2026/11/08 21:52:43.814065 [Node 3] event=RECEIVE_REPLY from=1 lamport=673 local=674 replies=1/2
2026/11/08 21:52:43.825579 [Node 3] event=RECEIVE_REQUEST from=1 lamport=674 local=675 state=WANTED
2026/11/08 21:52:43.825579 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:52:44.338659 [Node 3] event=RECEIVE_REPLY from=2 lamport=677 local=678 replies=2/2
2026/11/08 21:52:44.338659 [Node 3] event=ENTER_CS lamport=678
2026/11/08 21:52:45.339623 [Node 3] event=EXIT_CS lamport=678 deferred=1
2026/11/08 21:52:45.339623 [Node 3] event=SEND_REPLY to=1 lamport=679
2026/11/08 21:52:47.784360 [Node 3] event=RECEIVE_REQUEST from=2 lamport=678 local=680 state=RELEASED
2026/11/08 21:52:47.784360 [Node 3] event=SEND_REPLY to=2 lamport=681
2026/11/08 21:52:48.388782 [Node 3] event=RECEIVE_REQUEST from=1 lamport=683 local=684 state=RELEASED
2026/11/08 21:52:48.389424 [Node 3] event=SEND_REPLY to=1 lamport=685
2026/11/08 21:52:48.402490 [Node 3] event=REQUEST_CS lamport=686
2026/11/08 21:52:48.402490 [Node 3] event=SEND_REQUEST to=1 lamport=686
2026/11/08 21:52:48.402490 [Node 3] event=SEND_REQUEST to=2 lamport=686
2026/11/08 21:52:48.785579 [Node 3] event=RECEIVE_REPLY from=2 lamport=688 local=689 replies=1/2
2026/11/08 21:52:49.785985 [Node 3] event=RECEIVE_REPLY from=1 lamport=691 local=692 replies=2/2
2026/11/08 21:52:49.785985 [Node 3] event=ENTER_CS lamport=692
2026/11/08 21:52:50.786176 [Node 3] event=EXIT_CS lamport=692 deferred=0
2026/11/08 21:52:51.133723 [Node 3] event=RECEIVE_REQUEST from=2 lamport=690 local=693 state=RELEASED
2026/11/08 21:52:51.133723 [Node 3] event=SEND_REPLY to=2 lamport=694
2026/11/08 21:52:52.198320 [Node 3] event=REQUEST_CS lamport=695
2026/11/08 21:52:52.198385 [Node 3] event=SEND_REQUEST to=1 lamport=695
2026/11/08 21:52:52.198385 [Node 3] event=SEND_REQUEST to=2 lamport=695
2026/11/08 21:52:52.198917 [Node 3] event=RECEIVE_REPLY from=2 lamport=697 local=698 replies=1/2
2026/11/08 21:52:52.198917 [Node 3] event=RECEIVE_REPLY from=1 lamport=697 local=699 replies=2/2
2026/11/08 21:52:52.198917 [Node 3] event=ENTER_CS lamport=699
2026/11/08 21:52:52.881008 [Node 3] event=RECEIVE_REQUEST from=1 lamport=698 local=700 state=HELD
2026/11/08 21:52:52.881008 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:52:53.199962 [Node 3] event=EXIT_CS lamport=700 deferred=1
2026/11/08 21:52:53.199962 [Node 3] event=SEND_REPLY to=1 lamport=701
2026/11/08 21:52:55.466536 [Node 3] event=RECEIVE_REQUEST from=2 lamport=701 local=702 state=RELEASED
2026/11/08 21:52:55.466536 [Node 3] event=SEND_REPLY to=2 lamport=703
2026/11/08 21:52:56.993302 [Node 3] event=REQUEST_CS lamport=704
2026/11/08 21:52:56.993302 [Node 3] event=SEND_REQUEST to=1 lamport=704
2026/11/08 21:52:56.993302 [Node 3] event=SEND_REQUEST to=2 lamport=704
2026/11/08 21:52:56.993909 [Node 3] event=RECEIVE_REPLY from=2 lamport=707 local=708 replies=1/2
2026/11/08 21:52:56.993909 [Node 3] event=RECEIVE_REPLY from=1 lamport=706 local=709 replies=2/2
2026/11/08 21:52:56.993909 [Node 3] event=ENTER_CS lamport=709
2026/11/08 21:52:57.829182 [Node 3] event=RECEIVE_REQUEST from=1 lamport=707 local=710 state=HELD
2026/11/08 21:52:57.829182 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:52:57.994282 [Node 3] event=EXIT_CS lamport=710 deferred=1
2026/11/08 21:52:57.994282 [Node 3] event=SEND_REPLY to=1 lamport=711
2026/11/08 21:53:00.099526 [Node 3] event=RECEIVE_REQUEST from=1 lamport=713 local=714 state=RELEASED
2026/11/08 21:53:00.099526 [Node 3] event=SEND_REPLY to=1 lamport=715
2026/11/08 21:53:00.170788 [Node 3] event=REQUEST_CS lamport=716
2026/11/08 21:53:00.170818 [Node 3] event=SEND_REQUEST to=1 lamport=716
2026/11/08 21:53:00.170818 [Node 3] event=SEND_REQUEST to=2 lamport=716
2026/11/08 21:53:00.171375 [Node 3] event=RECEIVE_REPLY from=2 lamport=718 local=719 replies=1/2
2026/11/08 21:53:00.358929 [Node 3] event=RECEIVE_REQUEST from=2 lamport=719 local=720 state=WANTED
2026/11/08 21:53:00.358929 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:53:01.100212 [Node 3] event=RECEIVE_REPLY from=1 lamport=722 local=723 replies=2/2
2026/11/08 21:53:01.100212 [Node 3] event=ENTER_CS lamport=723
2026/11/08 21:53:02.101340 [Node 3] event=EXIT_CS lamport=723 deferred=1
2026/11/08 21:53:02.101392 [Node 3] event=SEND_REPLY to=2 lamport=724
2026/11/08 21:53:03.404129 [Node 3] event=RECEIVE_REQUEST from=1 lamport=723 local=725 state=RELEASED
2026/11/08 21:53:03.404129 [Node 3] event=SEND_REPLY to=1 lamport=726
2026/11/08 21:53:05.356477 [Node 3] event=REQUEST_CS lamport=727
2026/11/08 21:53:05.356513 [Node 3] event=SEND_REQUEST to=1 lamport=727
2026/11/08 21:53:05.356513 [Node 3] event=SEND_REQUEST to=2 lamport=727
2026/11/08 21:53:05.357054 [Node 3] event=RECEIVE_REPLY from=2 lamport=729 local=730 replies=1/2
2026/11/08 21:53:05.357054 [Node 3] event=RECEIVE_REPLY from=1 lamport=730 local=731 replies=2/2
2026/11/08 21:53:05.357054 [Node 3] event=ENTER_CS lamport=731
2026/11/08 21:53:05.943536 [Node 3] event=RECEIVE_REQUEST from=2 lamport=730 local=732 state=HELD
2026/11/08 21:53:05.944099 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:53:06.357310 [Node 3] event=EXIT_CS lamport=732 deferred=1
2026/11/08 21:53:06.357310 [Node 3] event=SEND_REPLY to=2 lamport=733
2026/11/08 21:53:07.934340 [Node 3] event=RECEIVE_REQUEST from=1 lamport=733 local=734 state=RELEASED
2026/11/08 21:53:07.934340 [Node 3] event=SEND_REPLY to=1 lamport=735
2026/11/08 21:53:09.378293 [Node 3] event=REQUEST_CS lamport=736
2026/11/08 21:53:09.378339 [Node 3] event=SEND_REQUEST to=1 lamport=736
2026/11/08 21:53:09.378339 [Node 3] event=SEND_REQUEST to=2 lamport=736
2026/11/08 21:53:09.378924 [Node 3] event=RECEIVE_REPLY from=1 lamport=739 local=740 replies=1/2
2026/11/08 21:53:09.378924 [Node 3] event=RECEIVE_REPLY from=2 lamport=738 local=741 replies=2/2
2026/11/08 21:53:09.378924 [Node 3] event=ENTER_CS lamport=741
2026/11/08 21:53:09.676010 [Node 3] event=RECEIVE_REQUEST from=2 lamport=739 local=742 state=HELD
2026/11/08 21:53:09.676010 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:53:10.379664 [Node 3] event=EXIT_CS lamport=742 deferred=1
2026/11/08 21:53:10.379711 [Node 3] event=SEND_REPLY to=2 lamport=743
2026/11/08 21:53:11.569590 [Node 3] event=REQUEST_CS lamport=744
2026/11/08 21:53:11.569638 [Node 3] event=SEND_REQUEST to=1 lamport=744
2026/11/08 21:53:11.569638 [Node 3] event=SEND_REQUEST to=2 lamport=744
2026/11/08 21:53:11.570234 [Node 3] event=RECEIVE_REPLY from=2 lamport=746 local=747 replies=1/2
2026/11/08 21:53:11.570234 [Node 3] event=RECEIVE_REPLY from=1 lamport=746 local=748 replies=2/2
2026/11/08 21:53:11.570234 [Node 3] event=ENTER_CS lamport=748
2026/11/08 21:53:12.308074 [Node 3] event=RECEIVE_REQUEST from=1 lamport=747 local=749 state=HELD
2026/11/08 21:53:12.308074 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:53:12.571287 [Node 3] event=EXIT_CS lamport=749 deferred=1
2026/11/08 21:53:12.571332 [Node 3] event=SEND_REPLY to=1 lamport=750
2026/11/08 21:53:14.241356 [Node 3] event=REQUEST_CS lamport=751
2026/11/08 21:53:14.241356 [Node 3] event=SEND_REQUEST to=1 lamport=751
2026/11/08 21:53:14.241356 [Node 3] event=SEND_REQUEST to=2 lamport=751
2026/11/08 21:53:14.241946 [Node 3] event=RECEIVE_REPLY from=1 lamport=753 local=754 replies=1/2
2026/11/08 21:53:14.241946 [Node 3] event=RECEIVE_REPLY from=2 lamport=753 local=755 replies=2/2
2026/11/08 21:53:14.241946 [Node 3] event=ENTER_CS lamport=755
2026/11/08 21:53:14.648115 [Node 3] event=RECEIVE_REQUEST from=2 lamport=754 local=756 state=HELD
2026/11/08 21:53:14.648115 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:53:15.242428 [Node 3] event=EXIT_CS lamport=756 deferred=1
2026/11/08 21:53:15.242475 [Node 3] event=SEND_REPLY to=2 lamport=757
2026/11/08 21:53:16.295397 [Node 3] event=RECEIVE_REQUEST from=1 lamport=757 local=758 state=RELEASED
2026/11/08 21:53:16.295397 [Node 3] event=SEND_REPLY to=1 lamport=759
2026/11/08 21:53:17.622912 [Node 3] event=RECEIVE_REQUEST from=2 lamport=761 local=762 state=RELEASED
2026/11/08 21:53:17.622912 [Node 3] event=SEND_REPLY to=2 lamport=763
2026/11/08 21:53:18.840093 [Node 3] event=RECEIVE_REQUEST from=1 lamport=764 local=765 state=RELEASED
2026/11/08 21:53:18.840723 [Node 3] event=SEND_REPLY to=1 lamport=766
2026/11/08 21:53:18.872771 [Node 3] event=REQUEST_CS lamport=767
2026/11/08 21:53:18.872813 [Node 3] event=SEND_REQUEST to=1 lamport=767
2026/11/08 21:53:18.872813 [Node 3] event=SEND_REQUEST to=2 lamport=767
2026/11/08 21:53:18.873340 [Node 3] event=RECEIVE_REPLY from=2 lamport=769 local=770 replies=1/2
2026/11/08 21:53:19.841708 [Node 3] event=RECEIVE_REPLY from=1 lamport=771 local=772 replies=2/2
2026/11/08 21:53:19.841708 [Node 3] event=ENTER_CS lamport=772
2026/11/08 21:53:20.842218 [Node 3] event=EXIT_CS lamport=772 deferred=0
2026/11/08 21:53:21.981175 [Node 3] event=RECEIVE_REQUEST from=2 lamport=770 local=773 state=RELEASED
2026/11/08 21:53:21.981175 [Node 3] event=SEND_REPLY to=2 lamport=774
2026/11/08 21:53:23.042352 [Node 3] event=REQUEST_CS lamport=775
2026/11/08 21:53:23.042352 [Node 3] event=SEND_REQUEST to=1 lamport=775
2026/11/08 21:53:23.042352 [Node 3] event=SEND_REQUEST to=2 lamport=775
2026/11/08 21:53:23.042927 [Node 3] event=RECEIVE_REPLY from=2 lamport=778 local=779 replies=1/2
2026/11/08 21:53:23.042927 [Node 3] event=RECEIVE_REPLY from=1 lamport=777 local=780 replies=2/2
2026/11/08 21:53:23.042927 [Node 3] event=ENTER_CS lamport=780
2026/11/08 21:53:23.249318 [Node 3] event=RECEIVE_REQUEST from=1 lamport=778 local=781 state=HELD
2026/11/08 21:53:23.249318 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:53:24.043619 [Node 3] event=EXIT_CS lamport=781 deferred=1
2026/11/08 21:53:24.043619 [Node 3] event=SEND_REPLY to=1 lamport=782
2026/11/08 21:53:26.216897 [Node 3] event=RECEIVE_REQUEST from=1 lamport=784 local=785 state=RELEASED
2026/11/08 21:53:26.216897 [Node 3] event=SEND_REPLY to=1 lamport=786
2026/11/08 21:53:26.397761 [Node 3] event=RECEIVE_REQUEST from=2 lamport=787 local=788 state=RELEASED
2026/11/08 21:53:26.397761 [Node 3] event=SEND_REPLY to=2 lamport=789
2026/11/08 21:53:26.778923 [Node 3] event=REQUEST_CS lamport=790
2026/11/08 21:53:26.778923 [Node 3] event=SEND_REQUEST to=1 lamport=790
2026/11/08 21:53:26.778923 [Node 3] event=SEND_REQUEST to=2 lamport=790
2026/11/08 21:53:27.218626 [Node 3] event=RECEIVE_REPLY from=1 lamport=792 local=793 replies=1/2
2026/11/08 21:53:28.219146 [Node 3] event=RECEIVE_REPLY from=2 lamport=795 local=796 replies=2/2
2026/11/08 21:53:28.219684 [Node 3] event=ENTER_CS lamport=796
2026/11/08 21:53:29.220719 [Node 3] event=EXIT_CS lamport=796 deferred=0
2026/11/08 21:53:30.136332 [Node 3] event=RECEIVE_REQUEST from=1 lamport=794 local=797 state=RELEASED
2026/11/08 21:53:30.136839 [Node 3] event=SEND_REPLY to=1 lamport=798
2026/11/08 21:53:31.461728 [Node 3] event=REQUEST_CS lamport=799
2026/11/08 21:53:31.461770 [Node 3] event=SEND_REQUEST to=2 lamport=799
2026/11/08 21:53:31.461770 [Node 3] event=SEND_REQUEST to=1 lamport=799
2026/11/08 21:53:31.462316 [Node 3] event=RECEIVE_REPLY from=2 lamport=801 local=802 replies=1/2
2026/11/08 21:53:31.462316 [Node 3] event=RECEIVE_REPLY from=1 lamport=802 local=803 replies=2/2
2026/11/08 21:53:31.462316 [Node 3] event=ENTER_CS lamport=803
2026/11/08 21:53:31.714403 [Node 3] event=RECEIVE_REQUEST from=2 lamport=802 local=804 state=HELD
2026/11/08 21:53:31.714403 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:53:32.462733 [Node 3] event=EXIT_CS lamport=804 deferred=1
2026/11/08 21:53:32.462733 [Node 3] event=SEND_REPLY to=2 lamport=805
2026/11/08 21:53:33.627070 [Node 3] event=REQUEST_CS lamport=806
2026/11/08 21:53:33.627070 [Node 3] event=SEND_REQUEST to=1 lamport=806
2026/11/08 21:53:33.627070 [Node 3] event=SEND_REQUEST to=2 lamport=806
2026/11/08 21:53:33.627668 [Node 3] event=RECEIVE_REPLY from=2 lamport=808 local=809 replies=1/2
2026/11/08 21:53:33.627668 [Node 3] event=RECEIVE_REPLY from=1 lamport=808 local=810 replies=2/2
2026/11/08 21:53:33.627668 [Node 3] event=ENTER_CS lamport=810
2026/11/08 21:53:34.468272 [Node 3] event=RECEIVE_REQUEST from=2 lamport=809 local=811 state=HELD
2026/11/08 21:53:34.468272 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:53:34.627891 [Node 3] event=EXIT_CS lamport=811 deferred=1
2026/11/08 21:53:34.627891 [Node 3] event=SEND_REPLY to=2 lamport=812
2026/11/08 21:53:34.639947 [Node 3] event=RECEIVE_REQUEST from=1 lamport=812 local=813 state=RELEASED
2026/11/08 21:53:34.639947 [Node 3] event=SEND_REPLY to=1 lamport=814
2026/11/08 21:53:37.050969 [Node 3] event=REQUEST_CS lamport=815
2026/11/08 21:53:37.050969 [Node 3] event=SEND_REQUEST to=1 lamport=815
2026/11/08 21:53:37.050969 [Node 3] event=SEND_REQUEST to=2 lamport=815
2026/11/08 21:53:37.051566 [Node 3] event=RECEIVE_REPLY from=2 lamport=817 local=818 replies=1/2
2026/11/08 21:53:37.051566 [Node 3] event=RECEIVE_REPLY from=1 lamport=818 local=819 replies=2/2
2026/11/08 21:53:37.051566 [Node 3] event=ENTER_CS lamport=819
2026/11/08 21:53:38.052183 [Node 3] event=EXIT_CS lamport=819 deferred=0
2026/11/08 21:53:39.039847 [Node 3] event=RECEIVE_REQUEST from=2 lamport=818 local=820 state=RELEASED
2026/11/08 21:53:39.039847 [Node 3] event=SEND_REPLY to=2 lamport=821
2026/11/08 21:53:39.296639 [Node 3] event=RECEIVE_REQUEST from=1 lamport=821 local=822 state=RELEASED
2026/11/08 21:53:39.296639 [Node 3] event=SEND_REPLY to=1 lamport=823
2026/11/08 21:53:40.936457 [Node 3] event=REQUEST_CS lamport=824
2026/11/08 21:53:40.936519 [Node 3] event=SEND_REQUEST to=1 lamport=824
2026/11/08 21:53:40.936519 [Node 3] event=SEND_REQUEST to=2 lamport=824
2026/11/08 21:53:40.937048 [Node 3] event=RECEIVE_REPLY from=2 lamport=827 local=828 replies=1/2
2026/11/08 21:53:41.042232 [Node 3] event=RECEIVE_REPLY from=1 lamport=828 local=829 replies=2/2
2026/11/08 21:53:41.042232 [Node 3] event=ENTER_CS lamport=829
2026/11/08 21:53:42.043249 [Node 3] event=EXIT_CS lamport=829 deferred=0
2026/11/08 21:53:42.693477 [Node 3] event=RECEIVE_REQUEST from=1 lamport=829 local=830 state=RELEASED
2026/11/08 21:53:42.693500 [Node 3] event=SEND_REPLY to=1 lamport=831
2026/11/08 21:53:43.366391 [Node 3] event=RECEIVE_REQUEST from=2 lamport=832 local=833 state=RELEASED
2026/11/08 21:53:43.366391 [Node 3] event=SEND_REPLY to=2 lamport=834
2026/11/08 21:53:43.532445 [Node 3] event=REQUEST_CS lamport=835
2026/11/08 21:53:43.532502 [Node 3] event=SEND_REQUEST to=1 lamport=835
2026/11/08 21:53:43.532502 [Node 3] event=SEND_REQUEST to=2 lamport=835
2026/11/08 21:53:43.693733 [Node 3] event=RECEIVE_REPLY from=1 lamport=837 local=838 replies=1/2
2026/11/08 21:53:44.694563 [Node 3] event=RECEIVE_REPLY from=2 lamport=840 local=841 replies=2/2
2026/11/08 21:53:44.694563 [Node 3] event=ENTER_CS lamport=841
2026/11/08 21:53:45.695544 [Node 3] event=EXIT_CS lamport=841 deferred=0
2026/11/08 21:53:46.442171 [Node 3] event=RECEIVE_REQUEST from=2 lamport=841 local=842 state=RELEASED
2026/11/08 21:53:46.442171 [Node 3] event=SEND_REPLY to=2 lamport=843
2026/11/08 21:53:47.115312 [Node 3] event=RECEIVE_REQUEST from=1 lamport=844 local=845 state=RELEASED
2026/11/08 21:53:47.115312 [Node 3] event=SEND_REPLY to=1 lamport=846
2026/11/08 21:53:47.555583 [Node 3] event=REQUEST_CS lamport=847
2026/11/08 21:53:47.555583 [Node 3] event=SEND_REQUEST to=1 lamport=847
2026/11/08 21:53:47.555583 [Node 3] event=SEND_REQUEST to=2 lamport=847
2026/11/08 21:53:47.556179 [Node 3] event=RECEIVE_REPLY from=2 lamport=849 local=850 replies=1/2
2026/11/08 21:53:48.444398 [Node 3] event=RECEIVE_REPLY from=1 lamport=850 local=851 replies=2/2
2026/11/08 21:53:48.444398 [Node 3] event=ENTER_CS lamport=851
2026/11/08 21:53:48.917249 [Node 3] event=RECEIVE_REQUEST from=2 lamport=850 local=852 state=HELD
2026/11/08 21:53:48.917804 [Node 3] event=DEFER_REPLY to=2
2026/11/08 21:53:49.445554 [Node 3] event=EXIT_CS lamport=852 deferred=1
2026/11/08 21:53:49.445605 [Node 3] event=SEND_REPLY to=2 lamport=853
2026/11/08 21:53:51.800334 [Node 3] event=REQUEST_CS lamport=854
2026/11/08 21:53:51.800334 [Node 3] event=SEND_REQUEST to=1 lamport=854
2026/11/08 21:53:51.800334 [Node 3] event=SEND_REQUEST to=2 lamport=854
2026/11/08 21:53:51.800334 [Node 3] event=RECEIVE_REPLY from=1 lamport=856 local=857 replies=1/2
2026/11/08 21:53:51.800334 [Node 3] event=RECEIVE_REPLY from=2 lamport=856 local=858 replies=2/2
2026/11/08 21:53:51.800334 [Node 3] event=ENTER_CS lamport=858
2026/11/08 21:53:52.246978 [Node 3] event=RECEIVE_REQUEST from=1 lamport=857 local=859 state=HELD
2026/11/08 21:53:52.246978 [Node 3] event=DEFER_REPLY to=1
2026/11/08 21:53:52.801544 [Node 3] event=EXIT_CS lamport=859 deferred=1
2026/11/08 21:53:52.801544 [Node 3] event=SEND_REPLY to=1 lamport=860
2026/11/08 21:53:52.802711 [Node 3] event=RECEIVE_REQUEST from=2 lamport=860 local=861 state=RELEASED
2026/11/08 21:53:52.802711 [Node 3] event=SEND_REPLY to=2 lamport=862
2026/11/08 21:53:55.579388 [Node 3] event=REQUEST_CS lamport=863
2026/11/08 21:53:55.579388 [Node 3] event=SEND_REQUEST to=2 lamport=863
2026/11/08 21:53:55.579388 [Node 3] event=SEND_REQUEST to=1 lamport=863
2026/11/08 21:53:55.580526 [Node 3] event=RECEIVE_REPLY from=1 lamport=865 local=866 replies=1/2
2026/11/08 21:53:55.580526 [Node 3] event=RECEIVE_REPLY from=2 lamport=866 local=867 replies=2/2
2026/11/08 21:53:55.580526 [Node 3] event=ENTER_CS lamport=867
2026/11/08 21:53:55.904672 [Node 3] event=RECEIVE_REQUEST from=2 lamport=867 local=868 state=HELD
2026/11/08 21:53:55.904672 [Node 3] event=DEFER_REPLY to=2
```