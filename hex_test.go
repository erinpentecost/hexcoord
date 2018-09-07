package hexcoord_test

import (
	"fmt"
	"testing"

	"github.com/erinpentecost/hexcoord"
	"github.com/stretchr/testify/assert"
)

func TestHexHashIdentity(t *testing.T) {
	o1 := hexcoord.HexOrigin()
	o2 := hexcoord.HexOrigin()
	assert.Equal(t, o1, o2, "Origin copy is not equal.")
}

func TestHexAdd(t *testing.T) {
	h1 := hexcoord.Hex{
		Q: 5,
		R: 5,
	}
	h2 := hexcoord.Hex{
		Q: 9,
		R: 9,
	}
	hsum := h1.Add(h2)
	hexpected := hexcoord.Hex{
		Q: 14,
		R: 14,
	}

	assert.Equal(t, hexpected, hsum, "Hex add is not correct.")
}

func TestHexDistance(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	closeHexes := hexcoord.HexOrigin().HexArea(done, 1)
	for h := range closeHexes {
		if h == hexcoord.HexOrigin() {
			assert.Equal(t, 0, h.DistanceTo(hexcoord.HexOrigin()))
		} else {
			assert.Equal(t, 1, h.DistanceTo(hexcoord.HexOrigin()), fmt.Sprintf("Hex distance to %v is wrong.", h))
		}
	}
}

func testDirectionEquality(t *testing.T, testOrigin hexcoord.Hex) {
	for a := -9; a < 9; a++ {
		if a == 0 {
			continue
		}
		for d := -9; d < 9; d++ {
			dh := hexcoord.HexDirection(d).Multiply(a).Add(testOrigin)

			rh := hexcoord.HexDirection(3 + d).Multiply(-1 * a).Add(testOrigin)

			assert.Equal(t, dh, rh, "Reversed distance hexes are not equal.")

			oh := hexcoord.HexDirection(3 + d).Multiply(a).Add(testOrigin)

			assert.NotEqual(t, dh, oh, fmt.Sprintf("Opposite hexes are equal with d=%v.", d))

			assert.Equal(t, 2*testOrigin.DistanceTo(oh), dh.DistanceTo(oh), "Distance is not expected for opposite hexes.")
		}
	}
}

func TestDirectionEquality(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	testHexes := hexcoord.HexOrigin().HexArea(done, 10)
	for h := range testHexes {
		testDirectionEquality(t, h)
	}
}

func TestFractionalConversion(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	testHexes := hexcoord.HexOrigin().HexArea(done, 10)
	for h := range testHexes {
		frac := h.ToHexFractional()
		recast := frac.ToHex()
		assert.Equal(t, h, recast)
	}
}