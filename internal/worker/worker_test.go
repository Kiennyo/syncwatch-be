package worker

import (
	"testing"
	"time"
)

//nolint:revive,cognitive-complexity
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
		t.Run(tt.name, func(_ *testing.T) {
			if tt.wantErr {
				Background(func() {
					panic("This is a test panic")
				})
			} else {
				Background(func() {
					time.Sleep(time.Second)
				})
			}
			Wait()
		})
	}
}
