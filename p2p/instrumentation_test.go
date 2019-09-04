package p2p_test

import (
	"testing"

	. "github.com/DCNT-Hammer/dcnt/p2p"
)

func TestRegisterPrometheus(t *testing.T) {
	RegisterPrometheus()
	RegisterPrometheus()
}
