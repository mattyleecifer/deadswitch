package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
		fmt.Println("Failed to get definitions - make sure all flags are set or settings file is present.")
		fmt.Println(err)
		os.Exit(1)
	}

	go ds.mainTimer() // start main timer
	// start secondary timers
	go ds.halfTimer()
	go ds.quarterTimer()
	go ds.writeToFileEveryHour() // basically keeps it updated with hours left

	// start server
	http.HandleFunc("/", ds.hauth)
	fmt.Println("Running GUI on http://127.0.0.1"+ds.port, "(ctrl-click link to open)")

	log.Fatal(http.ListenAndServe(ds.port, nil))
	// log.Fatal(http.ListenAndServeTLS(port, "certificate.crt", "private.key", nil)) // TLS
}

func (ds *switchdefinitions) sendemail(subject, message string, alert bool) {
	m := mail.NewMessage()

	m.SetHeader("From", ds.owner)

	if alert {
		m.SetHeader("To", ds.owner)
	} else {
		m.SetHeader("To", ds.recipients...)
		for _, filename := range ds.files {
			m.Attach(filename)
		}
	}

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", message)

	d := mail.NewDialer("smtp.gmail.com", 587, ds.owner, ds.key)

	if err := d.DialAndSend(m); err != nil {

		panic(err)

	}
}
