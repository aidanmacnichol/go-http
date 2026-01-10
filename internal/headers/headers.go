package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) SetHeader(key string, val string) {
	h[key] = val
}

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[strings.ToLower(key)]
	return val, ok
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return len(crlf), true, nil
	}

	fieldName, fieldValue, err := fieldLineFromString(string(data[:idx]))
	if err != nil {
		return 0, false, err
	}

	value, ok := h[fieldName]
	if ok {
		h[fieldName] = value + ", " + fieldValue
	} else {
		h[fieldName] = fieldValue
	}

	return idx + 2, false, nil
}

func fieldLineFromString(str string) (fieldName string, fieldValue string, err error) {
	trimmed := strings.TrimSpace(str)

	fieldName, fieldValue, found := strings.Cut(trimmed, ":")

	if len(fieldName) != len(strings.TrimSpace(fieldName)) {
		return "", "", fmt.Errorf("invalid field line format: space between field name and semicolon.")
	}
	if !found {
		return "", "", fmt.Errorf("invalid field line format")
	}

	fieldValue = strings.TrimSpace(fieldValue)

	fieldName = strings.ToLower(fieldName)

	err = validFieldName(fieldName)
	if err != nil {
		return "", "", err
	}

	return fieldName, fieldValue, nil
}

func validFieldName(fieldName string) error {
	if len(fieldName) == 0 {
		return fmt.Errorf("empty field name")
	}
	// Could do a LUT to make O(1) lookup time.
	specialChars := "!#$%&'*+-.^_`|~"
	for _, c := range fieldName {
		isAlphaNumeric := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
		isSpecial := strings.ContainsRune(specialChars, c)

		if !isAlphaNumeric && !isSpecial {
			return fmt.Errorf("invalid characters in field name %c", c)
		}
	}
	return nil
}
