package lib

type IQueue[T any] interface {
	Push(item T)
	Pop() T
	IsEmpty() bool
	Size() int
	All() []T
}

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Pop() T {
	if q.IsEmpty() {
		panic("not item in queue")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

func (q *Queue[T]) All() []T {
	return q.items
}

func NewQueue[T any](data []T) IQueue[T] {
	q := &Queue[T]{}
	for i := 0; i < len(data); i++ {
		q.Push(data[i])
	}

	return q
}
