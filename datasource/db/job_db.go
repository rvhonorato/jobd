package db

import (
	"os"

	"github.com/golang/glog"
	scribble "github.com/nanobox-io/golang-scribble"
)

var Client *scribble.Driver

var NAME = os.Getenv("DB_PATH")

func init() {
	if NAME == "" {
		glog.Warning("DB_PATH not set, using default `./db`")
		NAME = "./db"
	}
}

func InitDB() error {

	database, _ := scribble.New(NAME, nil)

	Client = database
	return nil
}
