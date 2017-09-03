package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const defaultDuration = 1 * time.Minute

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
	start := time.Now()

	finish := start.Add(defaultDuration)

	fmt.Printf("Start timer.\n\n")

	countDown(finish)

	_ = exec.Command("mpg123", "data/ringing.mp3").Run()
}
