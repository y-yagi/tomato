package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/olekukonko/tablewriter"
)

var nothingMsg = "Tomato is nothing (=･x･=)\n"

func showTodayTomatoes(outStream io.Writer) error {
	tomatoes, err := selectTomatos(now.BeginningOfDay(), now.EndOfDay())
	if err != nil {
		return err
	}

	if len(tomatoes) == 0 {
		fmt.Fprintf(outStream, nothingMsg)
		return nil
	}

	table := tablewriter.NewWriter(outStream)
	table.SetHeader([]string{"id", "Created", "Tag"})
	var values = []string{}

	for i, tomato := range tomatoes {
		values = append(values, strconv.Itoa(i+1))
		values = append(values, tomato.CreatedAt.Format("15:04"))
		values = append(values, tomato.Tag)
		table.Append(values)
		values = nil
	}

	table.Render()

	return nil
}

func showTomatoes(outStream io.Writer, showRange string) error {
	var start time.Time
	var end time.Time

	if showRange == "today" || showRange == "t" {
		return showTodayTomatoes(outStream)
	} else if showRange == "all" || showRange == "a" {
		start = time.Date(2000, 01, 01, 00, 00, 00, 0, time.Now().Location())
	} else if showRange == "week" || showRange == "w" {
		start = now.BeginningOfWeek()
	} else if showRange == "mohth" || showRange == "m" {
		start = now.BeginningOfMonth()
	} else {
		msg := fmt.Sprintf("'%s' is invalid argument. Please specify 'today', 'week', 'month' or 'all'.", showRange)
		return errors.New(msg)
	}

	end = time.Now()

	tagSummaries, err := selectTagSummary(start, end)
	if err != nil {
		return err
	}

	if len(tagSummaries) == 0 {
		fmt.Fprintf(outStream, nothingMsg)
		return nil
	}

	table := tablewriter.NewWriter(outStream)
	table.SetHeader([]string{"Count", "Tag"})
	var values = []string{}

	for _, tagSummary := range tagSummaries {
		values = append(values, strconv.Itoa(tagSummary.Count))
		values = append(values, tagSummary.Tag)
		table.Append(values)
		values = nil
	}

	table.Render()

	return nil
}
