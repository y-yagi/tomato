package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/0xAX/notificator"
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
)

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

	repo := tomato.NewRepository(cfg.DataBase)
	err = repo.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	timer := tomato.NewPomodoroTimer(outStream, notify, repo, finishSound)

	if len(show) != 0 {
		err = timer.Show(show)
		if err != nil {
			fmt.Fprintf(outStream, "Error: %v\n", err)
			exitCode = 1
			return
		}
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

func init() {
	err := configure.Load("tomato", &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	finishSound = filepath.Join(configure.ConfigDir("tomato"), "ringing.mp3")
	if !osext.IsExist(finishSound) {
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