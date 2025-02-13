// cmd/client/main.go
package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/talhaajaved/blockchain-assignment1/shared"
)

func main() {
	var serverAddr string
	var useTLS bool
	flag.StringVar(&serverAddr, "server", "localhost:8000", "Coordinator server address")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS for secure connection")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Matrix Operation Client")
	fmt.Println("-----------------------")

	fmt.Print("Enter operation (add, multiply, transpose): ")
	opInput, _ := reader.ReadString('\n')
	opInput = strings.TrimSpace(opInput)
	operation := shared.Operation(opInput)

	fmt.Println("Enter Matrix A (enter rows; separate numbers by spaces, empty line to finish):")
	matrixA := readMatrix(reader)

	var matrixB shared.Matrix
	if operation == shared.Add || operation == shared.Multiply {
		fmt.Println("Enter Matrix B (same format as Matrix A):")
		matrixB = readMatrix(reader)
	}

	req := &shared.MatrixOperationRequest{
		Operation: operation,
		MatrixA:   matrixA,
		MatrixB:   matrixB,
	}
	log.Printf("Client: Sending request: %+v", req)

	var client *rpc.Client
	var err error
	if useTLS {
		config := &tls.Config{InsecureSkipVerify: true}
		conn, err := tls.Dial("tcp", serverAddr, config)
		if err != nil {
			log.Fatalf("Client: Failed to connect to coordinator: %v", err)
		}
		client = rpc.NewClient(conn)
	} else {
		client, err = rpc.Dial("tcp", serverAddr)
		if err != nil {
			log.Fatalf("Client: Failed to connect to coordinator: %v", err)
		}
	}
	defer client.Close()

	var resp shared.MatrixOperationResponse
	err = client.Call("Coordinator.HandleRequest", req, &resp)
	if err != nil {
		log.Fatalf("Client: RPC call failed: %v", err)
	}
	log.Printf("Client: Received response: %+v", resp)

	if resp.Error != "" {
		fmt.Println("Error:", resp.Error)
	} else {
		fmt.Println("Result:")
		printMatrix(resp.Result)
	}
}

func readMatrix(reader *bufio.Reader) shared.Matrix {
	var matrix shared.Matrix
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		fields := strings.Fields(line)
		row := make([]int, len(fields))
		for i, field := range fields {
			num, err := strconv.Atoi(field)
			if err != nil {
				fmt.Printf("Invalid number '%s', using 0.\n", field)
				num = 0
			}
			row[i] = num
		}
		matrix = append(matrix, row)
	}
	return matrix
}

func printMatrix(matrix shared.Matrix) {
	for _, row := range matrix {
		fmt.Println(row)
	}
}
