package goenv_test

import (
	"testing"

	"github.com/go-stack/stack"
)

func TestSigpanic(t *testing.T) {
	t.Parallel()
	sp := goenv.Sigpanic()
	if got, want := sp.Name(), "runtime.sigpanic"; got != want {
		t.Errorf("got == %v, want == %v", got, want)
	}
}
