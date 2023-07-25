package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQHMarshalBinary(t *testing.T) {
	qh := DNSHeader{
		ID:      0xabcd,
		Flags:   DNS_REPLY | RECURSION_AVAIL | RECURSION_DESIRED,
		QDCount: 1,
		ANCount: 4,
		NSCount: 0,
		ARCount: 0,
	}

	encoded, err := qh.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xab, 0xcd, 0x81, 0x80, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}, encoded)

	unMarshalledHeader, err := NewQueryHeaderFromBinary(encoded)
	assert.NoError(t, err)
	assert.Equal(t, qh, *unMarshalledHeader)
}

func TestFlags(t *testing.T) {
	qh := DNSHeader{
		ID:      0xabcd,
		QDCount: 1,
		ANCount: 4,
		NSCount: 0,
		ARCount: 0,
		Flags:   DNS_REPLY | RECURSION_AVAIL | RECURSION_DESIRED,
	}

	assert.Truef(t, qh.RecursionAvailable(), "Expect RecursionAvailable to be true")
	assert.True(t, qh.RecursionDesired(), "Expect RecursionDesired to be true")
	assert.True(t, qh.Reply(), "Expect Reply to be true")
	assert.False(t, qh.Query(), "Expect Query to be false")
	assert.Equal(t, 0, qh.OPCode())
	assert.Equal(t, 0, qh.ReturnCode())
}

func TestQHUnmarshalBinary(t *testing.T) {
	expected := &DNSHeader{
		ID:      0xc9e2,
		Flags:   DNS_REPLY | OPCODE_QUERY | RECURSION_DESIRED | RECURSION_AVAIL,
		QDCount: 1,
		ANCount: 4,
		NSCount: 0,
		ARCount: 1,
	}
	b := []byte{0xc9, 0xe2, 0x81, 0x80, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01}

	header, err := NewQueryHeaderFromBinary(b)
	assert.NoError(t, err)
	assert.Equal(t, expected, header)
}

func TestDecodeNameOrPointer(t *testing.T) {
	testCases := []struct {
		Expected string
		Index    int
		Raw      []byte
	}{
		{"joshkeller.dev", 0, []byte{0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r', 0x03, 'd', 'e', 'v', 0x00}},
		{"google.com", 0, []byte{0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00}},
		{"joshkeller.dev", 24,
			[]byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x00, 0x00, 0x00,
				/* 24 */ 0xc0, 0x04},
		},
		{"www.joshkeller.dev", 24,
			[]byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x00, 0x00, 0x00,
				/* 24 */ 0x03, 'w', 'w', 'w',
				/* 28 */ 0xc0, 0x04},
		},
	}

	for _, tc := range testCases {
		decoded, err := DecodeNameOrPointer(tc.Raw, tc.Index)
		assert.NoError(t, err)
		assert.Equal(t, tc.Expected, decoded)
	}

}
