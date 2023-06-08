package main

import "math/rand"

func main() {
	maxConcurrent := 10

	jobs := make(chan int)
	limit := make(chan struct{}, maxConcurrent)
	done := make(chan bool, 1)

	// randomly queue jobs
	go func() {
		for {
			jobsToCreate := rand.Intn(1000)
			for i := 0; i < jobsToCreate; i++ {
				println("adding jobId", i)
				jobs <- i
			}
			generateRandomSleepTime := rand.Intn(10000)
			// println("sleeping for", generateRandomSleepTime)
			sleepTime := generateRandomSleepTime

			if sleepTime == 100 {
				// println("process done")
				done <- true
				break
			}
		}
	}()

	go func() {
		// process jobs
		i := 0
		for {
			if i >= maxConcurrent {
				i = 0
			}
			i++
			// add to the limit channel. When the channel is full, this blocks the loop
			limit <- struct{}{} // this is a struct because it takes up 0 bytes
			go func(id int) {
				// println("starting go routine", id, "with job")
				processJob(<-jobs, id)
				<-limit
			}(i)
		}
	}()

	<-done

	for {
		if len(jobs) <= 0 || len(limit) <= 0 {
			println("all jobs processed", len(limit), len(jobs))
			break
		}
	}
}

func processJob(jobId int, id int) {
	println("processing jobId", jobId, "with go routine", id)
}
