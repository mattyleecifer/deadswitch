package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type switchdefinitions struct {
	days        int
	timeleft    float32
	recipients  []string
	message     string
	auth        string // password for deadswitch
	files       []string
	owner       string // to send reminder email
	key         string // gmail key
	port        string
	mainTimerCh chan time.Time
}

type saveDefinitions struct {
	Days       int
	Timeleft   float32 // restart should only read this
	Recipients []string
	Message    string
	Auth       string
	Files      []string
	Owner      string
	Key        string
	Port       string
}

func main() {
	ds, err := readSwitchDefinitionsFromFile() // this checks flags as well
	if err != nil {
		fmt.Println(err)
	}

	ds.writeSwitchDefinitionsToFile() // immediately save data

	go ds.mainTimer()                // start main timer
	go ds.writeToFileEverySixHours() // basically keeps it updated with days left
	// start secondary timers

	// start server
	http.HandleFunc("/", ds.hauth)
	fmt.Println("Running GUI on http://127.0.0.1"+ds.port, "(ctrl-click link to open)")

	log.Fatal(http.ListenAndServe(ds.port, nil))
	// log.Fatal(http.ListenAndServeTLS(port, "certificate.crt", "private.key", nil)) // TLS
}
