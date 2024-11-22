// Package main provides the main entry point for the jobd application
package main

import (
	"flag"
	router "jobd/controllers"
	"jobd/datasource/db"
	"jobd/services"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func init() {
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("stderrthreshold", "DEBUG")
	_ = flag.Set("v", "2")
	flag.Parse()
}

func main() {

	errDB := db.InitDB()
	if errDB != nil {
		log.Fatal(errDB)
	}

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	s.Every(1).Seconds().Do(services.RunTasks)
	s.Every(30).Seconds().Do(services.UpdateSlurmlJobs)
	s.Every(1).Hours().Do(services.ClearOldJobs)
	s.StartAsync()

	r := router.SetupRouter()
	_ = r.Run()

}
