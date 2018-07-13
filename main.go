// https://cstack.github.io/db_tutorial/parts/part2.html
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Unrecognized = "unrecognized"
	Insert       = "insert"
	Success      = "success"
	Select       = "select"
)

type Statement struct{ Type string }

func doMetaCommand(input string) string {
	var output string
	if input == ".exit" {
		os.Exit(0)
	} else {
		output = Unrecognized
	}
	return output
}

func executeStatement(statement *Statement) {
	switch statement.Type {
	case Insert:
		fmt.Println("This is where we would do an insert.")
		break
	case Select:
		fmt.Println("This is where we would do a select.")
	}
}

func prepareStatement(input string, statement *Statement) string {
	if input[:6] == Insert {
		statement.Type = Insert
		return Success
	}

	if input[:6] == Select {
		statement.Type = Select
		return Success
	}

	return Unrecognized
}

func main() {
	for {
		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, _ := consoleReader.ReadString('\n')
		input = strings.ToLower(input)
		input = strings.TrimSuffix(input, "\n")

		if input[:1] == "." {
			switch doMetaCommand(input) {
			case Success:
				continue
			case Unrecognized:
				fmt.Printf("Unrecognized command '%s'.\n", input)
				continue
			}
		}

		var statement Statement
		switch prepareStatement(input, &statement) {
		case Success:
			break
		case Unrecognized:
			fmt.Printf("Unrecognized keyword at start of '%s'.\n", input)
			continue
		}

		executeStatement(&statement)
		fmt.Println("Executed.")
	}

}
