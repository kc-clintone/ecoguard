package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

var (
	Users     []User
	usersLock = &sync.Mutex{}
	usersFile = "users.json"
)

func LoadUsers() error {
	if _, err := os.Stat(usersFile); os.IsNotExist(err) {
		Users = []User{}
		return nil
	}
	data, err := ioutil.ReadFile(usersFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Users)
}

func SaveUsers() error {
	data, err := json.MarshalIndent(Users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(usersFile, data, 0644)
}