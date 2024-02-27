package main

import (
	"fmt"
	"os"
	"time"
)

func (ds *switchdefinitions) resetTimers() {
	currenttime := time.Now()
	ds.hoursleft = ds.days * 24
	ds.mainTimerCh <- currenttime
	ds.halfTimerCh <- currenttime
	ds.quarterTimerCh <- currenttime
	ds.sendemail("Deadswitch reset", "This is a message to alert you that your deadswitch was just reset.", true)
}

func (ds *switchdefinitions) mainTimer() {
	fulltimerlength := ds.days * 24

	var timerlength int
	if ds.hoursleft != 0 {
		timerlength = ds.hoursleft
	} else {
		ds.hoursleft = ds.days * 24
		timerlength = fulltimerlength
	}
	t := time.NewTicker(time.Hour * time.Duration(timerlength)) // Restarting may have different first timer
	// t := time.NewTicker(time.Minute * time.Duration(timerlength)) // for testing

	defer t.Stop()
	for {
		select {
		case <-ds.mainTimerCh:
			// reset timer
			t = time.NewTicker(time.Hour * time.Duration(fulltimerlength))
			// t = time.NewTicker(time.Minute * time.Duration(fulltimerlength)) // for testing

			fmt.Println("Timer reset at", time.Now())
		case <-t.C:
			fmt.Println("Timer expired")
			// send email
			ds.sendemail("Deadswitch", ds.message, false)
			fmt.Println("Deadswitch activated at ", time.Now())
			os.Exit(0)
		}
	}
}

func (ds *switchdefinitions) halfTimer() {
	fulltimerlength := ds.days * 12

	var timerlength int
	if ds.hoursleft != 0 {
		if ds.hoursleft > fulltimerlength {
			timerlength = ds.hoursleft - fulltimerlength
		} else {
			<-ds.halfTimerCh // blocks until timer is reset
			timerlength = fulltimerlength
			fmt.Println("Half timer reset at", time.Now())
		}
	} else {
		timerlength = fulltimerlength
	}
	t := time.NewTicker(time.Hour * time.Duration(timerlength)) // Restarting may have different first timer
	// t := time.NewTicker(time.Minute * time.Duration(timerlength)) // for testing

	defer t.Stop()
	for {
		select {
		case <-ds.halfTimerCh:
			// reset timer
			t = time.NewTicker(time.Hour * time.Duration(fulltimerlength))
			// t = time.NewTicker(time.Minute * time.Duration(fulltimerlength)) // for testing

			fmt.Println("Half timer reset at", time.Now())
		case <-t.C:
			fmt.Println("Half timer expired")
			// send email
			ds.sendemail("Deadswitch Reminder", "Your deadswitch has reached its halfway mark. Please remember to log in and reset the switch", true)
			fmt.Println("Reminder sent at ", time.Now())
			t.Stop()
		}
	}
}

func (ds *switchdefinitions) quarterTimer() {
	quartertime := ds.days * 6
	fulltimerlength := ds.days * 18 // 3/4 time

	var timerlength int
	if ds.hoursleft != 0 {
		if ds.hoursleft > quartertime {
			timerlength = ds.hoursleft - quartertime
		} else {
			<-ds.quarterTimerCh // blocks until timer is reset
			timerlength = fulltimerlength
			fmt.Println("3/4 timer reset at", time.Now())
		}
	} else {
		timerlength = fulltimerlength
	}
	t := time.NewTicker(time.Hour * time.Duration(timerlength)) // Restarting may have different first timer
	// t := time.NewTicker(time.Minute * time.Duration(timerlength)) // for testing

	defer t.Stop()
	for {
		select {
		case <-ds.quarterTimerCh:
			// reset timer
			t = time.NewTicker(time.Hour * time.Duration(fulltimerlength))
			// t = time.NewTicker(time.Minute * time.Duration(fulltimerlength)) // for testing

			fmt.Println("3/4 timer reset at", time.Now())
		case <-t.C:
			fmt.Println("3/4 timer expired")
			// send email
			ds.sendemail("Deadswitch Reminder", "Your deadswitch has reached its 3/4 mark. Please remember to log in and reset the switch", true)
			fmt.Println("Reminder sent at ", time.Now())
			t.Stop()
		}
	}
}

func (ds *switchdefinitions) writeToFileEveryHour() {
	ticker := time.NewTicker(time.Hour)
	// ticker := time.NewTicker(time.Minute) // for testing
	for {
		<-ticker.C
		ds.hoursleft -= 1
		ds.writeSwitchDefinitionsToFile()
	}
}
