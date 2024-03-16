package worker

import (
	"sync"
	"testing"
	"time"
)

func TestBackground(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Normal Function",
			wantErr: false,
		},
		{
			name:    "Panic Prone Function",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			w := New(&wg)

			if tt.wantErr {
				w.Background(func() {
					panic("This is a test panic")
				})
			} else {
				w.Background(func() {
					time.Sleep(time.Second)
				})
			}

			wg.Wait()
		})
	}
}
