package lib

import (
	"errors"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
)

type IWorkerQueue interface {
	Match(task *dao.Tasks) error
	OrderQueue(task *[]dao.Tasks) error
	Size() int
}

var (
	ErrAllWorkerAreCrossed = errors.New("all workers are crossed")
	ErrNotHaveWorker       = errors.New("not have workers available")
)

type WorkerQueue struct {
	Crossed []IWorker
	Main    IQueue[IWorker]
}

func (workerQueue *WorkerQueue) Match(task *dao.Tasks) error {

	if task.Status == constant.TaskStatus_Canceled {
		task.MemberId = nil
		task.Person = nil
		return nil
	}

	if len(workerQueue.Crossed) > 0 {
		for i, worker := range workerQueue.Crossed {
			if err := worker.AddTask(task); err == nil {
				workerQueue.Crossed = append(workerQueue.Crossed[:i], workerQueue.Crossed[i+1:]...) // delete worker crossed
				workerQueue.Main.Push(worker)
				return nil
			}
		}
	}

	if workerQueue.Main.IsEmpty() && len(workerQueue.Crossed) == 0 {
		return ErrNotHaveWorker
	} else if workerQueue.Main.IsEmpty() {
		return ErrAllWorkerAreCrossed
	}

	for i := 0; i < workerQueue.Main.Size(); i++ {
		worker := workerQueue.Main.Pop()
		if err := worker.AddTask(task); err == nil {
			workerQueue.Main.Push(worker)
			return nil
		} else {
			workerQueue.Crossed = append(workerQueue.Crossed, worker) // skip queue
		}
	}

	return ErrNotHaveWorker

}

// tasks must softed date Desc
func (workerQueue *WorkerQueue) OrderQueue(tasks *[]dao.Tasks) error {
	items := workerQueue.Main.All()
	for _, task := range *tasks {
		for i, item := range items {
			if item.GetId() == *task.MemberId {
				items = append(items[:i], items[i+1:]...)
				items = append(items, item)
			}
		}
	}

	workerQueue.Main = NewQueue(items)
	return nil
}

func (WorkerQueue *WorkerQueue) Size() int {
	return WorkerQueue.Main.Size() + len(WorkerQueue.Crossed)
}

func NewWorkerQueue(workers []IWorker) IWorkerQueue {
	return &WorkerQueue{
		Main: NewQueue(workers),
	}
}
