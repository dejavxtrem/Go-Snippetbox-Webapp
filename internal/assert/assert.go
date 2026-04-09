package assert

import (
	"reflect"
	"testing"
)

func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got == want {
		t.Errorf("got: %v; expected values to be different", got)
	}
}

func Equal[T any](t *testing.T, got, want T) {
	t.Helper()

	// if got != want {
	// 	t.Errorf("got: %v; want: %v", got, want)
	// }

	// And call the isEqual() function below instead of using the != comparison
	// operator.
	if !isEqual(got, want) {
		//t.Errorf("got: %v; want: %v", got, want)
		t.Errorf("got: true; want: false")
	}
}

// Add a helper to check that a value is true.
func True(t *testing.T, got bool) {
	t.Helper()

	if !got {
		t.Errorf("got: %t; want: true", got)
	}
}

// Add a helper to check that a value is nil, using the isNil() function below.
func Nil(t *testing.T, got any) {
	t.Helper()

	// if !isNil(got) {
	// 	t.Errorf("got: %v; want: nil", got)
	// }

	if got == nil {
		t.Errorf("got: nil; want: non-nil")
	}
}

func NotNil(t *testing.T, got any) {
	t.Helper()

	if got == nil {
		t.Errorf("got: nil; want: non-nil")
	}
}

func isEqual[T any](got, want T) bool {
	// First check if both values are nil using the isNil() function below.
	if isNil(got) && isNil(want) {
		return true
	}

	// Otherwise use reflect.DeepEqual to check if they are the same.
	return reflect.DeepEqual(got, want)
}

func isNil(v any) bool {
	// Returns true if v equals nil.
	if v == nil {
		return true
	}

	// Use reflection to check the underlying type of v, and return true if it
	// is a nullable type (e.g. pointer, map or slice) with a value of nil.
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	}

	// Other types like string, bool, int are never nil.
	return false
}
