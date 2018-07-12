package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")

		input, _ := consoleReader.ReadString('\n')

		input = strings.ToLower(input)

		if strings.HasPrefix(input, "exit") {
			os.Exit(0)
		}
	}
}
