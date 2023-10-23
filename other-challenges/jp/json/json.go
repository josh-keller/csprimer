package json

import (
	"errors"
	"fmt"
)

type Obj interface{}

var ErrMissingClosingBrace = errors.New("No matching '}' found.")

func Parse(b []byte) (Obj, error) {
	return parseJson(b, 0)
}

func parseJson(b []byte, i int) (Obj, error) {
	// Next should be an opening char
	i, err := consumeWhitespace(b, i)
	if err != nil {
		return nil, errors.New("File is only whitepace")
	}

	switch b[i] {
	case '{':
		return parseObject(b, i+1)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid character at position %d", i))
	}
}

func consumeWhitespace(b []byte, i int) (int, error) {
	for ; i < len(b); i++ {
		if !isWhiteSpace(b[i]) {
			return i, nil
		}
	}

	return i, errors.New("Reached end of file")
}

func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func parseObject(b []byte, i int) (Obj, error) {
	i, err := consumeWhitespace(b, i)
	if err != nil {
		return nil, errors.New("Unexpected end of file")
	}

	switch b[i] {
	case '}':
		return struct{}{}, nil
	default:
		return nil, ErrMissingClosingBrace
	}
}
