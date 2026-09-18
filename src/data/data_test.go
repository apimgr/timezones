package data

import (
	"errors"
	"os"
	"testing"
)

// TestReadFile covers the happy path (an embedded file that exists) and the
// error path (a file that was never embedded).
func TestReadFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{name: "existing embedded file", file: "timezones.json", wantErr: false},
		{name: "missing file", file: "does-not-exist.json", wantErr: true},
		{name: "empty name", file: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ReadFile(tt.file)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ReadFile(%q) expected error, got nil", tt.file)
				}
				if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, os.ErrInvalid) {
					// embed.FS returns *fs.PathError wrapping fs.ErrNotExist
					// or fs.ErrInvalid; just ensure we got *some* error above.
					return
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadFile(%q) unexpected error: %v", tt.file, err)
			}
			if len(data) == 0 {
				t.Fatalf("ReadFile(%q) returned empty data", tt.file)
			}
		})
	}
}
