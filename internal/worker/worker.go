package worker

import (
	"fmt"
	"log/slog"
	"sync"
)

var wg sync.WaitGroup

func Background(fn func()) {
	wg.Add(1)

	// Launch a background goroutine.
	go func() {
		defer wg.Done()

		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				slog.Error(fmt.Errorf("failed to execute background task: %s", err).Error())
			}
		}()

		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}

func Wait() {
	wg.Wait()
}
