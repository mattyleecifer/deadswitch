package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-mail/mail"
)

type switchdefinitions struct {
	days        int
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
	Recipients []string
	Message    string
	Auth       string
	Files      []string
	Owner      string
	Key        string
	Port       string
}

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

func main() {
	ds, err := readSwitchDefinitionsFromFile() // this checks flags as well
	if err != nil {
		fmt.Println(err)
	}

	ds.writeSwitchDefinitionsToFile() // immediately save data

	go ds.mainTimer()

	http.HandleFunc("/", ds.hauth)
	http.HandleFunc("/auth/", ds.hauth)
	fmt.Println("Running GUI on http://127.0.0.1"+ds.port, "(ctrl-click link to open)")

	log.Fatal(http.ListenAndServe(ds.port, nil))
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

func (ds switchdefinitions) hauth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth := r.FormValue("auth")
		if auth == ds.auth {
			// reset timer
			// test := <-ds.timerCh
			// fmt.Println(test)
			fmt.Println("resetting")
			ds.resetTimers()
		}
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		render(w, `<!DOCTYPE html>
		<html>
			<head>
				<meta name="HandheldFriendly" content="true" />
				<meta name="MobileOptimized" content="320" />
				<meta name="viewport" content="initial-scale=1.0, maximum-scale=1.0, width=device-width, user-scalable=no" />
				<script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>
				<title>Dead Switch</title>
			</head>
			<body>
				<div class="container">
					<div class="centertext" id="inputbox">
						Input auth:<br>
						<input id="auth" name="auth" type="password">
						<button hx-post="/auth/" hx-target="html" hx-swap="none" hx-include="#auth" hx-trigger="click, keydown[keyCode==13&&shiftKey!=true] from:#inputbox">Submit</button>
					</div>
				</div>
			</body>
		</html>`, nil)
	}
}

func render(w http.ResponseWriter, html string, data any) {
	// Render the HTML template
	// fmt.Println("Rendering...")
	w.WriteHeader(http.StatusOK)
	tmpl, err := template.New(html).Parse(html)
	if err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ds switchdefinitions) resetTimers() {
	ds.mainTimerCh <- time.Now()
}

func (ds switchdefinitions) writeSwitchDefinitionsToFile() error {
	filename := "deadswitchsettings.json"

	saveData := saveDefinitions{
		Days:       ds.days,
		Recipients: ds.recipients,
		Message:    ds.message,
		Auth:       ds.auth,
		Files:      ds.files,
		Owner:      ds.owner,
		Key:        ds.key,
		Port:       ds.port,
	}

	data, err := json.Marshal(saveData)
	if err != nil {
		return err
	}

	fmt.Println(data)

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Switch definitions written to file:", filename)
	return nil
}

func readSwitchDefinitionsFromFile() (*switchdefinitions, error) {
	filename := "deadswitchsettings.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		switchDef := getflags()
		return &switchDef, err
	}
	var saveData saveDefinitions
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		return nil, err
	}

	switchDef := switchdefinitions{
		days:        saveData.Days,
		recipients:  saveData.Recipients,
		message:     saveData.Message,
		auth:        saveData.Auth,
		files:       saveData.Files,
		owner:       saveData.Owner,
		key:         saveData.Key,
		port:        saveData.Port,
		mainTimerCh: make(chan time.Time),
	}

	return &switchDef, nil
}
