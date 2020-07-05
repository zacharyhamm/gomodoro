package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"time"
	"math"
)

type GomodoroState int
const (
	Initial GomodoroState = 0
	Resting GomodoroState = 1
	Working GomodoroState = 2
	Toggle GomodoroState = 3
)

const WorkMinutes int64 = 23
const RestMinutes int64 = 7

type SystrayUpdate struct {
	state GomodoroState
	stateChangeTime time.Time
}

func main() {
	systray.Run(onReady, nil)
}

func timeSince(state *SystrayUpdate) (int64, int64) {
	elapsed := time.Since(state.stateChangeTime)

	elapsedSeconds := int64(math.Round(elapsed.Seconds()))
	elapsedMinutes := int64(elapsedSeconds / 60)

	return elapsedMinutes, elapsedSeconds - elapsedMinutes * 60
}

func updateSystray(state *SystrayUpdate) (int64, int64) {
	if state.stateChangeTime.IsZero() || state.state == Initial {
		systray.SetTitle("🍅");
		return 0, 0
	}

	minutes, seconds := timeSince(state)
	var stateString string
	switch state.state {
	case Resting:
		stateString = "☕"
	case Working:
		stateString = "🔨"
	}

	systray.SetTitle(fmt.Sprintf("%s%dm%ds", stateString, minutes, seconds))

	return minutes, seconds
}

func onReady() {
	systray.SetTitle("🍅")
	systray.SetTooltip("Gomodoro");
	mStartStop := systray.AddMenuItem("Start/Stop", "Start/Stop")
	mQuitOrig := systray.AddMenuItem("Quit", "Bye bye")

	globalState := &SystrayUpdate{Initial, time.Unix(0, 0)}

	go func() {
		<- mQuitOrig.ClickedCh
		systray.Quit()
		fmt.Println("Quit");
	}()

	tick := make(chan int)
	go func() {
		j := 1
		for {
			time.Sleep(1025 * time.Millisecond)
			tick <- j
			j += 1
		}
	}()

	stateChange := make(chan GomodoroState)

	go func() {
		for {
			<- mStartStop.ClickedCh
			stateChange <- Toggle
		}
	}()

	go func() {
		for {
			state := <- stateChange
			globalState.stateChangeTime = time.Now()

			if state == Toggle {
				if globalState.state == Initial || globalState.state == Resting {
					globalState.state = Working
				} else {
					globalState.state = Resting
				}
			} else {
				globalState.state = state
			}

			updateSystray(globalState)
		}
	}()

	go func() {
		for {
			<- tick
			minutes, _ := updateSystray(globalState)

			if globalState.state == Resting && minutes >= RestMinutes {
				stateChange <- Initial
				continue
			}

			if globalState.state == Working && minutes >= WorkMinutes {
				stateChange <- Resting 
			}
		}
	}()
}

