package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func (ds *switchdefinitions) writeSwitchDefinitionsToFile() error {
	filename := "deadswitchsettings.json"

	saveData := saveDefinitions{
		Days:       ds.days,
		Hoursleft:  ds.hoursleft,
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

	// fmt.Println(data)

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Switch definitions written to file:", filename)
	return nil
}

func getDefinitions() (*switchdefinitions, error) {
	switchDef, err := getflags() // get flags first
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting flags, attempting to read file...")
		filename := "deadswitchsettings.json"
		data, err := ioutil.ReadFile(filename)
		var saveData saveDefinitions
		err = json.Unmarshal(data, &saveData)
		if err != nil {
			return nil, err
		}

		fmt.Println("Getting settings from file...")

		switchDef = &switchdefinitions{
			days:       saveData.Days,
			hoursleft:  saveData.Hoursleft,
			recipients: saveData.Recipients,
			message:    saveData.Message,
			auth:       saveData.Auth,
			files:      saveData.Files,
			owner:      saveData.Owner,
			key:        saveData.Key,
			port:       saveData.Port,
		}
	}

	fmt.Println("Deadswitch settings: ", switchDef)

	switchDef.mainTimerCh = make(chan time.Time)
	switchDef.halfTimerCh = make(chan time.Time)
	switchDef.quarterTimerCh = make(chan time.Time)

	return switchDef, nil
}

func getflags() (*switchdefinitions, error) {
	var definitions switchdefinitions
	var err error

	definitions.port = ":3451"

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
			definitions.days, err = strconv.Atoi(arg)
			if err != nil {
				return nil, err
			}
		case "-message":
			definitions.message = arg
		case "-owner":
			definitions.owner = arg
		case "-auth":
			definitions.auth = arg
		case "-recipient":
			definitions.recipients = append(definitions.recipients, arg)
		case "-file":
			definitions.files = append(definitions.files, arg)
		case "-key":
			definitions.key = arg
		case "-port":
			definitions.port = ":" + arg
		}
	}

	var flagsmissing string
	if definitions.days == 0 {
		flagsmissing += "-days atleast1day "
	}

	if definitions.key == "" {
		flagsmissing += "-key 'yourgmailkey' "
	}

	if definitions.owner == "" {
		flagsmissing += "-owner youremailaddress "
	}

	if len(definitions.recipients) == 0 {
		flagsmissing += "-recipient recipientemailaddress "
	}

	if flagsmissing != "" {
		return nil, fmt.Errorf("Flags missing: %s", flagsmissing)
	}

	fmt.Println("Definitions from flags: ", definitions)

	if definitions.auth == "" {
		fmt.Println("Warning: No auth phrase set (set flag -auth resetpassword)")
	}

	if definitions.message == "" {
		fmt.Println("Warning: Message is empty (set flag -message 'your message here')")
	}

	definitions.writeSwitchDefinitionsToFile()

	return &definitions, nil
}
