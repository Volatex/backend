package cron

import (
	"context"
	"time"
)

type PrintStrategiesUseCase interface {
	Execute(ctx context.Context) error
}

type PrintStrategiesJob struct {
	useCase PrintStrategiesUseCase
	ticker  *time.Ticker
	done    chan bool
}

func NewPrintStrategiesJob(useCase PrintStrategiesUseCase) *PrintStrategiesJob {
	return &PrintStrategiesJob{
		useCase: useCase,
		ticker:  time.NewTicker(1 * time.Minute),
		done:    make(chan bool),
	}
}

func (j *PrintStrategiesJob) Start() {
	go func() {
		for {
			select {
			case <-j.ticker.C:
				if err := j.useCase.Execute(context.Background()); err != nil {
					println("Error executing print strategies job:", err.Error())
				}
			case <-j.done:
				j.ticker.Stop()
				return
			}
		}
	}()
}

func (j *PrintStrategiesJob) Stop() {
	j.done <- true
}
