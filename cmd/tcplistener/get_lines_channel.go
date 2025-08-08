package main

import (
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		line := ""

		for {
			bs := make([]byte, 8)
			n, err := f.Read(bs)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}

			parts := strings.Split(string(bs[:n]), "\n")

			if len(parts) > 1 {
				line += parts[0]
				lines <- line
				line = parts[1]
			} else {
				line += parts[0]
			}
		}

		if len(line) > 0 {
			lines <- line
		}
	}()
	return lines
}
