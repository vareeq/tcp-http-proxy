# Go TCP HTTP Proxy Server

This project demonstrates how to proxy a TCP connection through an HTTP server and return the response back to the TCP client.

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