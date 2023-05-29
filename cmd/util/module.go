package util

import (
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

// Exposed interface
type Module struct {
	Name       func() string
	Enable     func(scheduler *gocron.Scheduler) error
	Disable    func(scheduler *gocron.Scheduler)
	Initialise func() error
}

// Internal structure
type ModuleInternal struct {
	ModuleName string

	Channel string
	Sender  SenderOptions

	IntervalFunction   func(*gocron.Scheduler) *gocron.Scheduler
	ExecFunction       func(sendMessage func(body BodyElement))
	InitialiseFunction func() error

	Job    *gocron.Job
	Logger *zap.Logger
}

// RegisterModule options
type ModuleRegistration struct {
	// Module Name
	Name string

	// Target Slack channel for the `sendMessage` wrapper
	Channel string

	// Sender username and image for the `sendMessage` wrapper
	Sender SenderOptions

	// Interval getter
	IntervalFunction func(*gocron.Scheduler) *gocron.Scheduler

	// Execution function
	ExecFunction func(sendMessage func(body BodyElement))

	Initialise func() error
}

func RegisterModule(opts ModuleRegistration) Module {
	module := ModuleInternal{
		ModuleName:         opts.Name,
		Channel:            opts.Channel,
		Sender:             opts.Sender,
		IntervalFunction:   opts.IntervalFunction,
		ExecFunction:       opts.ExecFunction,
		InitialiseFunction: opts.Initialise,
		Logger:             Logger.With(zap.String("module", opts.Name)),
	}

	result := Module{
		Name: func() string { return module.ModuleName },
		Enable: func(scheduler *gocron.Scheduler) error {
			job, err := module.IntervalFunction(scheduler).Do(
				func() {
					sendMessageFn := func(body BodyElement) {
						err := SendMessage(module.Channel, body, &module.Sender)
						if err != nil {
							module.Logger.Warn("unable to send message",
								zap.String("channel", module.Channel),
								zap.Any("sender", module.Sender),
								zap.Error(err),
							)
						}
					}

					module.Logger.Info("executing function")
					module.ExecFunction(sendMessageFn)
				},
			)

			if err != nil {
				return err
			}

			module.Job = job
			module.Logger.Info("module enabled")
			return nil
		},
		Disable: func(scheduler *gocron.Scheduler) {
			scheduler.RemoveByReference(module.Job)
			module.Job = nil
			module.Logger.Info("module disabled")
		},
		Initialise: func() error {
			if module.InitialiseFunction == nil {
				return nil
			}

			err := module.InitialiseFunction()
			if err != nil {
				return err
			}

			module.Logger.Info("module initialised")
			return nil
		},
	}

	return result
}
