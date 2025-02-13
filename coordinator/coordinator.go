// coordinator/coordinator.go
package coordinator

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/talhaajaved/blockchain-assignment1/shared"
)

// Coordinator receives client requests and delegates them to workers.
type Coordinator struct {
	WorkerManager *WorkerManager
	UseTLS        bool
	TLSConfig     *tls.Config
}

// NewCoordinator creates a new Coordinator with the given worker addresses.
// When useTLS is true, tlsConfig is used for dialing workers.
func NewCoordinator(workerAddresses []string, useTLS bool, tlsConfig *tls.Config) *Coordinator {
	wm := NewWorkerManager(workerAddresses)
	return &Coordinator{
		WorkerManager: wm,
		UseTLS:        useTLS,
		TLSConfig:     tlsConfig,
	}
}

// HandleRequest is the RPC method called by clients.
func (c *Coordinator) HandleRequest(req *shared.MatrixOperationRequest, resp *shared.MatrixOperationResponse) error {
	log.Printf("Coordinator: Received request: %+v", req)
	maxRetries := len(c.WorkerManager.Workers)
	triedWorker := make(map[string]bool)
	attempts := 0

	for attempts < maxRetries {
		// Get a worker not already tried for this request.
		worker := c.WorkerManager.getLeastBusyWorker(triedWorker)
		if worker == nil {
			log.Println("Coordinator: No available worker after excluding tried ones.")
			break
		}

		// Mark this worker as tried.
		triedWorker[worker.Address] = true

		log.Printf("Coordinator: Selected worker at %s", worker.Address)
		worker.increment()

		var conn net.Conn
		var dialErr error
		if c.UseTLS {
			// Use tls.DialWithDialer to set a timeout when dialing a TLS worker.
			conn, dialErr = tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", worker.Address, c.TLSConfig)
		} else {
			conn, dialErr = net.DialTimeout("tcp", worker.Address, 3*time.Second)
		}
		if dialErr != nil {
			log.Printf("Coordinator: Error dialing worker at %s: %v", worker.Address, dialErr)
			worker.decrement()
			attempts++
			continue
		}

		client := rpc.NewClient(conn)
		var workerResp shared.MatrixOperationResponse
		log.Printf("Coordinator: Calling WorkerService.Compute on worker %s", worker.Address)
		callErr := client.Call("WorkerService.Compute", req, &workerResp)
		client.Close()
		if callErr != nil {
			log.Printf("Coordinator: Error during RPC call to worker at %s: %v", worker.Address, callErr)
			worker.decrement()
			attempts++
			continue
		}

		log.Printf("Coordinator: Received response from worker at %s: %+v", worker.Address, workerResp)
		*resp = workerResp
		worker.decrement()
		return nil
	}

	log.Println("Coordinator: All workers failed to process the request.")
	return errors.New("all workers failed to process the request")
}
