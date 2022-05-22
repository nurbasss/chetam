package main

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

//connection to mongo db
func connect() (*mgo.Database, error) {
	host := "mongo:27017"
	dbName := "chetam"

	if session, err := mgo.Dial(host); err != nil {
		return nil, err
	} else {
		if err := session.Ping(); err != nil {
			return nil, err
		}

		db := session.DB(dbName)

		fmt.Println("chetam teper  db rabotaet!")

		return db, nil
	}
}
