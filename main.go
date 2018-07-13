// https://cstack.github.io/db_tutorial/parts/part2.html
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Unrecognized       = "unrecognized"
	Insert             = "insert"
	Success            = "success"
	Select             = "select"
	Error              = "error"
	ColumnUsernameSize = 32
	ColumnEmailSize    = 255
	TableFull          = "table full"
	SyntaxError        = "syntax error"
	RowSize            = 32 + ColumnUsernameSize + ColumnEmailSize
	PageSize           = 4096
	TableMaxPages      = 100
	RowsPerPage        = PageSize / RowSize
	TableMaxRows       = RowsPerPage * TableMaxPages
)

type Table struct {
	Pages   [TableMaxPages]*Row
	NumRows int
}

type Statement struct {
	Type        string
	RowToInsert Row
}

type Row struct {
	Id       uint32
	Username [ColumnUsernameSize]byte
	Email    [ColumnEmailSize]byte
}

func serializeRow(source *Row, destination *Row) {
	copy([]*Row{destination}, []*Row{source})
}

func deserializeRow(source *Row, destination *Row) {
	copy([]*Row{destination}, []*Row{source})
}

func rowSlot(table *Table, rowNum int) *Row {
	pageNum := rowNum / RowsPerPage
	page := table.Pages[pageNum]
	return page
}

func printRow(row *Row) {
	fmt.Printf("(%d, %s, %s)\n", row.Id, row.Username, row.Email)
}

func doMetaCommand(input string) string {
	var output string
	if input == ".exit" {
		os.Exit(0)
	} else {
		output = Unrecognized
	}
	return output
}

func executeInsert(statement *Statement, table *Table) string {
	if table.NumRows > TableMaxRows {
		return TableFull
	}
	rowToInsert := statement.RowToInsert
	page := rowSlot(table, table.NumRows)
	serializeRow(&rowToInsert, page)
	table.NumRows++

	return Success
}

func executeSelect(statement *Statement, table *Table) string {
	var row Row
	for i := 0; i < table.NumRows; i++ {
		page := rowSlot(table, i)
		deserializeRow(page, &row)
		printRow(&row)
	}
	return Success
}

func executeStatement(statement *Statement, table *Table) string {
	switch statement.Type {
	case Insert:
		return executeInsert(statement, table)
	case Select:
		return executeSelect(statement, table)
	}
	return SyntaxError
}

func prepareStatement(input string, statement *Statement) string {
	if input[:6] == Insert {
		statement.Type = Insert
		_, err := fmt.Sscanf(input, "insert %d %s %s", statement.RowToInsert.Id, statement.RowToInsert.Username, statement.RowToInsert.Email)
		if err != nil {
			return Error
		}
		return Success
	}

	if input[:6] == Select {
		statement.Type = Select
		return Success
	}

	return Unrecognized
}

func main() {
	var table Table
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
		case SyntaxError:
			fmt.Println("Syntax error. Could not parse statement.")
			continue
		case Unrecognized:
			fmt.Printf("Unrecognized keyword at start of '%s'.\n", input)
			continue
		}

		switch executeStatement(&statement, &table) {
		case Success:
			fmt.Println("Executed.")
			break
		case TableFull:
			fmt.Println("Error: Table full.")
			break
		}
	}

}
