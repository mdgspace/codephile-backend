package worker

import (
	"log"
	"sync"

	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/errors"
)

type User bson.ObjectId
type Website string

type Job struct {
	user        User
	websiteName Website
	//handler func which is called when job is performed
	handler func(user bson.ObjectId, website string) error
}

// Keeps track if a job corresponding to a user is already in the queue
type UserQueue struct {
	m map[User]bool
	sync.Mutex
}

var jobQueue chan Job
var userQueue UserQueue

func NewJob(user bson.ObjectId, websiteName string, handler func(user bson.ObjectId, website string) error) Job {
	return Job{user: User(user), websiteName: Website(websiteName), handler: handler}
}

func init() {
	jobQueue = make(chan Job, beego.AppConfig.DefaultInt("MAX_QUEUE_SIZE", 100))
	userQueue = UserQueue{m: map[User]bool{}}
	startWorkerCoRoutines()
}

func work() {
	for job := range jobQueue {
		err := job.handler(bson.ObjectId(job.user), string(job.websiteName))
		if err != nil {
			log.Println("unable to fetch submissions", err.Error())
		}
		userQueue.Lock()
		delete(userQueue.m, job.user)
		userQueue.Unlock()
	}
}

func startWorkerCoRoutines() {
	for i := 0; i < beego.AppConfig.DefaultInt("MAX_WORKER_POOL", 1); i++ {
		go work()
	}
}

func Enqueue(job Job) error {
	// Don't allow if a job corresponding to a user is already present
	if _, ok := userQueue.m[job.user]; ok {
		return errors.ErrJobQueueFull
	}
	select {
	case jobQueue <- job:
		userQueue.Lock()
		defer userQueue.Unlock()
		//jobQueue <- job
		userQueue.m[job.user] = true
		return nil
	//Check if queue is full
	default:
		return errors.ErrJobQueueFull
	}

}
