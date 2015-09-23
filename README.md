## Elevator Control

### Build

Please install [Go programming language](https://golang.org/doc/install)

Build

```
make
```

Run

```
make run
```

### Implementations

There are four pieces of this elevator simulation program, controller.go, elevator.go, pqueue.go and main.go.

__controller.go__

This file implements the data structure of a _Controller_ as well as the interface methods associated with it.

It has an event loop which listens on _status_ and _pickup_ events and performs actions accordingly.

__elevator.go__

This file implements the data structure of an _Elevator_ as well as the interface methods associated with it.

When it starts (by calling `elevator.Run()`) it will send status update to its controller and simulate the moving action based on pre-configured time intervals.

__pqueue.go__

This file implements a priority queue data structure.

__main.go__

This file bootstrap the simulation program.

It instantiates the instances of elevators and controller and generate _pickup_ events and send to the controller based on pre-configured time interval.

Here is a list of options and its default value:

```
NumOfElevators        = 4               // number of elevators to play with
PickupFrequency       = 1 * time.Second // pickup request frequency
ElevatorSpeed         = 4 * time.Second // move to next floor every 5 seconds
StatusUpdateFrequency = 3 * time.Second // elevator status update frequency
NumOfFloors           = 24
```

### Scheduling

__FCFS__

Serves pickup request in order. For example, suppose there are a sequence of pickup requests like {from: 0, to: 10}, {from: 20, to: 0}, {from: 2, to: 8}, the elevator actually going to not pick up from 2nd floor, just because the pickup from 20th floor happened first, which is very inefficient! In worst case senario, the elevator can bounce between low floors to high floors without picking up any passengers along the way.

__FCFS improvements__

There are some heuristics can be applied to help with the race condition caused by FCFS.

1. Pick Idle elevators to serve pickup requests (and choose the one that is closest to the pickup floor).
2. Pick up passengers along the way (same direction). Also choose the elevator that has minimum stops from its current floor to the pickup floor.
3. If none of above applies, add the pickup to a waiting list. And if any of the elevator sends a status change event, process the waiting list.

__More Improvements__

3 options here:

1. Always move idle elevators to the 1st floor, because the 1st floor usually have a larger traffic, is it?
2. Always move idle elevators to the top floor since gravity can help save energy and along the way it can pickup passengers. (implemented)
3. Move majority of the idle elevators to the 1st floor and the rest to the top floor.


### Data Structures

A _Controller_ is defined as

```go
type Controller struct {
	Elevators []*Elevator
	PickupQueue []Pickup // a waiting list
	EventChan chan Event // receives pickup event and elevator's status update event
}
```

The EventChan is the main communication channel between controller thread and all the elevators' threads. It also recieves pickup events generated from yet another thread.

An _Elevator_ is defined as

```go
type Elevator struct {
	ID           int
	CurrentFloor int
	NextFloor    int
	TasksQueue   TaskPQ // priority queue
	Status       EStatus
	Controller   chan Event // used to communicate with elevator controller
}
```

Each instance of elevator owns a priority queue (called _TasksQueue_) that is used to manage the elevators' stops. The priority of each _Task_ (stop) is determined based on the number of stops between the elevator's current floor and the pickup floor and the priorities of the queue will gets updated whenever the elevator moved to a new floor.