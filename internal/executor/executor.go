package executor

import (
	"awesomeProject/internal/nosync"
	"awesomeProject/internal/tester"
	"awesomeProject/internal/withchan"
	"awesomeProject/internal/withmutex"
	"awesomeProject/internal/withwg"
	"time"
)

type TestFunction func(messageCount, messageSize, threadCount int) tester.TestResult

var algorithms = map[string]TestFunction{
	"nosync": func(messageCount, messageSize, threadCount int) tester.TestResult {
		start := time.Now()
		lost := nosync.Test(messageCount, messageSize, threadCount)
		return tester.TestResult{
			MsgAmount:    messageCount,
			MsgSize:      messageSize,
			Threads:      threadCount,
			Duration:     time.Since(start),
			MessagesLost: lost,
		}
	},
	"waitgroup": func(messageCount, messageSize, threadCount int) tester.TestResult {
		start := time.Now()
		withwg.Test(messageCount, messageSize, threadCount)
		duration := time.Since(start)
		return tester.TestResult{
			MsgAmount:  messageCount,
			MsgSize:    messageSize,
			Threads:    threadCount,
			Duration:   duration,
			Deadlocked: duration >= 5*time.Second,
		}
	},
	"channels": func(messageCount, messageSize, threadCount int) tester.TestResult {
		start := time.Now()
		withchan.Test(messageCount, messageSize, threadCount)
		duration := time.Since(start)
		return tester.TestResult{
			MsgAmount:  messageCount,
			MsgSize:    messageSize,
			Threads:    threadCount,
			Duration:   duration,
			Deadlocked: duration >= 5*time.Second,
		}
	},
	"mutex": func(messageCount, messageSize, threadCount int) tester.TestResult {
		start := time.Now()
		withmutex.Test(messageCount, messageSize, threadCount)
		duration := time.Since(start)
		return tester.TestResult{
			MsgAmount:  messageCount,
			MsgSize:    messageSize,
			Threads:    threadCount,
			Duration:   duration,
			Deadlocked: duration >= 5*time.Second,
		}
	},
}

func Execute(algName string, messageCount, messageSize, threadCount int) tester.TestResult {
	algorithm, exists := algorithms[algName]
	if !exists {
		panic("Unrecognized algorithm: " + algName)
	}
	return algorithm(messageCount, messageSize, threadCount)
}
