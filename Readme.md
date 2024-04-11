# Go TCP HTTP Proxy Server

This project demonstrated proxy-ing tcp connection into a http server and return the response back to tcp

## Prerequisites

- Go 1.16 or later

## Running Server

```bash
go run server/server.go
```

will run the tcp and http server

## Running Client

```bash
go run client/client.go
```

will run the tcp client. it includes a cli prompt to enter the amount to simulate transaction