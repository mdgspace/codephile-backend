package db

import (
	"github.com/globalsign/mgo"
	"log"
)

type Collection struct {
	s          *mgo.Session
	db         *mgo.Database
	name       string
	Collection *mgo.Collection
}

func (c *Collection) Connect() {
	c.s = service.Session()
	database := *c.s.DB("")
	c.db = &database
	collection := *c.db.C(c.name)
	c.Collection = &collection
}

func NewCollectionSession(name string) *Collection {
	var c = Collection{
		name: name,
	}
	c.Connect()
	return &c
}

func NewUserCollectionSession() *Collection {
	return NewCollectionSession("coduser")
}

func (c *Collection) Close() {
	service.Close(c)
}

func (c *Collection) DropDatabase() {
	log.Println("Dropping database ...")
	err := c.db.DropDatabase()
	if err != nil {
		panic(err)
	}
}
