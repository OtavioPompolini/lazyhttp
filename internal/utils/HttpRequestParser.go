package utils

import (
	"errors"
	"net/http"
	"strings"
)

// This was made with AI... Will refactor later
func ParseHttpRequest(input string) (*http.Request, error) {
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")

	if len(lines) == 0 {
		return nil, errors.New("empty request")
	}

	requestLine := strings.TrimSpace(lines[0])
	parts := strings.Fields(requestLine)
	rawURL := parts[1]

	headers := make(map[string]string)
	var bodyLines []string
	foundBlankLine := false

	for _, line := range lines[1:] {
		line = strings.TrimRight(line, "\r")
		if !foundBlankLine {
			if strings.TrimSpace(line) == "" {
				foundBlankLine = true
				continue
			}

			headerParts := strings.SplitN(line, ":", 2)
			if len(headerParts) != 2 {
				continue
			}

			key := strings.TrimSpace(headerParts[0])
			value := strings.TrimSpace(headerParts[1])
			headers[key] = value
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	body := strings.Join(bodyLines, "\n")

	httpRequest, err := http.NewRequest(parts[0], rawURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		httpRequest.Header.Add(k, v)
	}

	return httpRequest, nil
}
