package convert_test

import (
	"image"
	"testing"

	"github.com/bmcszk/gptrts/pkg/convert"
)

func TestToPointer(t *testing.T) {
	t.Run("converts int to pointer", func(t *testing.T) {
		value := 42
		ptr := convert.ToPointer(value)

		if ptr == nil {
			t.Fatal("expected pointer, got nil")
		}
		if *ptr != value {
			t.Errorf("expected %d, got %d", value, *ptr)
		}
	})

	t.Run("converts string to pointer", func(t *testing.T) {
		value := "hello"
		ptr := convert.ToPointer(value)

		if ptr == nil {
			t.Fatal("expected pointer, got nil")
		}
		if *ptr != value {
			t.Errorf("expected %s, got %s", value, *ptr)
		}
	})

	t.Run("converts struct to pointer", func(t *testing.T) {
		value := image.Rect(1, 2, 3, 4)
		ptr := convert.ToPointer(value)

		if ptr == nil {
			t.Fatal("expected pointer, got nil")
		}
		if *ptr != value {
			t.Errorf("expected %+v, got %+v", value, *ptr)
		}
	})

	t.Run("converts zero value to pointer", func(t *testing.T) {
		var value int
		ptr := convert.ToPointer(value)

		if ptr == nil {
			t.Fatal("expected pointer, got nil")
		}
		if *ptr != 0 {
			t.Errorf("expected 0, got %d", *ptr)
		}
	})
}
