package config

import (
	"os"
	"strconv"
)

type Config struct {
	MessageCount  int
	MessageSize   int
	ThreadsAmount int

	AlgName   string
	BenchMode bool
	SaveCsv   bool
}

func LoadConfig() (*Config, error) {
	config := &Config{
		MessageCount:  getNumberFromEnv("MESSAGE_COUNT"),
		MessageSize:   getNumberFromEnv("MESSAGE_SIZE"),
		ThreadsAmount: getNumberFromEnv("THREADS"),
		BenchMode:     getBoolFromEnv("BENCH_MODE"),
		SaveCsv:       getBoolFromEnv("SAVE_CSV"),
	}

	config.AlgName, _ = os.LookupEnv("ALG_NAME")

	return config, nil
}

func getNumberFromEnv(key string) int {
	rawNumber, _ := os.LookupEnv(key)
	parsedNumber, _ := strconv.Atoi(rawNumber)
	return parsedNumber
}

func getBoolFromEnv(key string) bool {
	rawValue, _ := os.LookupEnv(key)
	parsedBool, _ := strconv.ParseBool(rawValue)
	return parsedBool
}
