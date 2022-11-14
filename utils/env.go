package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load environment variables from environment files. Defaults to loading from .env.
func Load(paths ...string) error {
	if len(paths) == 0 {
		paths = []string{".env"}
	}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error reading %v: %w", path, err)
		}
		s := bufio.NewScanner(file)
		var i int
		for s.Scan() {
			i++
			line := s.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) < 2 {
				return fmt.Errorf("missing equal sign on line %v in %v", i, path)
			}
			if err := os.Setenv(parts[0], parts[1]); err != nil {
				return fmt.Errorf("error setting environment variable for line %v in %v: %w", i, path, err)
			}
		}
		if err := s.Err(); err != nil {
			return fmt.Errorf("error scanning %v: %w", path, err)
		}
		_ = file.Close()
	}
	return nil
}
