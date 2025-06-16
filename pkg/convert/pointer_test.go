package convert_test

import (
	"image"
	"testing"

	"github.com/bmcszk/fogofgo/pkg/convert"
)

func TestToPointer(t *testing.T) {
	t.Run("converts int to pointer", func(t *testing.T) {
		testIntToPointer(t)
	})

	t.Run("converts string to pointer", func(t *testing.T) {
		testStringToPointer(t)
	})

	t.Run("converts struct to pointer", func(t *testing.T) {
		testStructToPointer(t)
	})

	t.Run("converts zero value to pointer", func(t *testing.T) {
		testZeroValueToPointer(t)
	})
}

func testIntToPointer(t *testing.T) {
	t.Helper()
	value := 42
	ptr := convert.ToPointer(value)

	if ptr == nil {
		t.Fatal("expected pointer, got nil")
	}
	if *ptr != value {
		t.Errorf("expected %d, got %d", value, *ptr)
	}
}

func testStringToPointer(t *testing.T) {
	t.Helper()
	value := "hello"
	ptr := convert.ToPointer(value)

	if ptr == nil {
		t.Fatal("expected pointer, got nil")
	}
	if *ptr != value {
		t.Errorf("expected %s, got %s", value, *ptr)
	}
}

func testStructToPointer(t *testing.T) {
	t.Helper()
	value := image.Rect(1, 2, 3, 4)
	ptr := convert.ToPointer(value)

	if ptr == nil {
		t.Fatal("expected pointer, got nil")
	}
	if *ptr != value {
		t.Errorf("expected %+v, got %+v", value, *ptr)
	}
}

func testZeroValueToPointer(t *testing.T) {
	t.Helper()
	var value int
	ptr := convert.ToPointer(value)

	if ptr == nil {
		t.Fatal("expected pointer, got nil")
	}
	if *ptr != 0 {
		t.Errorf("expected 0, got %d", *ptr)
	}
}
