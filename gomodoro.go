package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"time"
)

func main() {
	systray.Run(onReady, nil)
}

func onReady() {
	systray.SetTitle("Gomodoro");
	systray.SetTooltip("Gomodoro");
	mQuitOrig := systray.AddMenuItem("Quit", "Bye bye")
	mStartStop := systray.AddMenuItem("Start", "Start");
	go func() {
		<- mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Quit");
	}()

	tick := make(chan int)
	state := make(chan int)
	go func() {
		j := 1
		for {
			time.Sleep(1 * time.Second)
			tick <- j
			j += 1
		}
	}()

	go func() {
		currentState := 0
		for {
			<- mStartStop.ClickedCh

			switch currentState {
			case 0:
				currentState = 1
				systray.SetTitle("Working...")
			case 1:
				currentState = 0
				systray.SetTitle("Resting")
			}

			state <- currentState
		}
	}()

	go func() {
		for {
			select {
			case <- tick:
				fmt.Println("tick")
			case <- state:
				fmt.Println("state switch")
			}
		}
	}()
}

