package types_test

import (
	"testing"

	. "github.com/jsnfwlr/filamate/internal/types"
)

func TestPointerOf(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		str := "test"
		ptr := PointerOf(str)

		if ptr == nil {
			t.Error("expected non-nil pointer")
		}

		if *ptr != str {
			t.Errorf("expected %s, got %s", str, *ptr)
		}
	})

	t.Run("int", func(t *testing.T) {
		val := 42
		ptr := PointerOf(val)

		if ptr == nil {
			t.Error("expected non-nil pointer")
		}

		if *ptr != val {
			t.Errorf("expected %d, got %d", val, *ptr)
		}
	})

	t.Run("bool", func(t *testing.T) {
		val := true
		ptr := PointerOf(val)

		if ptr == nil {
			t.Error("expected non-nil pointer")
		}

		if *ptr != val {
			t.Errorf("expected %t, got %t", val, *ptr)
		}
	})

	t.Run("empty_string", func(t *testing.T) {
		str := ""
		ptr := PointerOf(str)

		if ptr == nil {
			t.Error("expected non-nil pointer")
		}

		if *ptr != str {
			t.Errorf("expected empty string, got %s", *ptr)
		}
	})

	t.Run("zero_values", func(t *testing.T) {
		zeroInt := PointerOf(0)
		if *zeroInt != 0 {
			t.Errorf("expected 0, got %d", *zeroInt)
		}

		zeroFloat := PointerOf(0.0)
		if *zeroFloat != 0.0 {
			t.Errorf("expected 0.0, got %f", *zeroFloat)
		}

		zeroBool := PointerOf(false)
		if *zeroBool != false {
			t.Errorf("expected false, got %t", *zeroBool)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type testStruct struct {
			Field string
		}

		val := testStruct{Field: "test"}
		ptr := PointerOf(val)

		if ptr == nil {
			t.Error("expected non-nil pointer")
		}

		if ptr.Field != val.Field {
			t.Errorf("expected %s, got %s", val.Field, ptr.Field)
		}
	})

	t.Run("slice", func(t *testing.T) {
		val := []int{1, 2, 3}
		ptr := PointerOf(val)

		if ptr == nil {
			t.Error("expected non-nil pointer")
		}

		if len(*ptr) != len(val) {
			t.Errorf("expected length %d, got %d", len(val), len(*ptr))
		}

		for i, v := range val {
			if (*ptr)[i] != v {
				t.Errorf("expected %d at index %d, got %d", v, i, (*ptr)[i])
			}
		}
	})
}
