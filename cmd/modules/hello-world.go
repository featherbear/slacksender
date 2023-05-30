package modules

import (
	. "gradbot/util"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

func helloWorld() Module {
	return RegisterModule(ModuleRegistration{
		Name:    "Hello World",
		Channel: "#awong6-test-private",
		Sender: SenderOptions{
			Name:  "Hello World",
			Emoji: "indestructible",
		},
		IntervalFunction: func(s *gocron.Scheduler) *gocron.Scheduler {
			return s.Every(1).Minute()
		},
		ExecFunction: func(Logger *zap.Logger, sendMessage func(body BodyElement)) {
			sendMessage(PlainMessage("Hello, World!"))
		},
	})
}
