package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	defer file.Close()

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
