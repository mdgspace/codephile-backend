package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Service struct {
	baseClient  *mongo.Client
	queue       chan int
	URL         string
	Open        int
}

var service Service

func (s *Service) New() error {
	var err error
	s.queue = make(chan int, maxPool)
	for i := 0; i < maxPool; i = i + 1 {
		s.queue <- 1
	}
	s.Open = 0
	s.baseClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(s.URL))
	if err = s.baseClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	return err
}

func (s *Service) Client() *mongo.Client {
	<-s.queue
	s.Open++
	return s.baseClient
}

func (s *Service) Close(c *Collection) {
	s.queue <- 1
	s.Open--
}
