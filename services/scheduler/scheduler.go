package scheduler

import (
	"github.com/jasonlvhit/gocron"
	"log"
)

func task() {
	 //Updates profiles and submissions
	 log.Println("HELLO")
}

func taskWithParams(a int, b string) {
     //
}

func StartScheduling() {
	gocron.Every(5).Seconds().Do(task)
	gocron.Every(1).Day().Do(task)   //task to be done everyday
	<- gocron.Start()
}