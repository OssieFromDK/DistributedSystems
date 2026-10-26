# Chit Chat

## How to run

Run every command from the project root. Start the server first:

```
go run ./server
```

Then start each client in its own terminal with a unique name:

```
go run ./client alice
go run ./client bob
go run ./client charlie
```

Type a message and press Enter to publish it
Type `/quit` to leave the chat as a client