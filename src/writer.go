package main

import (
	"io"
	"os"
	"sync"
)

type OutputWriter struct {
	mu sync.Mutex
	w io.Writer
}

func (w *OutputWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// do something to data

	len, err := os.Stdout.Write([]byte(data))
	if err != nil {
		return len, err
	}
	return len, nil
}
