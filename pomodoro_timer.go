package tomato

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"github.com/chzyer/readline"
	"github.com/jinzhu/now"
	"github.com/olekukonko/tablewriter"
	"github.com/y-yagi/goext/osext"
	"github.com/y-yagi/goext/strext"
)

// PomodoroTimer is a timer module.
type PomodoroTimer struct {
	out         io.Writer
	notify      *notificator.Notificator
	repo        *Repository
	sound       string
	historyFile string
}

// NewPomodoroTimer creates a new timer.
func NewPomodoroTimer(out io.Writer, notify *notificator.Notificator, repo *Repository, sound string, historyFile string) *PomodoroTimer {
	timer := &PomodoroTimer{out: out, notify: notify, repo: repo, sound: sound, historyFile: historyFile}
	timer.init()
	return timer
}

// Run pomodoro timer.
func (timer *PomodoroTimer) Run() error {
	var tag string
	var err error
	done := make(chan bool)

	fmt.Fprint(timer.out, "Start task.\n")

	timer.countDown(taskDuration)

	if timer.notify != nil {
		timer.notify.Push("Tomato", "Pomodoro finished!", "", notificator.UR_CRITICAL)
	}

	if len(timer.sound) > 0 {
		_ = exec.Command("mpg123", timer.sound).Start()
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "Tags: ",
		InterruptPrompt: "^C",
		Stdout:          timer.out,
		HistoryFile:     timer.historyFile,
	})

	if err != nil {
		return err
	}
	defer l.Close()

	go func() {
		for {
			tag, err = l.Readline()
			if err == readline.ErrInterrupt {
				done <- true
				return
			}

			if !strext.IsBlank(tag) {
				done <- true
				return
			}

			fmt.Fprint(timer.out, "Please input non empty value\nTag: ")
		}
	}()

	for {
		select {
		case <-done:
			if len(tag) != 0 {
				timer.repo.createTomato(tag)
			}
			return err
		case <-time.After(60 * time.Second):
			if timer.notify != nil {
				timer.notify.Push("Tomato", "Please input tag", "", notificator.UR_CRITICAL)
			}
		}
	}
}

// ShortRest take a short rest.
func (timer *PomodoroTimer) ShortRest() {
	timer.rest(restDuration)
}

// LongRest take a long rest.
func (timer *PomodoroTimer) LongRest() {
	timer.rest(longRestDuration)
}

// IsStarted check already started or not.
func (timer *PomodoroTimer) IsStarted() bool {
	output, _ := exec.Command("ps", "--no-headers", "-C", "tomato").Output()

	return strings.Count(string(output), "tomato") > 1
}

func (timer *PomodoroTimer) rest(duration time.Duration) {
	done := make(chan bool)
	fmt.Fprintf(timer.out, "\nStart rest.\n")

	timer.countDown(duration)

	if timer.notify != nil {
		timer.notify.Push("Tomato", "Break is over!", "", notificator.UR_CRITICAL)
	}

	if len(timer.sound) > 0 {
		_ = exec.Command("mpg123", timer.sound).Start()
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(timer.out, "\nPlease press the Enter key for start next tomato.\n")

	go func() {
		scanner.Scan()
		done <- true
	}()

	for {
		select {
		case <-done:
			return
		case <-time.After(60 * time.Second):
			if timer.notify != nil {
				timer.notify.Push("Tomato", "Please press the Enter key for start next tomato.", "", notificator.UR_CRITICAL)
			}
		}
	}
}

// Show past pomodoro.
func (timer *PomodoroTimer) Show(showRange string) error {
	var start time.Time
	var end time.Time

	nothingMsg := "Tomato is nothing (=･x･=)\n"
	detectedRange := timer.detectRange(showRange)

	switch detectedRange {
	case "today":
		return timer.showToday()
	case "all":
		start = time.Date(2000, 01, 01, 00, 00, 00, 0, time.Now().Location())
	case "week":
		start = now.BeginningOfWeek()
	case "month":
		start = now.BeginningOfMonth()
	default:
		msg := fmt.Sprintf("'%s' is invalid argument. Please specify 'today', 'week', 'month' or 'all'.", detectedRange)
		return errors.New(msg)
	}

	end = time.Now()

	tagSummaries, err := timer.repo.selectTagSummary(start, end)
	if err != nil {
		return err
	}

	if len(tagSummaries) == 0 {
		fmt.Fprintf(timer.out, nothingMsg)
		return nil
	}

	table := tablewriter.NewWriter(timer.out)
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

func (timer *PomodoroTimer) showToday() error {
	nothingMsg := "Tomato is nothing (=･x･=)\n"

	tomatoes, err := timer.repo.selectTomatos(now.BeginningOfDay(), now.EndOfDay())
	if err != nil {
		return err
	}

	if len(tomatoes) == 0 {
		fmt.Fprintf(timer.out, nothingMsg)
		return nil
	}

	table := tablewriter.NewWriter(timer.out)
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

func (timer *PomodoroTimer) detectRange(showRange string) string {
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

func (timer *PomodoroTimer) formatMinutes(timeLeft time.Duration) string {
	minutes := int(timeLeft.Minutes())
	seconds := int(timeLeft.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func (timer *PomodoroTimer) countDown(duration time.Duration) {
	start := time.Now()
	finish := start.Add(duration)

	for {
		timeLeft := -time.Since(finish)
		if timeLeft < 0 {
			fmt.Fprint(timer.out, "Countdown: ", timer.formatMinutes(0), "   \r")
			return
		}

		fmt.Fprint(timer.out, "Countdown: ", timer.formatMinutes(timeLeft), "   \r")
		time.Sleep(1 * time.Second)
	}
}

func (timer *PomodoroTimer) init() {
	if !osext.IsExist(timer.historyFile) {
		tags, _ := timer.repo.selectTags()
		ioutil.WriteFile(timer.historyFile, []byte(strings.Join(tags, "\n")), 0644)
	}
}
