package main

import (
	"context"
	"sync"
	"time"

	"github.com/tombenke/go-12f-common/apprun"
	"github.com/tombenke/go-12f-common/examples/scheduler/timer"
	"github.com/tombenke/go-12f-common/examples/scheduler/worker"
	"github.com/tombenke/go-12f-common/healthcheck"
	"github.com/tombenke/go-12f-common/log"
)

type Application struct {
	config        *Config
	wg            *sync.WaitGroup
	err           error
	currentTimeCh chan (time.Time)

	// The internal components of the application
	components []apprun.LifecycleManager
}

func (a *Application) Startup(ctx context.Context, wg *sync.WaitGroup) error {
	log.Logger.Infof("Application Startup")

	// Inject the central waitgroup into the application object
	a.wg = wg

	// Startup the internal components
	for c := range a.components {
		if err := a.components[c].Startup(context.Background(), a.wg); err != nil {
			return err
		}

	}

	// Check if Application is ready to serve service calls
	_ = a.Check(context.Background())

	a.err = nil
	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	log.Logger.Infof("Application Shutdown")
	a.closeChannels()

	a.err = healthcheck.ServiceNotAvailableError{}

	// Shutdown the internal components
	for c := range a.components {
		if err := a.components[c].Shutdown(context.Background()); err != nil {
			return err
		}

	}
	return nil
}

// Close the inter-component communication channels
func (a *Application) closeChannels() {
	if a.currentTimeCh != nil {
		close(a.currentTimeCh)
	}
}

func (a *Application) Check(ctx context.Context) error {
	log.Logger.Infof("Application Check")
	return a.err
}

func NewApplication(config *Config) (apprun.LifecycleManager, error) {
	log.Logger.Infof("Application.Config: %+v", *config)
	// Create channel(s) for inter-component communication
	currentTimeCh := make(chan (time.Time))

	return &Application{
		config: config,
		err:    healthcheck.ServiceNotAvailableError{},
		components: []apprun.LifecycleManager{
			timer.NewTimer(&config.timer, currentTimeCh),
			worker.NewWorker(&config.worker, currentTimeCh),
		},
	}, nil
}
