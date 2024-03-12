package worker

import (
	"fmt"
	"log/slog"
	"sync"
)

type Worker struct {
	wg *sync.WaitGroup
}

func New(wg *sync.WaitGroup) *Worker {
	return &Worker{
		wg: wg,
	}
}

func (w *Worker) Background(fn func()) {
	w.wg.Add(1)

	// Launch a background goroutine.
	go func() {
		defer w.wg.Done()

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
