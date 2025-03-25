package assert

import "testing"

// NilErr asserts that the given err value is nil.
func NilErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("received non-nil error: %s", err)
	}
}
