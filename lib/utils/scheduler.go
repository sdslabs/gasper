package utils

import "time"

// Scheduler deals with running and managing various tasks at defined intervals
type Scheduler struct {
	interval time.Duration
	task     func()
	mutex    bool
}

// NewScheduler returns a pointer to a Scheduler object
func NewScheduler(interval time.Duration, task func()) *Scheduler {
	return &Scheduler{
		interval: interval,
		task:     task,
	}
}

// Run starts scheduling the given task
func (s *Scheduler) Run() {
	if s.mutex {
		return
	}
	s.mutex = true
	ticker := time.NewTicker(s.interval)
	for range ticker.C {
		if s.mutex {
			go s.task()
		} else {
			return
		}
	}
}

// RunAsync starts scheduling the given task in a non-blocking manner
func (s *Scheduler) RunAsync() {
	go s.Run()
}

// Terminate stops the scheduler
func (s *Scheduler) Terminate() {
	s.mutex = false
}
