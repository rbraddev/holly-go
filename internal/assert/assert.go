package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; expected: %v", actual, expected)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}

func ExpError(t *testing.T, actual error) {
	t.Helper()

	if actual == nil {
		t.Errorf("got: nil; expected: %v", actual)
	}
}
