package engine_test

import (
	"testing"

	. "github.com/DCNT-Hammer/dcnt/engine"
)

func TestRegisterPrometheusTwice(t *testing.T) {
	// prometheus will panic if this fails
	RegisterPrometheus()
	RegisterPrometheus()
}
