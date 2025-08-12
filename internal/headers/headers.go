package headers

import (
	"bytes"
	"errors"
	"strings"

	"github.com/Waterbootdev/http-from-tcp/internal/commen"
)

type Headers map[string]string

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

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

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

	h[key] = value

	return numberBytesReaded, false, err
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
