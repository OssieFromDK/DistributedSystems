package main

import (
	"fmt"
	"time"
)

type ChopstickChV2 struct {
	id int
	request chan chan bool
	put chan bool
}

func ChopstickV2(id int) ChopstickChV2 {
	return ChopstickChV2{
		id: id,
		request: make(chan chan bool),
		put: make(chan bool),
	}
}

func RunChopstick(stick ChopstickChV2) {
	for {
		replyChan := <-stick.request
		replyChan <- true
		<-stick.put
	}
}

func PhilosopherV2(name string, leftID, rightID int, leftStick, rightStick ChopstickChV2) {
	for i := 0; i < 3; i++ {
		fmt.Println("Philosopher", name, "is thinking")
		time.Sleep(300 * time.Millisecond)

		if leftID < rightID {
			takeStick(leftStick)
			fmt.Println("Philosopher", name, "took left chopstick", leftID)

			takeStick(rightStick)
			fmt.Println("Philosopher", name, "took right chopstick", rightID)
		} else {
			takeStick(rightStick)
			fmt.Println("Philosopher", name, "took right chopstick", rightID)

			takeStick(leftStick)
			fmt.Println("Philosopher", name, "took left chopstick", leftID)
		}

		fmt.Println("Philosopher", name, "is eating (", i+1, "/3 )")
		time.Sleep(300 * time.Millisecond)

		putStick(leftStick)
		putStick(rightStick)
		fmt.Println("Philosopher", name, "put both chopsticks back")
	}

	fmt.Println("Philosopher", name, "is done eating")
}

func takeStick(stick ChopstickChV2) {
	reply := make(chan bool)
	stick.request <- reply
	<-reply
}

func putStick(stick ChopstickChV2) {
	stick.put <- true
}

func main() {
	stick1 := ChopstickV2(0)
	stick2 := ChopstickV2(1)
	stick3 := ChopstickV2(2)
	stick4 := ChopstickV2(3)
	stick5 := ChopstickV2(4)

	go RunChopstick(stick1)
	go RunChopstick(stick2)
	go RunChopstick(stick3)
	go RunChopstick(stick4)
	go RunChopstick(stick5)

	go PhilosopherV2("F1", 0, 1, stick1, stick2)
	go PhilosopherV2("F2", 1, 2, stick2, stick3)
	go PhilosopherV2("F3", 2, 3, stick3, stick4)
	go PhilosopherV2("F4", 3, 4, stick4, stick5)
	go PhilosopherV2("F5", 4, 0, stick5, stick1)

	time.Sleep(9 * time.Second)
}
