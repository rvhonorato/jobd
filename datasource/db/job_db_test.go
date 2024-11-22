package db

import (
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	// Delete the database after the test
	defer os.RemoveAll(NAME)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "init db",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitDB(); (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
