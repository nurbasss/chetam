package main

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

//connection to mongo db
func connect() (*mgo.Database, error) {
	//host := "mongo:27017"
	host := "mongodb://localhost:27100"
	/* poka localno zapuskaite docker 
	ewe ne nastroil dlya go app 
	no mongo na dockere*/
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
