package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// const defaultDuration = 1 * time.Minute
const defaultDuration = 5 * time.Second

func formatMinutes(timeLeft time.Duration) string {
	minutes := int(timeLeft.Minutes())
	seconds := int(timeLeft.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func countDown(target time.Time) {
	for range time.Tick(100 * time.Millisecond) {
		timeLeft := -time.Since(target)
		if timeLeft < 0 {
			fmt.Print("Countdown: ", formatMinutes(0), "   \r")
			return
		}
		fmt.Fprint(os.Stdout, "Countdown: ", formatMinutes(timeLeft), "   \r")
		os.Stdout.Sync()
	}
}

func main() {
	var tag string
	start := time.Now()

	finish := start.Add(defaultDuration)

	fmt.Printf("Start timer.\n\n")

	countDown(finish)

	_ = exec.Command("mpg123", "data/ringing.mp3").Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nTag: ")

	if !scanner.Scan() {
		// Finish without tag
		os.Exit(0)
	}

	if scanner.Err() != nil {
		fmt.Printf("Error: %v\n", scanner.Err())
		os.Exit(1)
	}

	tag = scanner.Text()
	fmt.Printf("Tag: %s\n", tag)
}
