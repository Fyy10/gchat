# gchat

A simple C/S architecture chat app written in Go.

## Build

### Go build

Build server:

```bash
go build ./cmd/server
```

Build client:

```bash
go build ./cmd/client
```

### Make

Build server & client:

```bash
make
```

The executables will be in the `./bin/` folder

Clean up builds:

```bash
make clean
```

## Usage

Server:

```bash
./server -h # print help message
./server -ip 0.0.0.0 -port 8080 # run the server to listen on 0.0.0.0:8080
```

Client:

```bash
./client -h # print help message
./client -ip 127.0.0.1 -port 8080 # run the client to connect to 127.0.0.1:8080
```

### Client commands

- `help`: show the supported commands
- `who`: display all the online users
- `whoami`: display the current username
- `rename`: change the username
- `@user`: talk to an online user privately

Please type `help` in client for more details.
