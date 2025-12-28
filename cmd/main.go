package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"in-memory-key-value-db/internal/app/compute"
	"in-memory-key-value-db/internal/app/pipeline"
	"in-memory-key-value-db/internal/storage/engine"
	"os"
	"strings"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	storage := engine.NewEngine(logger)
	processor := compute.NewCompute()
	pl := pipeline.NewPipeline(processor, storage)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Simple KV store started.")

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		err := pl.Execute(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}
}
