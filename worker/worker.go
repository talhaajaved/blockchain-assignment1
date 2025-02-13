// worker/worker.go
package worker

import (
	"errors"
	"log"

	"github.com/talhaajaved/blockchain-assignment1/shared"
)

// WorkerService provides matrix operation services.
type WorkerService struct{}

// Compute performs the requested matrix operation.
func (w *WorkerService) Compute(req *shared.MatrixOperationRequest, resp *shared.MatrixOperationResponse) error {
	log.Printf("Worker: Received compute request: %+v", req)
	switch req.Operation {
	case shared.Add:
		result, err := addMatrices(req.MatrixA, req.MatrixB)
		if err != nil {
			resp.Error = err.Error()
			log.Printf("Worker: Error in addition: %v", err)
			return nil
		}
		resp.Result = result
	case shared.Multiply:
		result, err := multiplyMatrices(req.MatrixA, req.MatrixB)
		if err != nil {
			resp.Error = err.Error()
			log.Printf("Worker: Error in multiplication: %v", err)
			return nil
		}
		resp.Result = result
	case shared.Transpose:
		resp.Result = transposeMatrix(req.MatrixA)
	default:
		err := errors.New("unsupported operation")
		resp.Error = err.Error()
		log.Printf("Worker: Unsupported operation: %v", req.Operation)
	}
	log.Println("Worker: Completed compute request.")
	return nil
}

func addMatrices(a, b shared.Matrix) (shared.Matrix, error) {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) || len(a[0]) != len(b[0]) {
		return nil, errors.New("matrices dimensions do not match for addition")
	}
	rows, cols := len(a), len(a[0])
	result := make(shared.Matrix, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = a[i][j] + b[i][j]
		}
	}
	return result, nil
}

func multiplyMatrices(a, b shared.Matrix) (shared.Matrix, error) {
	if len(a) == 0 || len(b) == 0 || len(a[0]) != len(b) {
		return nil, errors.New("matrices dimensions do not match for multiplication")
	}
	rowsA, colsA, colsB := len(a), len(a[0]), len(b[0])
	result := make(shared.Matrix, rowsA)
	for i := 0; i < rowsA; i++ {
		result[i] = make([]int, colsB)
		for j := 0; j < colsB; j++ {
			sum := 0
			for k := 0; k < colsA; k++ {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}
	return result, nil
}

func transposeMatrix(a shared.Matrix) shared.Matrix {
	if len(a) == 0 {
		return shared.Matrix{}
	}
	rows, cols := len(a), len(a[0])
	result := make(shared.Matrix, cols)
	for i := 0; i < cols; i++ {
		result[i] = make([]int, rows)
		for j := 0; j < rows; j++ {
			result[i][j] = a[j][i]
		}
	}
	return result
}
