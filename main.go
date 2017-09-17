package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/0xAX/notificator"
)

func formatMinutes(timeLeft time.Duration) string {
	minutes := int(timeLeft.Minutes())
	seconds := int(timeLeft.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func countDown(outStream io.Writer, target time.Time) {
	for range time.Tick(100 * time.Millisecond) {
		timeLeft := -time.Since(target)
		if timeLeft < 0 {
			fmt.Fprint(outStream, "Countdown: ", formatMinutes(0), "   \r")
			return
		}
		fmt.Fprint(outStream, "Countdown: ", formatMinutes(timeLeft), "   \r")
		os.Stdout.Sync()
	}
}

func task(outStream io.Writer, notify *notificator.Notificator) error {
	start := time.Now()
	finish := start.Add(taskDuration)
	fmt.Fprint(outStream, "Start task.\n")

	countDown(outStream, finish)

	if notify != nil {
		notify.Push("Tomato", "Pomodoro finished!", "", notificator.UR_CRITICAL)
	}
	_ = exec.Command("mpg123", "data/ringing.mp3").Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(outStream, "\nTag: ")

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

func rest(outStream io.Writer, notify *notificator.Notificator, duration time.Duration) {
	start := time.Now()
	finish := start.Add(duration)
	fmt.Fprintf(outStream, "\nStart rest.\n")

	countDown(outStream, finish)

	if notify != nil {
		notify.Push("Tomato", "Break is over!", "", notificator.UR_CRITICAL)
	}
	_ = exec.Command("mpg123", "data/ringing.mp3").Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(outStream, "\nPlease press the Enter key for start next tomato.\n")
	scanner.Scan()
}

func run(args []string, outStream, errStream io.Writer) int {
	var show string
	flags := flag.NewFlagSet("tomato", flag.ExitOnError)
	flags.SetOutput(errStream)
	flags.StringVar(&show, "s", "", "Show your tomatoes. You can specify range, 'today', 'week', 'month' or 'all'.")
	flags.Parse(args[1:])

	err := initDB()
	if err != nil {
		fmt.Fprintf(outStream, "Error: %v\n", err)
		return 1
	}

	notify := notificator.New(notificator.Options{
		AppName: "Tomato",
	})

	if len(show) != 0 {
		if !contains([]string{"today", "week", "month", "all"}, show) {
			fmt.Printf("'%s' is invalid argument. Please specify 'today', 'week', 'month' or 'all'.\n", show)
			return 1
		}

		err = showTomatoes(outStream, show)
		if err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			return 1
		}
		return 0
	}

	for i := 1; ; i++ {
		err = task(outStream, notify)
		if err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			return 1
		}

		if i%4 == 0 {
			rest(outStream, notify, longRestDuration)
		} else {
			rest(outStream, notify, restDuration)
		}
	}
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
