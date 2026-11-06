package db

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestMissingTask(t *testing.T) {
	if err := Init(filepath.Join(t.TempDir(), "scheduler.db")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(Close)
	for _, tc := range []struct {
		name string
		call func() error
	}{
		{"get", func() error { _, err := GetTask("1"); return err }},
		{"update", func() error { return UpdateTask(&Task{ID: "1"}) }},
		{"delete", func() error { return DeleteTask("1") }},
		{"update date", func() error { return UpdateDate("20270101", "1") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, ErrNotFound) {
				t.Fatalf("error = %v, want ErrNotFound", err)
			}
		})
	}
}
