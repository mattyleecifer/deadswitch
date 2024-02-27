package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-mail/mail"
)

type switchdefinitions struct {
	days           int
	hoursleft      int
	recipients     []string
	message        string
	auth           string // password for deadswitch
	files          []string
	owner          string // to send reminder email
	key            string // gmail key
	port           string
	mainTimerCh    chan time.Time
	halfTimerCh    chan time.Time
	quarterTimerCh chan time.Time
}

type saveDefinitions struct {
	Days       int
	Hoursleft  int // restart should only read this
	Recipients []string
	Message    string
	Auth       string
	Files      []string
	Owner      string
	Key        string
	Port       string
}

func main() {
	ds, err := getDefinitions() // this checks flags as well
	if err != nil {
		fmt.Println(err)
	}

	ds.writeSwitchDefinitionsToFile() // immediately save data

	go ds.mainTimer()            // start main timer
	go ds.writeToFileEveryHour() // basically keeps it updated with days left
	// start secondary timers
	go ds.halfTimer()
	go ds.quarterTimer()

	// start server
	http.HandleFunc("/", ds.hauth)
	fmt.Println("Running GUI on http://127.0.0.1"+ds.port, "(ctrl-click link to open)")

	log.Fatal(http.ListenAndServe(ds.port, nil))
	// log.Fatal(http.ListenAndServeTLS(port, "certificate.crt", "private.key", nil)) // TLS
}

func (ds *switchdefinitions) sendemail(subject, message string) {
	m := mail.NewMessage()

	m.SetHeader("From", ds.owner)

	m.SetHeader("To", ds.recipients...)

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", message)

	for _, filename := range ds.files {
		m.Attach(filename)
	}

	d := mail.NewDialer("smtp.gmail.com", 587, ds.owner, ds.key)

	if err := d.DialAndSend(m); err != nil {

		panic(err)

	}
}
