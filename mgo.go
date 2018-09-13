package db

import (
	"github.com/ciazhar/config"
	"gopkg.in/mgo.v2"
	"log"
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

func CreateIndex(collection string, attr ...string) error {
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

func Find(collection string, query interface{},payload ...interface{}) error {
	return Mongo.C(collection).Find(query).All(&payload)
}

func FindWithPagingAndSorting(collection string, query interface{},skip,limit int,sort string,payload ...interface{}) error {
	return Mongo.C(collection).Find(query).Skip((skip-1)*limit).Limit(limit).Sort(sort).All(&payload)
}

func FindId(collection string,id string, payload ...interface{}) error {
	return Mongo.C(collection).FindId(bson.ObjectIdHex(id)).One(&payload)
}

func Insert(collection string, payload ...interface{}) error {
	for _,e := range payload {
		if err := Mongo.C(collection).Insert(e);err!=nil{
			return err
		}
	}
	return nil
}

func UpdateId(collection string,id string, q interface{}) error  {
	err := Mongo.C(collection).UpdateId(bson.ObjectIdHex(id), &q)
	return err
}

func SoftDelete(collection string,id string) error {
	var record interface{}
	if err := Mongo.C(collection).FindId(bson.ObjectIdHex(id)).One(&record); err!=nil{
		return err
	}

	err := Mongo.C(collection).UpdateId(bson.ObjectIdHex(id), &record)
	return err
}

func RemoveId(collection string,id string) error {
	err := Mongo.C(collection).RemoveId(bson.ObjectIdHex(id))
	return err
}
