package main

import (
	"os"
	"time"

	"gradbot/modules"
	. "gradbot/util"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

func main() {
	InitLogger()

	timezoneString := os.Getenv("TIMEZONE")
	if len(timezoneString) == 0 {
		Logger.Fatal("timezone environment variable not set")
	}

	loc, err := time.LoadLocation(timezoneString)
	if err != nil {
		Logger.Fatal("could not parse timezone", zap.Error(err))
	}

	scheduler := gocron.NewScheduler(loc)

	mods := []Module{modules.HelloWorld()}

	for _, module := range mods {
		module.Initialise()
	}

	for _, module := range mods {
		module.Enable(scheduler)
	}

	Logger.Info("Starting scheduler")
	scheduler.StartBlocking()

	// web.Test()
}
