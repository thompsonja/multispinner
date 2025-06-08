package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/thompsonja/multispinner"
)

func main() {
	// Create a new spinner with default configuration
	spinner := multispinner.Create(multispinner.DefaultConfig())

	// Register spinners for each task
	spinner1 := spinner.Register()
	spinner2 := spinner.Register()
	spinner3 := spinner.Register()

	// Start all spinners
	spinner.Start(spinner1, "Task 1: Processing...")
	spinner.Start(spinner2, "Task 2: Processing...")
	spinner.Start(spinner3, "Task 3: Processing...")

	// Create a WaitGroup to wait for all tasks
	var wg sync.WaitGroup
	wg.Add(3) // We have 3 tasks

	// Run tasks in goroutines
	go func() {
		defer wg.Done()
		// Simulate task 1
		time.Sleep(2 * time.Second)
		spinner.Stop(spinner1, "Task 1: Completed successfully")
	}()

	go func() {
		defer wg.Done()
		// Simulate task 2
		time.Sleep(3 * time.Second)
		spinner.Stop(spinner2, "Task 2: Completed successfully")
	}()

	go func() {
		defer wg.Done()
		// Simulate task 3
		time.Sleep(1 * time.Second)
		spinner.StopWithError(spinner3, "Task 3: Failed with error")
	}()

	// Wait for all tasks to complete
	wg.Wait()
	fmt.Println("\nAll tasks completed!")
}
