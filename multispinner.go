package multispinner

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Spinner represents a multi-spinner instance
type Spinner struct {
	mu           sync.Mutex
	spinners     []*spinnerInfo
	successColor string
	failureColor string
	frequency    time.Duration
	stopChan     chan struct{}
	currentLine  int
	activeCount  int
	startLine    int
	updateChan   chan updateMsg
}

type spinnerInfo struct {
	message string
	index   int
	active  bool
	lineNum int
}

type updateMsg struct {
	index   int
	message string
	action  string // "start", "stop", "error", "update"
}

// Config holds the configuration for creating a new spinner
type Config struct {
	SuccessColor string
	FailureColor string
	Frequency    time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		SuccessColor: "\033[32m", // Green
		FailureColor: "\033[31m", // Red
		Frequency:    100 * time.Millisecond,
	}
}

// Create creates a new spinner instance with the given configuration
func Create(config Config) *Spinner {
	if config.Frequency == 0 {
		config.Frequency = DefaultConfig().Frequency
	}
	if config.SuccessColor == "" {
		config.SuccessColor = DefaultConfig().SuccessColor
	}
	if config.FailureColor == "" {
		config.FailureColor = DefaultConfig().FailureColor
	}

	// Reset terminal state
	fmt.Print("\033[?25h") // Show cursor
	fmt.Print("\033[0m")   // Reset all attributes

	// Save current cursor position
	fmt.Print("\033[s")

	s := &Spinner{
		spinners:     make([]*spinnerInfo, 0),
		successColor: config.SuccessColor,
		failureColor: config.FailureColor,
		frequency:    config.Frequency,
		stopChan:     make(chan struct{}),
		currentLine:  0,
		startLine:    0, // We'll use relative positioning from saved position
		updateChan:   make(chan updateMsg),
	}

	// Start the spinner goroutine
	go s.run()

	return s
}

// run is the main spinner goroutine that handles all updates
func (s *Spinner) run() {
	chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	firstRun := true

	for {
		select {
		case <-s.stopChan:
			// Clean up terminal state before exiting
			fmt.Print("\033[?25h") // Show cursor
			fmt.Print("\033[0m")   // Reset all attributes
			return
		case msg := <-s.updateChan:
			s.mu.Lock()
			switch msg.action {
			case "start":
				if s.activeCount == 0 {
					fmt.Print("\033[?25l") // Hide cursor
				}
				s.activeCount++
			case "stop", "error":
				s.activeCount--
				if s.activeCount == 0 {
					fmt.Print("\033[?25h") // Show cursor
					fmt.Print("\033[0m")   // Reset all attributes
				}
			}
			s.mu.Unlock()
		default:
			s.mu.Lock()
			if s.activeCount > 0 {
				if firstRun {
					firstRun = false
				} else {
					time.Sleep(s.frequency)
				}

				// Update all active spinners
				for _, spinner := range s.spinners {
					if spinner.active {
						// Move to saved position and then down by the spinner's line number
						fmt.Print("\033[u")
						if spinner.lineNum > 0 {
							fmt.Printf("\033[%dB", spinner.lineNum)
						}
						fmt.Printf("\033[K\033[0m%s %s", chars[i], spinner.message)
					}
				}
				i = (i + 1) % len(chars)
			}
			s.mu.Unlock()
		}
	}
}

// Register registers a new spinner and returns its index
func (s *Spinner) Register() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := len(s.spinners)
	s.spinners = append(s.spinners, &spinnerInfo{
		index:   index,
		lineNum: s.currentLine,
	})
	s.currentLine++
	return index
}

// Message updates the message for a specific spinner
func (s *Spinner) Message(index int, message string) {
	s.updateChan <- updateMsg{
		index:   index,
		message: message,
		action:  "update",
	}
}

// Start starts a spinner with the given message
func (s *Spinner) Start(index int, message string) {
	s.mu.Lock()
	if index >= len(s.spinners) {
		s.mu.Unlock()
		return
	}
	s.spinners[index].message = message
	s.spinners[index].active = true
	s.mu.Unlock()

	s.updateChan <- updateMsg{
		index:   index,
		message: message,
		action:  "start",
	}
}

// Stop stops a spinner with a success message
func (s *Spinner) Stop(index int, message string) {
	s.mu.Lock()
	if index >= len(s.spinners) {
		s.mu.Unlock()
		return
	}

	spinner := s.spinners[index]
	if !spinner.active {
		s.mu.Unlock()
		return
	}
	spinner.active = false
	s.mu.Unlock()

	s.updateChan <- updateMsg{
		index:   index,
		message: message,
		action:  "stop",
	}

	// Print the final message
	fmt.Print("\033[u")
	if spinner.lineNum > 0 {
		fmt.Printf("\033[%dB", spinner.lineNum)
	}
	fmt.Printf("\033[K%s %s\033[0m\n", s.successColor+"✓"+s.successColor, message)
	os.Stdout.Sync()
}

// StopWithError stops a spinner with an error message
func (s *Spinner) StopWithError(index int, error string) {
	s.mu.Lock()
	if index >= len(s.spinners) {
		s.mu.Unlock()
		return
	}

	spinner := s.spinners[index]
	if !spinner.active {
		s.mu.Unlock()
		return
	}
	spinner.active = false
	s.mu.Unlock()

	s.updateChan <- updateMsg{
		index:   index,
		message: error,
		action:  "error",
	}

	// Print the final message
	fmt.Print("\033[u")
	if spinner.lineNum > 0 {
		fmt.Printf("\033[%dB", spinner.lineNum)
	}
	fmt.Printf("\033[K%s %s\033[0m\n", s.failureColor+"✗"+s.failureColor, error)
	os.Stdout.Sync()
}
