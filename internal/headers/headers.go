package headers

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for k, v := range *h {
		cb(k, v)
	}
}

func (h *Headers) Get(key string) string {
	value := (*h)[key]
	return value
}

func validateFieldName(name string) (string, error) {
	lowerCase := strings.ToLower(name)

	matched, err := regexp.MatchString(`^[\w!#$%&'*+\-.^_`+"`"+`|~]+$`, name)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", nil
	}
	return lowerCase, nil
}

func parseHeader(fieldLine []byte) (string, string, error) {
	splits := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(splits) != 2 {
		return "", "", errors.New("invalid field line")
	}
	fmt.Println(string(splits[0]), string(splits[1]))
	name := (string(splits[0]))
	value := strings.TrimSpace(string(splits[1]))
	name, err := validateFieldName(name)
	if err != nil {
		return "", "", err
	}
	if name == "" {
		return "", "", errors.New("invalid field name")
	}

	fmt.Println("FINAL ++++", name, value)
	return name, value, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	fmt.Println("Headers Parse ++++", string(data))
	lineSepratorBytes := []byte("\r\n")
	speratorIndex := bytes.Index(data, lineSepratorBytes)
	DONE := false
	for {
		if speratorIndex == -1 {
			fmt.Println("/n not found", len(data))
			break
		}
		if speratorIndex == 0 {
			DONE = true
			return len(lineSepratorBytes), DONE, nil
		}

		name, value, err := parseHeader(data[:speratorIndex])
		if err != nil {
			return 0, false, err
		}
		if exisitingValue, ok := h[name]; ok {
			exisitingValue += ", " + value
			h[name] = exisitingValue
			fmt.Println("H value +++++", h[value])
		} else {
			h[name] = value
		}
		bytesConsumed := speratorIndex + len(lineSepratorBytes)
		if len(data) >= bytesConsumed+len(lineSepratorBytes) {
			if bytes.Equal(data[bytesConsumed:bytesConsumed+len(lineSepratorBytes)], lineSepratorBytes) {
				return bytesConsumed + len(lineSepratorBytes), true, nil
			}
		}

		return bytesConsumed, false, nil
	}

	return 0, DONE, nil
}
