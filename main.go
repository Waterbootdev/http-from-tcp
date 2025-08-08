package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func readLines() {
	file, err := os.Open("message.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	line := ""

	for {
		bs := make([]byte, 8)
		_, err = file.Read(bs)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		parts := strings.Split(string(bs), "\n")

		if len(parts) > 1 {
			line += parts[0]
			fmt.Printf("read: %s\n", line)
			line = parts[1]
		} else {
			line += parts[0]
		}
	}

	if len(line) > 0 {
		fmt.Printf("line: %s\n", line)
		line = ""
	}
}

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for {
		bs := make([]byte, 8)
		_, err = file.Read(bs)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Printf("read: %s\n", string(bs))
	}

}
