package db

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/GuillaumeOuint/telegram-discord-bridge/pkg/types"
)

type DB struct {
	Messages *[]*types.Message
	Mutex    *sync.Mutex
}

func NewDB() *DB {
	mutex := &sync.Mutex{}
	messages := []*types.Message{}
	db := &DB{
		Messages: &messages,
		Mutex:    mutex,
	}
	db.Load()
	go db.GC()
	go db.SaveLoop()
	return db
}

func (db *DB) AddMessage(message *types.Message) {
	db.Mutex.Lock()
	*db.Messages = append(*db.Messages, message)
	db.Mutex.Unlock()
}

func (db *DB) GetMessages() []*types.Message {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	return *db.Messages
}

func (db *DB) Save() {
	// save to messages.db
	db.Mutex.Lock()
	filename := os.Getenv("MESSAGE_DB")
	if filename == "" {
		filename = "messages.db"
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(db.Messages)
	if err != nil {
		panic(err)
	}
	db.Mutex.Unlock()
}

func (db *DB) Load() {
	db.Mutex.Lock()
	// load from messages.db
	filename := os.Getenv("MESSAGE_DB")
	if filename == "" {
		filename = "messages.db"
	}
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		db.Mutex.Unlock()
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	err = dec.Decode(db.Messages)
	if err != nil {
		panic(err)
	}
	db.Mutex.Unlock()
}

func (db *DB) GC() {
	for {
		time.Sleep(5 * time.Minute)
		db.Mutex.Lock()
		for i, message := range *db.Messages {
			dur := os.Getenv("MESSAGE_TTL")
			if dur == "" {
				dur = "72h"
			}
			ttl, err := time.ParseDuration(dur)
			if err != nil {
				panic(err)
			}
			if message.Date.Before(time.Now().Add(-ttl)) {
				*db.Messages = append((*db.Messages)[:i], (*db.Messages)[i+1:]...)
			}
		}
		db.Mutex.Unlock()
		fmt.Printf("MessageDB size: %v\n", len(*db.Messages))
	}
}

func (db *DB) SaveLoop() {
	for {
		time.Sleep(1 * time.Minute)
		db.Save()
	}
}
