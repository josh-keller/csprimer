package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountBytes(t *testing.T) {
	assert.Equal(t, 341836, CountBytes("test.txt"))
}

func TestCountLines(t *testing.T) {
	assert.Equal(t, 7137, CountLines("test.txt"))
}
