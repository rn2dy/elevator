package main

import "fmt"

type Task struct {
	priority int
	Floor    int // pickup floor if IsPickup, otherwise means dropoff floor
	ToFloor  int // only used if IsPickup
	IsPickup bool
}

func (t Task) String() string {
	if t.IsPickup {
		return fmt.Sprintf("[From: %d, To: %d]", t.Floor, t.ToFloor)	
	}
	return fmt.Sprintf("[To: %d]", t.Floor)
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

func (self *TaskPQ) Peek() interface{} {
	return (*self)[len(*self)-1]
}