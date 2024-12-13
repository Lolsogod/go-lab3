package benchmark

import (
	"awesomeProject/internal/nosync"
	"awesomeProject/internal/withchan"
	"awesomeProject/internal/withmutex"
	"awesomeProject/internal/withwg"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type TestResult struct {
	msgAmount    int
	msgSize      int
	threads      int
	noSync       time.Duration
	waitGroup    time.Duration
	channels     time.Duration
	mutex        time.Duration
	deadlocks    map[string]bool
	messagesLost int 
}

func Run() []TestResult {
	messageAmounts := []int{100, 1000, 10000}
	messageSizes := []int{8, 64, 256, 1024}
	threadCounts := []int{2, 3, 4, 6, 8, 10, 12, 50, 100}

	var results []TestResult

	fmt.Println("Performance Testing Results:")
	fmt.Println("============================")

	for _, msgAmount := range messageAmounts {
		for _, msgSize := range messageSizes {
			for _, threads := range threadCounts {
				result := TestResult{
					msgAmount: msgAmount,
					msgSize:   msgSize,
					threads:   threads,
					deadlocks: make(map[string]bool),
				}

				fmt.Printf("\nTesting with: %d messages, %d bytes each, %d threads\n",
					msgAmount, msgSize, threads)
				fmt.Println("----------------------------------------")

				start := time.Now()
				lost := nosync.Test(msgAmount, msgSize, threads) 
				result.noSync = time.Since(start)
				result.messagesLost = lost
				fmt.Printf("No Sync: %v (Lost: %d messages)\n", result.noSync, lost)

				time.Sleep(time.Second)

				start = time.Now()
				withwg.Test(msgAmount, msgSize, threads)
				result.waitGroup = time.Since(start)
				result.deadlocks["waitGroup"] = result.waitGroup >= 5*time.Second
				fmt.Printf("With WaitGroup: %v\n", result.waitGroup)

				time.Sleep(time.Second)

				start = time.Now()
				withchan.Test(msgAmount, msgSize, threads)
				result.channels = time.Since(start)
				result.deadlocks["channels"] = result.channels >= 5*time.Second
				fmt.Printf("With Channels: %v\n", result.channels)

				time.Sleep(time.Second)

				start = time.Now()
				withmutex.Test(msgAmount, msgSize, threads)
				result.mutex = time.Since(start)
				result.deadlocks["mutex"] = result.mutex >= 5*time.Second
				fmt.Printf("With Mutex: %v\n", result.mutex)

				time.Sleep(time.Second)

				results = append(results, result)
			}
		}
	}

	printSummaryTable(results)

	return results
}

func printSummaryTable(results []TestResult) {
	fmt.Println("\n\nSummary Table")
	fmt.Println("=============")

	header := "| Messages | Size (B) | Threads | NoSync (Lost) | WaitGroup | Channels | Mutex |"
	separator := "|----------|----------|----------|--------------|------------|-----------|--------|"
	fmt.Println(header)
	fmt.Println(separator)

	formatDuration := func(d time.Duration, deadlocked bool) string {
		if deadlocked {
			return "DEADLOCK"
		}
		if d == 0 {
			return "~0"
		}
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
	}

	for _, r := range results {
		noSyncResult := fmt.Sprintf("%s (%d lost)",
			formatDuration(r.noSync, false),
			r.messagesLost)

		fmt.Printf("| %-8d | %-8d | %-8d | %-12s | %-10s | %-9s | %-6s |\n",
			r.msgAmount,
			r.msgSize,
			r.threads,
			noSyncResult,
			formatDuration(r.waitGroup, r.deadlocks["waitGroup"]),
			formatDuration(r.channels, r.deadlocks["channels"]),
			formatDuration(r.mutex, r.deadlocks["mutex"]))
	}
}

func ExportToCSV(results []TestResult) error {
	file, err := os.Create("performance_results.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Messages", "Size (B)", "Threads",
		"NoSync (ms)", "Messages Lost",
		"WaitGroup (ms)", "WaitGroup Deadlock",
		"Channels (ms)", "Channels Deadlock",
		"Mutex (ms)", "Mutex Deadlock",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	formatDuration := func(d time.Duration, deadlocked bool) string {
		if deadlocked {
			return "DEADLOCK"
		}
		return fmt.Sprintf("%.2f", float64(d.Microseconds())/1000)
	}

	for _, r := range results {
		row := []string{
			strconv.Itoa(r.msgAmount),
			strconv.Itoa(r.msgSize),
			strconv.Itoa(r.threads),
			formatDuration(r.noSync, false),
			strconv.Itoa(r.messagesLost),
			formatDuration(r.waitGroup, r.deadlocks["waitGroup"]),
			strconv.FormatBool(r.deadlocks["waitGroup"]),
			formatDuration(r.channels, r.deadlocks["channels"]),
			strconv.FormatBool(r.deadlocks["channels"]),
			formatDuration(r.mutex, r.deadlocks["mutex"]),
			strconv.FormatBool(r.deadlocks["mutex"]),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
