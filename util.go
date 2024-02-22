package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-mail/mail"
)

func getflags() switchdefinitions {
	var days int
	var recipients []string
	var message string
	var auth string
	var files []string
	var owner string
	var key string

	port := ":3451"

	// range over args to get flags
	for index, flag := range os.Args {
		var arg string
		if index < len(os.Args)-1 {
			item := os.Args[index+1]
			if !strings.HasPrefix(item, "-") {
				arg = item
			}
		}

		switch flag {
		case "-days":
			// Set API key
			days, _ = strconv.Atoi(arg)
		case "-message":
			// Set home directory
			message = arg
		case "-owner":
			// Set home directory
			owner = arg
		case "-auth":
			// chats save to homeDir/Saves
			auth = arg
		case "-recipient":
			// chats save to homeDir/Saves
			recipients = append(recipients, arg)
		case "-file":
			// chats save to homeDir/Saves
			files = append(files, arg)
		case "-key":
			key = arg

		case "-port":
			port = ":" + arg
		}
	}

	if days == 0 {
		panic("Error - must be at least 1 day")
	}

	if auth == "" {
		fmt.Println("Warning: No auth phrase set")
	}

	if key == "" {
		panic("No key set")
	}

	if owner == "" {
		panic("Need to set sender/owner")
	}

	if len(recipients) == 0 {
		panic("No recipients set")
	}

	definitions := switchdefinitions{
		days:        days,
		recipients:  recipients,
		message:     message,
		auth:        auth,
		files:       files,
		owner:       owner,
		key:         key,
		port:        port,
		mainTimerCh: make(chan time.Time),
	}

	fmt.Println("Definitions from flags: ", definitions)

	return definitions
}

func (ds switchdefinitions) sendemail(subject, message string) {
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
