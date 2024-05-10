package main

type Queue[T any] struct {
	arr []T
}

func (q *Queue[T]) Push(val T) {
	q.arr = append(q.arr, val)
}

func (q *Queue[T]) Pop() (T, bool) {
	var val T
	if len(q.arr) == 0 {
		return val, false
	}

	val = q.arr[0]
	q.arr = q.arr[1:]

	return val, true
}

func (q *Queue[T]) Length() int {
	return len(q.arr)
}

func (q *Queue[T]) Peek() (T, bool) {
	var val T
	if len(q.arr) == 0 {
		return val, false
	}

	val = q.arr[0]

	return val, true
}
