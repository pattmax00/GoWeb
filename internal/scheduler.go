package app

import (
	"sync"
	"time"
)

type Scheduled struct {
	EveryReboot []func(app *Deps)
	EverySecond []func(app *Deps)
	EveryMinute []func(app *Deps)
	EveryHour   []func(app *Deps)
	EveryDay    []func(app *Deps)
	EveryWeek   []func(app *Deps)
	EveryMonth  []func(app *Deps)
	EveryYear   []func(app *Deps)
}

type Task struct {
	Funcs    []func(app *Deps)
	Interval time.Duration
}

func RunScheduledTasks(app *Deps, poolSize int, stop <-chan struct{}) {
	for _, f := range app.ScheduledTasks.EveryReboot {
		f(app)
	}

	tasks := []Task{
		{Funcs: app.ScheduledTasks.EverySecond, Interval: time.Second},
		{Funcs: app.ScheduledTasks.EveryMinute, Interval: time.Minute},
		{Funcs: app.ScheduledTasks.EveryHour, Interval: time.Hour},
		{Funcs: app.ScheduledTasks.EveryDay, Interval: 24 * time.Hour},
		{Funcs: app.ScheduledTasks.EveryWeek, Interval: 7 * 24 * time.Hour},
		{Funcs: app.ScheduledTasks.EveryMonth, Interval: 30 * 24 * time.Hour},
		{Funcs: app.ScheduledTasks.EveryYear, Interval: 365 * 24 * time.Hour},
	}

	var wg sync.WaitGroup
	runners := make([]chan bool, len(tasks))
	for i, task := range tasks {
		runner := make(chan bool, poolSize)
		runners[i] = runner
		wg.Add(1)
		go func(task Task, runner chan bool) {
			defer wg.Done()
			ticker := time.NewTicker(task.Interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					for _, f := range task.Funcs {
						runner <- true
						go func(f func(app *Deps)) {
							defer func() { <-runner }()
							f(app)
						}(f)
					}
				case <-stop:
					return
				}
			}
		}(task, runner)
	}

	wg.Wait()

	for _, runner := range runners {
		close(runner)
	}
}
