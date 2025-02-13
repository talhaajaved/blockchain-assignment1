// cmd/coordinator/main.go
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/talhaajaved/blockchain-assignment1/coordinator"
)

func main() {
	var port string
	var useTLS bool
	flag.StringVar(&port, "port", "8000", "Port to listen on for client RPC")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS for secure connection")
	flag.Parse()

	var tlsConfig *tls.Config
	if useTLS {
		certFile := "certs/server.crt"
		keyFile := "certs/server.key"
		// Load the TLS certificate for both listening and dialing.
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("Coordinator: Failed to load key pair: %v", err)
		}
		// For self-signed certificates, you may want to disable verification.
		tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
	}

	// Define the worker addresses.
	workerAddresses := []string{
		"localhost:9000",
		"localhost:9001",
		"localhost:9002",
	}

	// Create the Coordinator instance, passing the TLS configuration.
	coord := coordinator.NewCoordinator(workerAddresses, useTLS, tlsConfig)
	if err := rpc.Register(coord); err != nil {
		log.Fatalf("Coordinator: Error registering RPC service: %v", err)
	}

	addr := ":" + port
	var listener net.Listener
	var err error

	if useTLS {
		listener, err = tls.Listen("tcp", addr, tlsConfig)
		if err != nil {
			log.Fatalf("Coordinator: Failed to listen with TLS: %v", err)
		}
		fmt.Println("Coordinator running with TLS on port", port)
	} else {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("Coordinator: Failed to listen on port %s: %v", port, err)
		}
		fmt.Println("Coordinator running on port", port)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Coordinator: Error accepting connection:", err)
			continue
		}
		log.Println("Coordinator: Accepted new connection")
		go rpc.ServeConn(conn)
	}
}
