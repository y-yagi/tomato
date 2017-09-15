package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// const duration = 1 * time.Minute
const taskDuration = 5 * time.Second
const restDuration = 3 * time.Second
const longRestDuration = 5 * time.Second

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

func task() error {
	start := time.Now()
	finish := start.Add(taskDuration)
	fmt.Printf("Start task.\n")

	countDown(finish)

	_ = exec.Command("mpg123", "data/ringing.mp3").Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nTag: ")

	if !scanner.Scan() {
		return nil
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	tag := scanner.Text()
	createTomato(tag)

	return nil
}

func rest(duration time.Duration) {
	start := time.Now()
	finish := start.Add(duration)
	fmt.Printf("\nStart rest.\n")

	countDown(finish)

	_ = exec.Command("mpg123", "data/ringing.mp3").Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nPlease press the Enter key for start next tomato.\n")
	scanner.Scan()
}

func main() {
	const version = "0.1.0"
	var show string

	flags := flag.NewFlagSet("goma", flag.ExitOnError)
	flags.StringVar(&show, "s", "", "Show your tomatoes. You can specify range, 'today', 'week', 'month', 'all'.")
	flags.Parse(os.Args[1:])

	err := initDB()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(show) != 0 {
		err = showTomatoes(show)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	for i := 1; ; i++ {
		err = task()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if i%4 == 0 {
			rest(longRestDuration)
		} else {
			rest(restDuration)
		}
	}
}
