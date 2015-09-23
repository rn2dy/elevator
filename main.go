package main

import (
	"math/rand"
	"time"
)

var (
	NumOfElevators        = 4               // number of elevators to play with
	PickupFrequency       = 1 * time.Second // pickup request frequency
	ElevatorSpeed         = 4 * time.Second // move to next queued floor every 5 seconds
	StatusUpdateFrequency = 3 * time.Second // status update frequency
	NumOfFloors           = 24
)

// Simulation program
func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// create elevator instances
	var elevators []*Elevator
	for i := 0; i < NumOfElevators; i++ {
		elevators = append(elevators, NewElevator(i))
	}

	// create a controller instance and register elevators
	var controller = NewController()
	controller.RegisterElevators(elevators...)

	// start controller thread
	go func() {
		controller.Start()
	}()

	// start elevators in their own thread
	for _, el := range elevators {
		go func(elevator *Elevator) {
			elevator.Run()
		}(el)
	}

	// send pickup event to controller
	for {
		select {
		case <-time.Tick(PickupFrequency):
			from := rand.Intn(NumOfFloors)
			to := rand.Intn(NumOfFloors)
			if from != to {
				controller.EventChan <- Event{
					Name: "PICKUP",
					Data: Pickup{from, to},
				}	
			}
		}
	}
}
