package profiler

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log           = logrus.New()
	profileLevels = 0
)

func init() {
	// Log as JSON instead of the default ASCII formatter
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Output to a file instead of stderr
	file, err := os.OpenFile("profiler.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr", err)
	}
}

// Profiler is used to track a series of checkpoints to time how long code takes to execute.
type Profiler struct {
	name        string       // my name, so I can be identified later
	checkpoints []Checkpoint // every tick is stored here
	profile     bool         // whether I am active or not
	level       int          // my activation level
}

// Checkpoint holds the profiler data
type Checkpoint struct {
	name string
	time time.Time
}

func activationCheck(level int) bool {
	if level <= profileLevels {
		return true
	}
	return false
}

// New initialises a new profiler for data capture
func New(title string, level int) (cp *Profiler) {
	// level is a filter set by the programmer against the Flags given by command line argument.
	// the filters are so the programmers can decide at runtime which profilers to activate
	cp = &Profiler{
		name:        title,
		profile:     activationCheck(level),
		checkpoints: []Checkpoint{{name: "start", time: time.Now()}},
		level:       level,
	}
	return
}

// Tick records the current time onto the Profiler
func (p *Profiler) Tick(title string) {
	if p.profile {
		p.checkpoints = append(p.checkpoints, Checkpoint{title, time.Now()})
	}
}

// Finish prints out the checkpoint to log and stdout
func (p *Profiler) Finish() {
	if p.profile {
		totalTime := p.checkpoints[len(p.checkpoints)-1].time.Sub(p.checkpoints[0].time)
		msg := ""
		for i, checkpoint := range p.checkpoints {
			if i > 0 {
				msg += fmt.Sprintf("\t%s:%s", checkpoint.name, (checkpoint.time.Sub(p.checkpoints[i-1].time)))
			}
		}
		log.WithFields(logrus.Fields{"function": "Profiler", "Level": p.level, "Target": p.name, "Total": totalTime.String()}).Info(msg)
		if p.profile {
			msg += "\n"
		}
		fmt.Println("Total: ", totalTime.String(), "\t(", p.level, ") ", p.name, "\t", msg)
	}
}

// SetProfileLevel sets this packages recording level for recording data.
// all profilers at or below this level will record and report on data.
func SetProfileLevel(levels int) {
	msg := fmt.Sprintf(" - Enabling Profiling Level %d", levels)
	profileLevels = levels
	fmt.Println(msg)
	log.Info(msg)
}
