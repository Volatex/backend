package cron

import (
	"context"
	"time"
)

type CheckPricesJob struct {
	useCase CheckPricesUseCase
	ticker  *time.Ticker
	done    chan bool
}

func NewCheckPricesJob(useCase CheckPricesUseCase) *CheckPricesJob {
	return &CheckPricesJob{
		useCase: useCase,
		ticker:  time.NewTicker(1 * time.Minute),
		done:    make(chan bool),
	}
}

func (j *CheckPricesJob) Start() {
	go func() {
		for {
			select {
			case <-j.ticker.C:
				if err := j.useCase.Execute(context.Background()); err != nil {
					println("Error executing check prices job:", err.Error())
				}
			case <-j.done:
				j.ticker.Stop()
				return
			}
		}
	}()
}

func (j *CheckPricesJob) Stop() {
	j.done <- true
}
