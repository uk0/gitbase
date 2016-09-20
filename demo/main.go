package main

import (
	"github.com/mvader/gitql/sql"
	"github.com/mvader/gitql/git"
	"github.com/olekukonko/tablewriter"

	"fmt"
	"io"
	"os"
	"github.com/mvader/gitql/sql/plan"
	"github.com/mvader/gitql/sql/expression"
)

func main() {
	db := git.NewDatabase("https://github.com/mvader/gitql.git")
	tables := db.Relations()

	fmt.Println("SELECT * FROM commits;\n")
	var q sql.Node = tables["commits"]
	printQuery(q)


	fmt.Println("SELECT * FROM commits;\n")
	q = plan.NewProject(
		[]sql.Expression{
			expression.NewGetField(2, sql.String, "author_name"),
			expression.NewGetField(7, sql.String, "message"),
		},
		plan.NewFilter(
			expression.NewEquals(
				expression.NewGetField(2, sql.String, "author_name"),
				expression.NewLiteral("santi@mola.io", sql.String),
			),
			tables["commits"],
		)
	)
	printQuery(q)
}

func printQuery(q sql.Node) {
	w := tablewriter.NewWriter(os.Stdout)
	headers := []string{}
	s := q.Schema()
	for _, f := range s {
		headers = append(headers, f.Name)
	}
	w.SetHeader(headers)
	iter, err := q.RowIter()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for {
		row, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		rowStrings := []string{}
		for _, v := range row.Fields() {
			rowStrings = append(rowStrings, fmt.Sprintf("%v", v))
		}
		w.Append(rowStrings)
	}
	w.Render()
}
