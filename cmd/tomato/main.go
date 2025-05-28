package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"

	"github.com/0xAX/notificator"
	_ "github.com/mattn/go-sqlite3"
	"github.com/y-yagi/configure"
	"github.com/y-yagi/goext/osext"
	"github.com/y-yagi/tomato"
)

type config struct {
	DataBase string `toml:"database"`
}

var (
	cfg         config
	finishSound string

	//go:embed ringing.mp3
	soundData []byte
)

func init() {
	err := configure.Load("tomato", &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	finishSound = filepath.Join(configure.ConfigDir("tomato"), "ringing.mp3")
	if !osext.IsExist(finishSound) {
		err := os.WriteFile(finishSound, soundData, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	if len(cfg.DataBase) == 0 {
		cfg.DataBase = filepath.Join(configure.ConfigDir("tomato"), "tomato.db")
		configure.Save("tomato", cfg)
	}
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func run(args []string, outStream, errStream io.Writer) (exitCode int) {
	var show string
	var config bool
	var console bool
	var err error
	var noSound bool

	exitCode = 0

	flags := flag.NewFlagSet("tomato", flag.ExitOnError)
	flags.SetOutput(errStream)
	flags.StringVar(&show, "s", "", "Show your tomatoes. You can specify `range`, 'today', 'week', 'month' or 'all'.")
	flags.BoolVar(&config, "c", false, "Edit config.")
	flags.BoolVar(&console, "db", false, "Start a console for the database.")
	flags.BoolVar(&noSound, "no-sound", false, "Do not play sound when finished timer.")
	flags.Parse(args[1:])

	notify := notificator.New(notificator.Options{
		AppName: "Tomato",
	})

	if config {
		if err := cmdConfig(); err != nil {
			fmt.Fprintf(errStream, "Error: %v\n", err)
			exitCode = 1
			return
		}

		return
	}

	if console {
		if err = cmdConsole(cfg.DataBase); err != nil {
			fmt.Fprintf(errStream, "Error: %v\n", err)
			exitCode = 1
			return
		}

		return
	}

	repo := tomato.NewRepository(cfg.DataBase)
	err = repo.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if noSound {
		finishSound = ""
	}
	historyFile := filepath.Join(configure.ConfigDir("tomato"), "readline.tmp")
	timer := tomato.NewPomodoroTimer(outStream, notify, repo, finishSound, historyFile)

	if len(show) != 0 {
		err = timer.Show(show)
		if err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			exitCode = 1
			return
		}
		return
	}

	if timer.IsStarted() {
		fmt.Fprintf(outStream, "Tomato has already started.\n")
		exitCode = 1
		return
	}

	for i := 1; ; i++ {
		err = timer.Run()
		if err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			exitCode = 1
			return
		}

		if i%4 == 0 {
			timer.LongRest()
		} else {
			timer.ShortRest()
		}
	}
}

func cmdConsole(database string) error {
	cmd := exec.Command("sqlite3", database)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func cmdConfig() error {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vim"
	}

	return configure.Edit("tomato", editor)
}
