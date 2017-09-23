package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRndStr(t *testing.T) {
	assert := assert.New(t)
	for i := 1; i < 512; i++ {
		assert.Len(RndStr(i), i)
	}
}

func TestRndFloat(t *testing.T) {
	assert := assert.New(t)
	for i := 1; i < 1*1000*1000; i++ {
		f := RndFloat()
		assert.True(f <= 1.0, "Must be smaller or equal to 1")
		assert.True(f >= 0.0, "Must be bigger or equal to 0")
	}
}

func TestRndInt64(t *testing.T) {
	assert := assert.New(t)
	for i := 1; i < 1000; i++ {
		assert.NotPanics(func() { RndInt64() })
	}
}
