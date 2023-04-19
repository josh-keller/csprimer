package color_convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertString(t *testing.T) {
	assert.Equal(t, "rgb(254, 3, 10)", ConvertString("#fe030a"))
	assert.Equal(t, "rgb(15, 13, 239)", ConvertString("#0f0def"))

}

func TestHexToRGB(t *testing.T) {
	assert.Equal(t, [3]int{254, 3, 10}, HexToRGB(0xfe030a))
	assert.Equal(t, [3]int{15, 13, 239}, HexToRGB(0x0f0def))
}
