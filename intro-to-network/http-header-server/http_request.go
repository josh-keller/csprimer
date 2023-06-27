package main

import (
	"errors"
	"regexp"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func ParseRequest(raw []byte) (*HTTPRequest, error) {
	reqLineRE := regexp.MustCompile(`(?ms)^\s*(\w+)\s+([\w.\-\/]+)\s+(HTTP\/\d\.\d)\s*\r\n(.*)\r\n\r\n(.*)`)
	match := reqLineRE.FindSubmatch(raw)
	if len(match) == 0 {
		return nil, errors.New("No match")
	}
	return &HTTPRequest{
		Method:  string(match[1]),
		Path:    string(match[2]),
		Version: string(match[3]),
		Headers: parseHeaders(string(match[4])),
		Body:    string(match[5]),
	}, nil
}

func parseHeaders(raw string) map[string]string {
	headerRE := regexp.MustCompile(`(?mi)^\s*([\w\-]+):\s+(.*)$`)
	matches := headerRE.FindAllStringSubmatch(raw, -1)
	headers := make(map[string]string)
	for _, match := range matches {
		headers[match[1]] = strings.TrimSpace(match[2])
	}

	return headers
}
