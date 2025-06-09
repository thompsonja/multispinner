# MultiSpinner

A Go library for creating and managing multiple spinners in parallel using goroutines.

## Features

- Create multiple spinners that can run concurrently
- Customize colors for success and failure states
- Adjust spinner frequency
- Thread-safe operations
- Simple and intuitive API

## Installation

```bash
go get github.com/thompsonja/multispinner
```

## Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/thompsonja/multispinner"
)

func main() {
    // Create a new spinner with default configuration
    spinner := multispinner.Create(multispinner.DefaultConfig())

    // Register two spinners
    index1 := spinner.Register()
    index2 := spinner.Register()

    // Start the spinners
    spinner.Start(index1, "Processing task 1...")
    spinner.Start(index2, "Processing task 2...")

    // Simulate some work
    time.Sleep(2 * time.Second)

    // Update message for spinner 1
    spinner.Message(index1, "Almost done with task 1...")

    // Simulate more work
    time.Sleep(1 * time.Second)

    // Stop the spinners
    spinner.Stop(index1, "Task 1 completed successfully!")
    spinner.StopWithError(index2, fmt.Errorf("Task 2 failed: connection timeout"))
}
```

## API Reference

### Create

Creates a new spinner instance with the given configuration.

```go
spinner := multispinner.Create(multispinner.Config{
    SuccessColor: "\033[32m", // Green
    FailureColor: "\033[31m", // Red
    Frequency:    100 * time.Millisecond,
})
```

### Register

Registers a new spinner and returns its index.

```go
index := spinner.Register()
```

### Message

Updates the message for a specific spinner.

```go
spinner.Message(index, "New message")
```

### Start

Starts a spinner with the given message.

```go
spinner.Start(index, "Starting task...")
```

### Stop

Stops a spinner with a success message.

```go
spinner.Stop(index, "Task completed successfully!")
```

### StopWithError

Stops a spinner with an error message.

```go
spinner.StopWithError(index, fmt.Errorf("Task failed: %v", err))
```

