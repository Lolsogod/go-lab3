package tester

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type TestResult struct {
	MsgAmount    int
	MsgSize      int
	Threads      int
	Duration     time.Duration
	Deadlocked   bool
	MessagesLost int
}

type Buffer interface {
	Write(data interface{})
	Read() interface{}
}

type TestConfig struct {
	Name         string
	Buffer       Buffer
	MessageCount int
	MessageSize  int
	ThreadCount  int
	WaitGroup    *sync.WaitGroup
}

func RunBufferTest(config TestConfig) {
	fmt.Printf("\n=== Testing %s ===\n", config.Name)

	message := make([]byte, config.MessageSize)
	for i := range message {
		message[i] = 'A'
	}

	done := make(chan bool)

	go func() {
		messagesPerThread := config.MessageCount / config.ThreadCount
		extraMessages := config.MessageCount % config.ThreadCount

		for i := 0; i < config.ThreadCount; i++ {
			threadMessages := messagesPerThread
			if i < extraMessages {
				threadMessages++
			}

			go func(id int, count int) {
				for j := 0; j < count; j++ {
					config.Buffer.Write(message)
				}
			}(i, threadMessages)
		}

		for i := 0; i < config.ThreadCount; i++ {
			threadMessages := messagesPerThread
			if i < extraMessages {
				threadMessages++
			}

			go func(id int, count int) {
				for j := 0; j < count; j++ {
					config.Buffer.Read()
				}
			}(i, threadMessages)
		}

		if config.WaitGroup != nil {
			config.WaitGroup.Wait()
		}
		done <- true
	}()

	select {
	case <-done:
		fmt.Printf("=== %s completed successfully ===\n", config.Name)
	case <-time.After(5 * time.Second):
		fmt.Printf("=== %s appears deadlocked - moving on... ===\n", config.Name)
	}
}

func PrintTestResults(result TestResult) string {
	var output strings.Builder

	fmt.Fprintf(&output, "\nBuffer Test Results:\n")
	fmt.Fprintf(&output, "==================\n")
	fmt.Fprintf(&output, "Messages: %d\n", result.MsgAmount)
	fmt.Fprintf(&output, "Message Size: %d bytes\n", result.MsgSize)
	fmt.Fprintf(&output, "Threads: %d\n", result.Threads)
	fmt.Fprintf(&output, "Duration: %v\n", result.Duration)

	if result.Deadlocked {
		fmt.Fprintf(&output, "Status: DEADLOCKED\n")
	} else {
		fmt.Fprintf(&output, "Status: OK\n")
	}

	if result.MessagesLost > 0 {
		fmt.Fprintf(&output, "Messages Lost: %d (%.2f%%)\n",
			result.MessagesLost,
			float64(result.MessagesLost)/float64(result.MsgAmount)*100)
	}

	return output.String()
}
