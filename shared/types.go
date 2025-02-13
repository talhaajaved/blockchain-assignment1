// shared/types.go
package shared

// Matrix is a 2D slice of integers.
type Matrix [][]int

// Operation represents the type of matrix operation.
type Operation string

const (
	Add       Operation = "add"
	Multiply  Operation = "multiply"
	Transpose Operation = "transpose"
)

// MatrixOperationRequest is the RPC request type.
type MatrixOperationRequest struct {
	Operation Operation
	MatrixA   Matrix
	MatrixB   Matrix // Used for addition and multiplication.
}

// MatrixOperationResponse is the RPC response type.
type MatrixOperationResponse struct {
	Result Matrix
	Error  string
}
