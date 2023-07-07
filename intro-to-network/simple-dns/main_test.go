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

	var unMarshalledHeader DNSHeader
	err = unMarshalledHeader.UnmarshalBinary(encoded)
	assert.NoError(t, err)
	assert.Equal(t, qh, unMarshalledHeader)
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
