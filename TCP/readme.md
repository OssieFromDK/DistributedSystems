a) Packages and data structure
We have made a `Packet2` struct, which is the data structure we use, it has a `Type string`, which has the control flag identifying which pakcet is being sent like "SYN", "SYN-ACK", "ACK" etc. and a `Number uint32` which is the payload, and acts as the sequence number which the handshake increments

b) Threads vs processes
In our implementation we use goroutines which basically is like a thread
Its not realistic to use threads since threads are simply a local thing, on a real network the client and server don't live in the same process so they don't share any memory or channels like our threads here do, they are bound to some socket accepting connections etc., threads are just used here for simulating the protocol on our local machine

c) Message Reordering
This is the entire point of the `Number` field in our Packet struct. Its a sequence number which solves reordering.
On the reciever, we maintain what the next expected number in the sequence should be, and then if it arrives early, we can buffer it, and if its behind or a duplicate, we can drop it, but still acknowledge we got it
Since the current code is just a simple handshake we don't do this, but for a real network it would be how.

d) Message delay and loss
Our forwarder drops 20% of packages and delays the rest randomly by 0-1000ms, if a packet is dropped, the server/client has a timeout expecting and answer, if if it hasn't recieved it, it times out and retries, after 5 tries it gives up and closes the connection

e) Why the 3-way handshake matters
Its basically to make sure everything works, after step 1, the server knows the client can send, after step 2, the client knows the server can send and recieve, and after step 3 the server knows the client can recieve, both sides have then proved to each other that they can transmit and recieve data before any real data is being sent. The sequence number in the handshake is important for the things described above, like packet losses, and message reordering, both sides keeps track of the packet id's and knows if something is out of order, making sure all the data arrives properly in the right order, which is what TCP is for, reliably sending data