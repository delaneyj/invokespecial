package tests

import (
	. "invokespecial/pkg/invokespecial"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	p := Seq(Str("a"), Str("b"))
	r, err := Parse(p, "ab")
	assert.Nil(t, err)
	assert.Equal(t, NewPair("a", "b"), r)
}
