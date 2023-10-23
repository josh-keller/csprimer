package json

import (
	"errors"
	"fmt"
	"io"
)

type Obj interface{}

var ErrMissingClosingBrace = errors.New("No matching '}' found.")

func Parse(r io.ByteScanner) (Obj, error) {
	return parseJson(r, 0)
}

func parseJson(r io.ByteScanner, i int) (Obj, error) {
	// Next should be an opening char
	b, err := consumeWhitespace(r)
	if err != nil {
		return nil, errors.New("File is only whitepace")
	}

	switch b {
	case '{':
		return parseObject(r)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid character at position %d", i))
	}
}

func consumeWhitespace(r io.ByteScanner) (byte, error) {
	var b byte
	var err error

	for b, err = r.ReadByte(); err == nil; b, err = r.ReadByte() {
		if !isWhiteSpace(b) {
			return b, nil
		}
	}

	return b, errors.New("Reached end of file")
}

func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func parseObject(r io.ByteScanner) (Obj, error) {
	b, err := consumeWhitespace(r)
	if err != nil {
		return nil, errors.Join(errors.New("Parse object"), err)
	}

	switch b {
	case '}':
		return struct{}{}, nil
	default:
		return nil, ErrMissingClosingBrace
	}
}
