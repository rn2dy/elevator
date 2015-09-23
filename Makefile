PROGRAM=ElevatorControl

all: build

build:
	go build -o $(PROGRAM)

run: build
	./$(PROGRAM)
