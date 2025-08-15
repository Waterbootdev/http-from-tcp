package headers

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/Waterbootdev/http-from-tcp/internal/commen"
)

type Headers map[string]string

func (h Headers) IsContentLengthNot(contentLengthKey string, contentLength int) bool {
	if value, ok := h[contentLengthKey]; ok {
		return strconv.Itoa(contentLength) != value
	} else {
		return contentLength != 0
	}
}

func NewHeaders() Headers {
	return make(Headers)
}

const UPPERCASE_LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LOWERCASE_LETTERS = "abcdefghijklmnopqrstuvwxyz"
const DIGITS = "0123456789"
const SPECIAL_CHRACTERS = "!#$%&'*+-.^_`|~"
const VALID_RUNES = UPPERCASE_LETTERS + LOWERCASE_LETTERS + DIGITS + SPECIAL_CHRACTERS

func isNotValidToken(toValidate string) bool {

	if len(toValidate) == 0 {
		return true
	}

	for _, char := range toValidate {
		if !strings.ContainsRune(VALID_RUNES, char) {
			return true
		}
	}

	return false
}

func (h Headers) parseSingel(data []byte) (n int, done bool, err error) {

	crlfIndex := bytes.Index(data, []byte(commen.CRLF))

	if crlfIndex == -1 {
		return 0, false, nil
	}

	numberBytesReaded := crlfIndex + commen.LENGTH_CRLF

	if crlfIndex == 0 {
		return numberBytesReaded, true, nil
	}

	key, value, err := parseFieldLine(string(data[:crlfIndex]))

	if err != nil {
		return 0, false, err
	}

	if existingValue, ok := h[key]; ok {
		value = existingValue + ", " + value
	}

	h[key] = value

	return numberBytesReaded, false, err
}
func (h Headers) Parse(data []byte) (numberBytesParsed int, done bool, err error) {

	for !done {

		var lastNumberBytesParsed int

		lastNumberBytesParsed, done, err = h.parseSingel(data[numberBytesParsed:])

		if err != nil {

			for key := range h {
				delete(h, key)
			}

			return 0, done, err
		}

		if lastNumberBytesParsed == 0 {
			return numberBytesParsed, done, err
		}

		numberBytesParsed += lastNumberBytesParsed
	}

	return numberBytesParsed, done, err
}

func splitFieldLine(s string) (string, string, error) {
	s = strings.TrimSpace(s)
	parts := strings.SplitN(s, ":", 2)

	if len(parts) != 2 {
		return "", "", errors.New("cant	split field line")
	}

	return parts[0], strings.TrimSpace(parts[1]), nil
}

func parseFieldLine(s string) (string, string, error) {

	token, value, err := splitFieldLine(s)

	if err != nil {
		return "", "", err
	}

	if isNotValidToken(token) {
		return "", "", errors.New("key is not a valid token")
	}

	return strings.ToLower(token), value, nil
}

func (h Headers) HeadersString() string {

	var buffer bytes.Buffer
	buffer.WriteString("Headers:\r\n")

	for key, value := range h {
		buffer.WriteString("- ")
		buffer.WriteString(key)
		buffer.WriteString(": ")
		buffer.WriteString(value)
		buffer.WriteString("\r\n")
	}

	return buffer.String()
}
