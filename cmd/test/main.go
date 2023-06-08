package main

import "time"

type T int

func main() {
	var iter = 100

	var slice []T

	queue := make(chan T)

	// Create our data and send it into the queue.
	for i := 0; i < iter; i++ {
		go func(i int) {
			// random sleep time
			time.Sleep(time.Duration(i) * time.Millisecond)
			queue <- T(i)
		}(i)
	}

	for i := 0; i < iter; i++ {
		t := <-queue
		slice = append(slice, t)
	}

	// now prints off all 100 int values
	println(len(slice))
	// fmt.Println(slice)
}
