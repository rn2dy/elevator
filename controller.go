package main

import (
	"fmt"
	"math"
	"time"
)

// Event controller recieved event
type Event struct {
	Name string
	Data interface{}
}

type Pickup struct {
	FromFloor, ToFloor int
}

// Controller evelator controller
type Controller struct {
	Elevators []*Elevator
	PickupQueue []Pickup // a waiting list
	EventChan chan Event
}

// NewController constructor
func NewController() *Controller {
	return &Controller{EventChan: make(chan Event)}
}

// RegisterElevators registers elevator (setup communication channel)
func (ctrl *Controller) RegisterElevators(elevators ...*Elevator) {
	for _, el := range elevators {
		el.Controller = ctrl.EventChan
	}
	ctrl.Elevators = elevators
}

// Pickup select one elevator to pick up passenger
func (ctrl *Controller) Pickup(fromFloor, toFloor int) bool {
	var stops = math.MaxInt64
	var elevator *Elevator
	
	// choose idle elevator if there is one
	idleElevators := ctrl.Status(Idle)
	if len(idleElevators) > 0 {
		var delta = math.MaxInt64
		for _, el := range idleElevators {
			d := Abs(el.CurrentFloor - fromFloor)
			if d < delta {
				elevator = el
				delta = d
			}
		}
		elevator.Pickup(fromFloor, toFloor)
		return true
	}
	
	// otherwise
	// if request moving up
	if fromFloor < toFloor {
		upElevators := ctrl.Status(MovingUp)
		if len(upElevators) > 0 {
			for _, el := range upElevators {
				if el.CurrentFloor < fromFloor {
					s := el.Stops(fromFloor)
					if s < stops {
						elevator = el
						stops = s
					}
				}
			}
		}
	}
	
	// if request moving down
	if fromFloor > toFloor {
		downElevators := ctrl.Status(MovingDown)
		if len(downElevators) > 0 {
			for _, el := range downElevators {
				if el.CurrentFloor > fromFloor {
					s := el.Stops(fromFloor)
					if s < stops {
						elevator = el
						stops = s
					}	
				}
			}
		}
	}
	
	// if no match found put request into wait queue
	if elevator == nil {
		pickup := Pickup{fromFloor, toFloor}
		ctrl.PickupQueue = append(ctrl.PickupQueue, pickup)
		fmt.Printf("Add pickup to waiting queue (%d): %v\n", pickup, len(ctrl.PickupQueue))
		return false
	}
	elevator.Pickup(fromFloor, toFloor)
	return true
}

// Status query for all elevator matching status s
func (ctrl *Controller) Status(s EStatus) []*Elevator {
	var elevators []*Elevator
	for _, el := range ctrl.Elevators {
		if el.Status == s {
			elevators = append(elevators, el)
		}
	}
	return elevators
}

// Start simulation a event listener loop
func (ctrl *Controller) Start() {
	go func() {
		for {
			<-time.Tick(StatusUpdateFrequency)
			for _, el := range ctrl.Elevators {
				fmt.Println(el.StatusString())	
			}
		}
	} ()
	
	for {
		select {
		case event := <-ctrl.EventChan:
			switch event.Name {
			case "PICKUP":
				pickup := event.Data.(Pickup)
				fmt.Printf("[PICK UP] from: %d to: %d\n", pickup.FromFloor, pickup.ToFloor)
				ctrl.Pickup(pickup.FromFloor, pickup.ToFloor)
			case "STATUS_CHANGE":
				var queue []Pickup
				for _, pickup := range ctrl.PickupQueue {
					if !ctrl.Pickup(pickup.FromFloor, pickup.ToFloor) {
						queue = append(queue, pickup)
					}
				}
				ctrl.PickupQueue = queue
			}		
		}
	}
}

func Abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

