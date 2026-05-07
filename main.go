package main

import "fmt"

func worker(logs <-chan string) {
	// consumer
	for log := range logs {
		fmt.Println(log)
	}
}

func main() {
	logs := []string{
		"ERROR: DB failed",
		"INFO: user login",
	}

	logChan := make(chan string)

	go worker(logChan)

	// producer
	go func() {
		for _, log := range logs {
			logChan <- log
		}

		close(logChan)
	}()

}
