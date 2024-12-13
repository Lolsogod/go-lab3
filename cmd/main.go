package main

import (
	"awesomeProject/internal/benchmark"
	"awesomeProject/internal/config"
	"awesomeProject/internal/executor"
	"awesomeProject/internal/tester"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	if cfg.BenchMode {
		results := benchmark.Run()

		if cfg.SaveCsv {
			if err := benchmark.ExportToCSV(results); err != nil {
				fmt.Printf("\nError exporting to CSV: %v\n", err)
			} else {
				fmt.Printf("\nResults exported to performance_results.csv\n")
			}
		}
	} else {
		results := executor.Execute(cfg.AlgName, cfg.MessageCount, cfg.MessageSize, cfg.ThreadsAmount)
		fmt.Println(tester.PrintTestResults(results))
	}

}
