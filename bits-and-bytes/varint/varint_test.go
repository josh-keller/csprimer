package varint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	assert.Equal(t, []byte{0x0}, Encode(0))
	assert.Equal(t, []byte{0x1}, Encode(1))
	assert.Equal(t, []byte{0x96, 0x01}, Encode(150))
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}, Encode(^uint64(0)))
}

func TestDecode(t *testing.T) {
	assert.Equal(t, 1, Decode(Encode(1)))
}

func TestRead(t *testing.T) {
	assert.Equal(t, uint64(1), Read("./1.uint64"))
	assert.Equal(t, uint64(0x96), Read("./150.uint64"))
	assert.Equal(t, uint64(0xffffffffffffffff), Read("./maxint.uint64"))
}
