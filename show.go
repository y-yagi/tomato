package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
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

	detectedRange := detectRange(showRange)

	if detectedRange == "today" {
		return showTodayTomatoes(outStream)
	} else if detectedRange == "all" {
		start = time.Date(2000, 01, 01, 00, 00, 00, 0, time.Now().Location())
	} else if detectedRange == "week" {
		start = now.BeginningOfWeek()
	} else if detectedRange == "month" {
		start = now.BeginningOfMonth()
	} else {
		msg := fmt.Sprintf("'%s' is invalid argument. Please specify 'today', 'week', 'month' or 'all'.", detectedRange)
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

func detectRange(showRange string) string {
	validRanges := []string{"today", "week", "month", "all"}
	detectedRange := showRange
	for _, v := range validRanges {
		if strings.HasPrefix(v, showRange) {
			detectedRange = v
			break
		}
	}

	return detectedRange
}
