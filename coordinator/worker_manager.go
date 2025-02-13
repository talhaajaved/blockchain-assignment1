// coordinator/worker_manager.go
package coordinator

import "sync"

// Worker represents a worker process.
type Worker struct {
	Address    string
	ActiveJobs int
	mu         sync.Mutex
}

// WorkerManager maintains a list of workers.
type WorkerManager struct {
	Workers []*Worker
	mu      sync.Mutex
}

// NewWorkerManager creates a new WorkerManager from a list of addresses.
func NewWorkerManager(addresses []string) *WorkerManager {
	workers := make([]*Worker, 0, len(addresses))
	for _, addr := range addresses {
		workers = append(workers, &Worker{Address: addr})
	}
	return &WorkerManager{Workers: workers}
}

// getLeastBusyWorker returns the worker with the fewest active jobs that is NOT in the exclude map.
func (wm *WorkerManager) getLeastBusyWorker(exclude map[string]bool) *Worker {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	var selected *Worker
	for _, w := range wm.Workers {
		if exclude[w.Address] {
			continue
		}
		w.mu.Lock()
		if selected == nil || w.ActiveJobs < selected.ActiveJobs {
			if selected != nil {
				selected.mu.Unlock()
			}
			selected = w
		} else {
			w.mu.Unlock()
		}
	}
	return selected
}

func (w *Worker) increment() {
	w.mu.Lock()
	w.ActiveJobs++
	w.mu.Unlock()
}

func (w *Worker) decrement() {
	w.mu.Lock()
	if w.ActiveJobs > 0 {
		w.ActiveJobs--
	}
	w.mu.Unlock()
}
