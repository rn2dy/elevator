package main

import (
	"container/heap"
	"fmt"
	"time"
)

// EStatus elevator status
type EStatus int

const (
	Idle EStatus = iota
	MovingUp
	MovingDown
)

// Elevator struct
type Elevator struct {
	ID           int
	CurrentFloor int
	NextFloor    int
	TasksQueue   TaskPQ // priority queue
	Status       EStatus
	Controller   chan Event // used to communicate with elevator controller
}

// NewElevator contructor
func NewElevator(id int) *Elevator {
	el := &Elevator{
		ID:     id,
		Status: Idle,
	}
	heap.Init(&el.TasksQueue)
	return el
}


// Pickup schedule a pickup
func (el *Elevator) Pickup(fromFloor, toFloor int) {
	priority := Abs(el.CurrentFloor - fromFloor)
	heap.Push(&el.TasksQueue, Task{priority, fromFloor, toFloor, true})
	if el.TasksQueue.Len() == 1 {
		el.NextFloor = el.TasksQueue[0].Floor
	}
	el.SetStatus()
}

// NextStop move the elevator to the next stop floor
func (el *Elevator) NextStop() {
	if el.TasksQueue.Len() > 0 {
		task := heap.Pop(&el.TasksQueue).(Task)
		el.CurrentFloor = task.Floor
		if task.IsPickup {
			el.registerStop(task.ToFloor)
		}
		if el.TasksQueue.Len() > 0 {
			el.NextFloor = el.TasksQueue[0].Floor
		} else {
			el.NextFloor = el.CurrentFloor
		}
		el.SetStatus()
	}
}

func (el *Elevator) registerStop(floor int) {
	// skip floor if already registered
	for _, task := range el.TasksQueue {
		if task.Floor == floor {
			return
		}
	}
	// register a stop 
	priority := Abs(el.CurrentFloor - floor)
	heap.Push(&el.TasksQueue, Task{
		priority: priority, 
		Floor: floor, 
		IsPickup: false,
	})
}

// Stops calculates the stops between floor and current floor
func (el *Elevator) Stops(floor int) int {
	var stops int
	for _, task := range el.TasksQueue {
		if el.CurrentFloor > floor { // down
			if el.CurrentFloor > task.Floor && task.Floor > floor {
				stops++
			}
		} else { // up
			if el.CurrentFloor < task.Floor && task.Floor < floor {
				stops++
			}
		}
	}
	return stops
}

// SetStatus set elevator's status
func (el *Elevator) SetStatus() {
	oldStatus := el.Status
	if el.TasksQueue.Len() == 0 {
		el.Status = Idle
		// nothing to do, go to the top floor
		el.CurrentFloor = NumOfFloors - 1
		el.NextFloor = el.CurrentFloor
	} else {
		if el.NextFloor > el.CurrentFloor {
			el.Status = MovingUp
		} else {
			el.Status = MovingDown
		}
	}
	if oldStatus != el.Status {
		go func(){
 			el.Controller <- Event{Name: "STATUS_CHANGE"}
		}()
	}
}

// StatusString return a formatted string to show current status
func (el *Elevator) StatusString() string {
	var status = ""
	switch el.Status {
	case MovingUp:
		status = "Moving up..."
	case MovingDown:
		status = "Moving down..."
	case Idle:
		status = "Idling..."
	}
	return fmt.Sprintf("[Elevator_%d %s] Current floor: %d, Next floor: %d, Queue: %v",
		el.ID,
		status,
		el.CurrentFloor,
		el.NextFloor,
		el.TasksQueue)
}

// Run simulate elevator move up/down with a time ticking loop
func (el *Elevator) Run() {
	// send status update to controller every 3s
	go func() {
		for {
			<-time.Tick(StatusUpdateFrequency)
			el.Controller <- Event{
				Name: "STATUS",
				Data: el.StatusString(),
			}
		}
	}()

	// every 300ms the elevator moves to the next target floor
	for {
		select {
		case <-time.Tick(ElevatorSpeed):
			el.NextStop()
		}
	}
}
