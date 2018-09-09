package db

import (
	"github.com/ciazhar/config"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/ciazhar/animus/model"
	"gopkg.in/mgo.v2/bson"
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

func (m *MongoDB)Find(collection string, query interface{},skip,limit int,sort string) *mgo.Query {
	return Mongo.C(collection).Find(query).Skip((skip-1)*limit).Limit(limit).Sort(sort)
}

func (m *MongoDB)FindId(collection string,id string) *mgo.Query {
	return Mongo.C(collection).FindId(bson.ObjectIdHex(id))
}

func (m *MongoDB)Insert(collection string, payload ...interface{}) error {
	for _,e := range payload {
		if err := Mongo.C(collection).Insert(e);err!=nil{
			return err
		}
	}
	return nil
}

func (m *MongoDB)RemoveId(collection string,id string) error {
	err := Mongo.C(collection).RemoveId(bson.ObjectIdHex(id))
	return err
}

func (m *MongoDB)UpdateId(collection string,anime *model.Anime) error  {
	err := Mongo.C(collection).UpdateId(anime.ID, &anime)
	return err
}

func (m *MongoDB)CreateIndex(collection string, attr ...string) error {
	index := mgo.Index{
		Key:        attr,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := Mongo.C(collection).EnsureIndex(index)
	return err
}
