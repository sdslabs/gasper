package utils

import "time"

// Scheduler deals with running and managing various tasks at defined intervals
type Scheduler struct {
	interval    time.Duration
	task        func()
	stopTrigger chan bool
	running     bool
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
	if s.running {
		return
	}
	s.running = true
	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-ticker.C:
			s.task()
		case <-s.stopTrigger:
			ticker.Stop()
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
	s.running = false
	s.stopTrigger <- true
}
