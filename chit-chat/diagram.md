# Lamport timestamp trace

Client Y is already in the chat (it sent Join with L=1, the server broadcast at 2 and Y received it at 3). Client X then joins, publishes "Hello", Client Y answers "Hi X", and Client X leaves.

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

- The numbers to the left of Client X and to the right of Client Y are that client's Lamport clock after the event on that line. The max(...) notes next to the server show how the server's clock is updated.
- L= is the timestamp sent with a message.
- A client adds 1 to its clock before it sends a call. When a process receives a message, it sets its clock to max(own clock, L) + 1. The server broadcasts with its new clock value.
- Empty replies carry no timestamp, so they do not change any clock.
- "X joined" is the broadcast "Participant X joined Chit Chat at logical time 3" and "X left" is "Participant X left Chit Chat at logical time 12".

## Step by step

1. X sends Join with L=1 (0 + 1).
2. The server sets its clock to max(2, 1) + 1 = 3 and broadcasts "Participant X joined Chit Chat at logical time 3" to X and Y.
3. X sets its clock to max(1, 3) + 1 = 4. Y sets its clock to max(3, 3) + 1 = 4.
4. X sends Publish "Hello" with L=5 (4 + 1).
5. The server sets its clock to max(3, 5) + 1 = 6 and broadcasts "Hello" with L=6 to X and Y.
6. X sets its clock to max(5, 6) + 1 = 7. Y sets its clock to max(4, 6) + 1 = 7.
7. Y sends Publish "Hi X" with L=8 (7 + 1).
8. The server sets its clock to max(6, 8) + 1 = 9 and broadcasts "Hi X" with L=9 to X and Y.
9. X sets its clock to max(7, 9) + 1 = 10. Y sets its clock to max(8, 9) + 1 = 10.
10. X sends Leave with L=11 (10 + 1).
11. The server closes X's Join stream, sets its clock to max(9, 11) + 1 = 12 and broadcasts "Participant X left Chit Chat at logical time 12" to Y.
12. Y sets its clock to max(10, 12) + 1 = 13.
