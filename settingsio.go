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
	var switchDef switchdefinitions
	filename := "deadswitchsettings.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		switchDef = getflags() // If there isnt a settings file, read from flags
	} else {
		var saveData saveDefinitions
		err = json.Unmarshal(data, &saveData)
		if err != nil {
			return nil, err
		}

		switchDef = switchdefinitions{
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

	switchDef.mainTimerCh = make(chan time.Time)
	switchDef.halfTimerCh = make(chan time.Time)
	switchDef.quarterTimerCh = make(chan time.Time)

	return &switchDef, nil
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
		days:       days,
		recipients: recipients,
		message:    message,
		auth:       auth,
		files:      files,
		owner:      owner,
		key:        key,
		port:       port,
	}

	fmt.Println("Definitions from flags: ", definitions)

	return definitions
}
