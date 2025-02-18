package lib

import (
	"sync"

	"github.com/google/uuid"
)

type WorkerQueueMap struct {
	mu    sync.RWMutex // Use RWMutex for better read performance
	Items map[uuid.UUID]IWorkerQueue
}

func (wq *WorkerQueueMap) Set(key uuid.UUID, item IWorkerQueue) {
	wq.mu.Lock()
	wq.Items[key] = item
	wq.mu.Unlock()
}

func (wq *WorkerQueueMap) Get(key uuid.UUID) IWorkerQueue {
	wq.mu.RLock() // Use read lock for concurrent reads
	defer wq.mu.RUnlock()
	item, exists := wq.Items[key]
	if !exists {
		panic("not fount key in WorkerQueue")
	}
	return item
}

func NewWorkerQueueMap() *WorkerQueueMap {
	return &WorkerQueueMap{
		Items: make(map[uuid.UUID]IWorkerQueue),
	}
}
