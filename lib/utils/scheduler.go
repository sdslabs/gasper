package utils

import "time"

// Scheduler deals with running and managing various tasks at defined intervals
type Scheduler struct {
	interval time.Duration
	task     func()
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
	ticker := time.NewTicker(s.interval)
	for range ticker.C {
		go s.task()
	}
}

// RunAsync starts scheduling the given task in a non-blocking manner
func (s *Scheduler) RunAsync() {
	go s.Run()
}
