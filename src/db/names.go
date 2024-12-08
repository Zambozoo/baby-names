package db

import (
	"bufio"
	"fmt"
	"os"

	"math/rand"
)

func ReadNames(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open names file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var names []string
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read names file: %w", err)
	}

	for i := range names {
		j := rand.Intn(i + 1)
		names[i], names[j] = names[j], names[i]
	}

	return names, nil
}
