package main

type Task struct {
	priority int
	Floor    int
	ToFloor  int
	IsPickup bool
}

type TaskPQ []Task

func (self TaskPQ) Len() int { return len(self) }
func (self TaskPQ) Less(i, j int) bool {
	return self[i].priority < self[j].priority
}
func (self TaskPQ) Swap(i, j int)       { self[i], self[j] = self[j], self[i] }
func (self *TaskPQ) Push(x interface{}) { *self = append(*self, x.(Task)) }
func (self *TaskPQ) Pop() (popped interface{}) {
	popped = (*self)[len(*self)-1]
	*self = (*self)[:len(*self)-1]
	return
}
