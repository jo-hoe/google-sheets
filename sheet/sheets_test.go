package sheet

import "testing"

func Test_hasFlag(t *testing.T) {
	flags := O_CREATE | O_EXCL | O_TRUNC
	if !hasFlag(flags, O_CREATE) {
		t.Errorf("%d did not have O_CREATE", flags)
	}
	if !hasFlag(flags, O_EXCL) {
		t.Errorf("%d did not have O_CREATE", flags)
	}
	if !hasFlag(flags, O_TRUNC) {
		t.Errorf("%d did not have O_CREATE", flags)
	}
}

func Test_hasFlag_without_flags(t *testing.T) {
	var flags int
	if hasFlag(flags, O_CREATE) {
		t.Errorf("%d did unexpectedly find O_CREATE", flags)
	}
	if hasFlag(flags, O_EXCL) {
		t.Errorf("%d did unexpectedly find O_CREATE", flags)
	}
	if hasFlag(flags, O_TRUNC) {
		t.Errorf("%d did unexpectedly find O_CREATE", flags)
	}
}
