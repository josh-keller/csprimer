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
		Expected     string
		ExpectedNext int
		Index        int
		Raw          []byte
	}{
		{"joshkeller.dev", 16, 0, []byte{0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r', 0x03, 'd', 'e', 'v', 0x00, 0x00, 0x01, 0x00, 0x01}},
		{"google.com", 12, 0, []byte{0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, 0x00, 0x01, 0x00, 0x01}},
		{"joshkeller.dev", 26, 24,
			[]byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x00, 0x00, 0x00,
				/* 24 */ 0xc0, 0x04,
				/* QType, Class */ 0x00, 0x01, 0x00, 0x01},
		},
		{"www.joshkeller.dev", 30, 24,
			[]byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x00, 0x00, 0x00,
				/* 24 */ 0x03, 'w', 'w', 'w',
				/* 28 */ 0xc0, 0x04,
				/* QType, Class */ 0x00, 0x01, 0x00, 0x01},
		},
	}

	for _, tc := range testCases {
		decodedName, nextptr, err := DecodeNameOrPointer(tc.Raw, tc.Index)
		assert.NoError(t, err)
		assert.Equal(t, tc.Expected, decodedName)
		assert.Equal(t, tc.ExpectedNext, nextptr)
	}
}

func TestDecodeQuestionSection(t *testing.T) {
	testCases := []struct {
		Expected QuestionSection
		Index    int
		Raw      []byte
	}{
		{QuestionSection{"joshkeller.dev", 1, 1}, 0, []byte{0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r', 0x03, 'd', 'e', 'v', 0x00, 0x00, 0x01, 0x00, 0x01}},
		{QuestionSection{"google.com", 1, 1}, 0, []byte{0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, 0x00, 0x01, 0x00, 0x01}},
		{QuestionSection{"joshkeller.dev", 1, 1}, 24,
			[]byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x00, 0x00, 0x00,
				/* 24 */ 0xc0, 0x04,
				/* QType, Class */ 0x00, 0x01, 0x00, 0x01},
		},
		{QuestionSection{"www.joshkeller.dev", 1, 1}, 24,
			[]byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x00, 0x00, 0x00,
				/* 24 */ 0x03, 'w', 'w', 'w',
				/* 28 */ 0xc0, 0x04,
				/* QType, Class */ 0x00, 0x01, 0x00, 0x01},
		},
	}

	for _, tc := range testCases {
		decoded, _, err := DecodeQuestionSection(tc.Raw, tc.Index)
		assert.NoError(t, err)
		assert.Equal(t, tc.Expected, *decoded)
	}
}

func TestDecodeResourceRecord(t *testing.T) {
	testCases := []struct {
		Expected ResourceRecord
		Index    int
		Raw      []byte
	}{
		{Expected: ResourceRecord{"joshkeller.dev", 1, 1, 1, 4, []byte{1, 2, 3, 4}},
			Index: 0,
			Raw: []byte{
				0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r', 0x03, 'd', 'e', 'v', 0x00,
				0x00, 0x01, // Type
				0x00, 0x01, // Class
				0x00, 0x00, 0x00, 0x01, // TTL
				0x00, 0x04, // RDLength
				0x01, 0x02, 0x03, 0x04, // RData
			}},
		{Expected: ResourceRecord{"google.com", 1, 1, 1, 4, []byte{5, 6, 7, 8}},
			Index: 0,
			Raw: []byte{
				0x06, 'g', 'o', 'o', 'g', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00,
				0x00, 0x01,
				0x00, 0x01,
				0x00, 0x00, 0x00, 0x01, // TTL
				0x00, 0x04, // RDLength
				0x05, 0x06, 0x07, 0x08, // RData
			}},
		{Expected: ResourceRecord{"joshkeller.dev", 1, 1, 1, 4, []byte{1, 1, 1, 1}},
			Index: 24,
			Raw: []byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x00, 0x00, 0x00,
				/* 24 */ 0xc0, 0x04,
				/* Type, Class */ 0x00, 0x01, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x01, // TTL
				0x00, 0x04, // RDLength
				0x01, 0x01, 0x01, 0x01, // RData
			},
		},
		{Expected: ResourceRecord{"www.joshkeller.dev", 1, 1, 1, 4, []byte{1, 1, 1, 1}},
			Index: 24,
			Raw: []byte{0x00, 0x00, 0x00, 0x00,
				/* 4 */ 0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
				/* 15 */ 0x03, 'd', 'e', 'v',
				/* 19 */ 0x00, 0x00, 0x01, 0x00, 0x01,
				/* 24 */ 0x03, 'w', 'w', 'w',
				/* 28 */ 0xc0, 0x04,
				/* Type, Class */ 0x00, 0x01, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x01, // TTL
				0x00, 0x04, // RDLength
				0x01, 0x01, 0x01, 0x01, // RData
			},
		},
	}

	for _, tc := range testCases {
		decoded, _, err := DecodeResourceRecord(tc.Raw, tc.Index)
		assert.NoError(t, err)
		assert.Equal(t, tc.Expected, *decoded)
	}
}
