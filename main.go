package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/0xAX/notificator"
	tm "github.com/buger/goterm"
	"github.com/y-yagi/configure"
)

type config struct {
	DataBase string `toml:"database"`
}

var (
	cfg         config
	finishSound string
)

func formatMinutes(timeLeft time.Duration) string {
	minutes := int(timeLeft.Minutes())
	seconds := int(timeLeft.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func countDown(outStream io.Writer, target time.Time) {
	tm.Clear()
	for range time.Tick(100 * time.Millisecond) {
		tm.MoveCursor(1, 1)
		timeLeft := -time.Since(target)
		if timeLeft < 0 {
			tm.Print("Countdown: ", formatMinutes(0), "   \r")
			return
		}
		tm.Print("Countdown: ", formatMinutes(timeLeft), "   \r")
		tm.Flush()
	}
}

func task(outStream io.Writer, notify *notificator.Notificator) error {
	start := time.Now()
	finish := start.Add(taskDuration)

	countDown(outStream, finish)

	if notify != nil {
		notify.Push("Tomato", "Pomodoro finished!", "", notificator.UR_CRITICAL)
	}

	_ = exec.Command("mpg123", finishSound).Start()

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

func run(args []string, outStream, errStream io.Writer) (exitCode int) {
	var show string
	var config bool
	var err error

	flags := flag.NewFlagSet("tomato", flag.ExitOnError)
	flags.SetOutput(errStream)
	flags.StringVar(&show, "s", "", "Show your tomatoes. You can specify range, 'today', 'week', 'month' or 'all'.")
	flags.BoolVar(&config, "c", false, "Edit config.")
	flags.Parse(args[1:])

	notify := notificator.New(notificator.Options{
		AppName: "Tomato",
	})
	exitCode = 0

	if config {
		editor := os.Getenv("EDITOR")
		if len(editor) == 0 {
			editor = "vim"
		}

		if err := configure.Edit("tomato", editor); err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			exitCode = 1
			return
		}

		return
	}

	if len(show) != 0 {
		if !contains([]string{"today", "week", "month", "all"}, show) {
			fmt.Printf("'%s' is invalid argument. Please specify 'today', 'week', 'month' or 'all'.\n", show)
			exitCode = 1
			return
		}

		if show == "today" {
			err = showTodayTomatoes(outStream)
		} else {
			err = showTomatoes(outStream, show)
		}

		if err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			exitCode = 1
			return
		}
		return
	}

	for i := 1; ; i++ {
		err = task(outStream, notify)
		if err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			exitCode = 1
			return
		}

		if i%4 == 0 {
			rest(outStream, notify, longRestDuration)
		} else {
			rest(outStream, notify, restDuration)
		}
	}
}

func init() {
	err := configure.Load("tomato", &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = initDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	finishSound = filepath.Join(configure.ConfigDir("tomato"), "ringing.mp3")
	if !isExist(finishSound) {
		err := ioutil.WriteFile(finishSound, Assets.Files["/ringing.mp3"].Data, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

// go:generate go-assets-builder -s="/data" -o bindata.go data

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
