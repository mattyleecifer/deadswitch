package main

import (
	"fmt"
	"time"
)

func (ds switchdefinitions) resetTimers() {
	ds.mainTimerCh <- time.Now()
}

func (ds switchdefinitions) mainTimer() {
	fulltimerlength := float32(ds.days) * 24

	var timerlength float32
	if ds.timeleft != 0.0 {
		timerlength = ds.timeleft * 24
	} else {
		ds.timeleft = float32(ds.days)
		timerlength = fulltimerlength
	}
	t := time.NewTicker(time.Hour * time.Duration(timerlength)) // Restarting may have different first timer

	// t := time.NewTicker(time.Minute) // for testing
	defer t.Stop()
	for {
		select {
		case <-ds.mainTimerCh:
			// send email
			t = time.NewTicker(time.Hour * time.Duration(fulltimerlength))
			// t = time.NewTicker(time.Minute) // for testing

			fmt.Println("Timer reset at", time.Now())
			// reset_time := <-ds.timerCh
			// fmt.Println("Timer reset at", reset_time)
		case <-t.C:
			fmt.Println("Timer expired")
			// send email
			ds.sendemail("Deadswitch", ds.message)
			fmt.Println("Deadswitch activated at ", time.Now())
			panic("Deadswitch activated")
		}
	}
}

func (ds switchdefinitions) writeToFileEverySixHours() {
	ticker := time.NewTicker(time.Hour * 6)
	for {
		select {
		case <-ticker.C:
			ds.timeleft -= 0.25
			ds.writeSwitchDefinitionsToFile()
			ds.writeToFileEverySixHours()
		}
	}
}
