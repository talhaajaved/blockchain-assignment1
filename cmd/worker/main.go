package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/talhaajaved/blockchain-assignment1/worker"
)

func main() {
	var port string
	var useTLS bool
	flag.StringVar(&port, "port", "9000", "Port to listen on")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS for secure connection")
	flag.Parse()

	// Register the worker RPC service.
	ws := new(worker.WorkerService)
	if err := rpc.Register(ws); err != nil {
		log.Fatalf("Worker: Error registering RPC service: %v", err)
	}

	addr := ":" + port
	var listener net.Listener
	var err error

	if useTLS {
		// Load the certificate and key files from the certs directory.
		certFile := "certs/server.crt"
		keyFile := "certs/server.key"
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("Worker: Failed to load key pair: %v", err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", addr, config)
		if err != nil {
			log.Fatalf("Worker: Failed to listen with TLS: %v", err)
		}
		fmt.Println("Worker running with TLS on port", port)
	} else {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("Worker: Failed to listen on port %s: %v", port, err)
		}
		fmt.Println("Worker running on port", port)
	}

	// Accept connections and serve RPC requests.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Worker: Error accepting connection:", err)
			continue
		}
		log.Println("Worker: Accepted new connection")
		go rpc.ServeConn(conn)
	}
}
