package varint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const MAXUINT64 = ^uint64(0)

func TestEncode(t *testing.T) {
	assert.Equal(t, []byte{0x0}, Encode(0))
	assert.Equal(t, []byte{0x1}, Encode(1))
	assert.Equal(t, []byte{0x96, 0x01}, Encode(150))
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}, Encode(MAXUINT64))
}

func TestDecode(t *testing.T) {
	assert.Equal(t, uint64(1), Decode([]byte{0x01}))
	assert.Equal(t, uint64(150), Decode([]byte{0x96, 0x01}))
	assert.Equal(t, uint64(MAXUINT64), Decode([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}))
}

func TestRoundTrip(t *testing.T) {
	for i := uint64(0); i < 1<<24; i++ {
		assert.Equal(t, i, Decode(Encode(i)))
	}
}

func TestRead(t *testing.T) {
	assert.Equal(t, uint64(1), Read("./1.uint64"))
	assert.Equal(t, uint64(0x96), Read("./150.uint64"))
	assert.Equal(t, uint64(0xffffffffffffffff), Read("./maxint.uint64"))
}
