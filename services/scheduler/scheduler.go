package scheduler

import (
	"github.com/jasonlvhit/gocron"
	"github.com/mdg-iitr/Codephile/models"
	// "github.com/mdg-iitr/Codephile/models"

)

func task() {
	 //Updates profiles and submissions
	 _ = models.RefreshSubmissions()
	 //handle error
}

func taskWithParams(a int, b string) {
     //
}

func StartScheduling() {
	// task()
	// gocron.Every(10).Seconds().Do(task)
	gocron.Every(1).Day().Do(task)   //task to be done everyday
	<- gocron.Start()
}