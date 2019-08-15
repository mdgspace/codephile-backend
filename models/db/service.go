package db

import (
	"crypto/tls"
	"github.com/globalsign/mgo"
	"net"
	"time"
)

type Service struct {
	baseSession *mgo.Session
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
	dialInfo, err := mgo.ParseURL(s.URL)
	if err != nil {
		panic(err)
	}
	dialInfo.Timeout = 10 * time.Second
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	s.baseSession, err = mgo.DialWithInfo(dialInfo)
	return err
}

func (s *Service) Session() *mgo.Session {
	<-s.queue
	s.Open++
	return s.baseSession.Copy()
}

func (s *Service) Close(c *Collection) {
	c.db.s.Close()
	s.queue <- 1
	s.Open--
}
