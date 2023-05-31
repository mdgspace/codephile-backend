package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	s          *mongo.Client
	db         *mongo.Database
	name       string
	Collection *mongo.Collection
}

func (c *Collection) Connect() {
	c.s = service.Client()
	database := *c.s.Database("codephile")
	c.db = &database
	collection := *c.db.Collection(c.name)
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
	err := c.db.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
}
