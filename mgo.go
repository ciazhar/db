package db

import (
	"github.com/ciazhar/config"
	"gopkg.in/mgo.v2"
	"log"
)

var Mongo *mgo.Database

type MongoDB struct {
	Host   string
	Database string
}

func Init(c *config.Config) {
	m := MongoDB{}

	m.Host = c.Get("database").Get("host").String()
	m.Database = c.Get("database").Get("name").String()
	session, err := mgo.Dial(m.Host)
	if err != nil {
		log.Fatal(err)
	}
	Mongo = session.DB(m.Database)
}