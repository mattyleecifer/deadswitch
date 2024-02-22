package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

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
