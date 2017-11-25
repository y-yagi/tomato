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
	var tag string
	var err error
	done := make(chan bool)

	start := time.Now()
	finish := start.Add(taskDuration)
	fmt.Fprint(outStream, "Start task.\n")

	countDown(outStream, finish)

	if notify != nil {
		notify.Push("Tomato", "Pomodoro finished!", "", notificator.UR_CRITICAL)
	}

	_ = exec.Command("mpg123", finishSound).Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(outStream, "\nTag: ")

	go func() {
		for {
			scanner.Scan()

			if scanner.Err() != nil {
				err = scanner.Err()
				done <- true
				return
			}

			tag = scanner.Text()

			if !isBlank(tag) {
				done <- true
				return
			}

			fmt.Fprint(outStream, "Please input non empty value\nTag: ")
		}
	}()

	for {
		select {
		case <-done:
			if len(tag) != 0 {
				createTomato(tag)
			}
			return err
		case <-time.After(60 * time.Second):
			if notify != nil {
				notify.Push("Tomato", "Please input tag", "", notificator.UR_CRITICAL)
			}
		}
	}
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
		if !contains([]string{"today", "week", "month", "all", "t", "w", "m", "a"}, show) {
			fmt.Printf("'%s' is invalid argument. Please specify 'today', 'week', 'month' or 'all'.\n", show)
			exitCode = 1
			return
		}

		err = showTomatoes(outStream, show)
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
