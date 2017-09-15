package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jinzhu/now"
	"github.com/olekukonko/tablewriter"
)

var msg = "Tomato is nothing (=･x･=)\n"

func showTodayTomatoes() error {
	tomatoes, err := selectTomatos(now.BeginningOfDay(), now.EndOfDay())
	if err != nil {
		return err
	}

	if len(tomatoes) == 0 {
		fmt.Fprintf(os.Stdout, msg)
		return nil
	}

	w := os.Stdout
	table := tablewriter.NewWriter(w)
	var values = []string{}

	for i, tomato := range tomatoes {
		values = append(values, strconv.Itoa(i+1))
		values = append(values, tomato.CreatedAt.Format("2006-01-02 15:04"))
		values = append(values, tomato.Tag)
		table.Append(values)
		values = nil
	}

	table.Render()
	return nil

}

func showTomatoes(showRange string) error {
	if showRange == "today" {
		return showTodayTomatoes()
	}

	return nil
}
