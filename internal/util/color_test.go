package util

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor2Hex(t *testing.T) {
	c1 := color.RGBA{
		R: 255,
		B: 255,
		G: 255,
		A: 255,
	}

	c2 := color.RGBA{
		R: 125,
		B: 125,
		G: 125,
		A: 255,
	}

	c3 := color.RGBA{
		R: 200,
		B: 10,
		G: 125,
		A: 255,
	}

	c4 := color.RGBA{
		R: 0,
		B: 0,
		G: 0,
		A: 255,
	}

	assert.Equal(t, "#ffffff", ConvertColor2Hex(c1))
	assert.Equal(t, "#7d7d7d", ConvertColor2Hex(c2))
	assert.Equal(t, "#c87d0a", ConvertColor2Hex(c3))
	assert.Equal(t, "#000000", ConvertColor2Hex(c4))
}
