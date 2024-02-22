package main

import (
	"fmt"
	"time"
)

func (ds switchdefinitions) resetTimers() {
	ds.mainTimerCh <- time.Now()
}

func (ds switchdefinitions) mainTimer() {
	hours_in_days := ds.days * 24
	t := time.NewTicker(time.Hour * time.Duration(hours_in_days))
	// t := time.NewTicker(time.Minute) // for testing
	defer t.Stop()
	for {
		select {
		case <-ds.mainTimerCh:
			// send email
			t = time.NewTicker(time.Hour * time.Duration(hours_in_days))
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
