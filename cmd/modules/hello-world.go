package modules

import (
	. "gradbot/util"

	"github.com/go-co-op/gocron"
)

func HelloWorld() Module {
	return RegisterModule(ModuleRegistration{
		Name:    "Hello World",
		Channel: "C04M53GKLH3",
		Sender: SenderOptions{
			Name:  "Hello World",
			Emoji: "indestructible",
		},
		IntervalFunction: func(s *gocron.Scheduler) *gocron.Scheduler {
			return s.Every(10).Second()
		},
		ExecFunction: func(sendMessage func(body BodyElement)) {
			sendMessage(PlainMessage("Hello, World!"))
		},
	})
}
