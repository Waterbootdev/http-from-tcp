package server

import "strings"

func testSwapPrefix(target string, endPointPrefix string, targetPrefix string) *string {

	if strings.HasPrefix(target, endPointPrefix) {
		return swapPrefix(target, endPointPrefix, targetPrefix)
	}

	return nil
}

func swapPrefix(s string, oldPrefix, newPrefix string) *string {
	r := strings.Replace(s, oldPrefix, newPrefix, 1)
	return &r
}
