// Package cron This provides the Kernel a managed Cron service.
//
// This is simply a wrapper around gopkg.in/robfig/cron.v2
//
package cron

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-kernel/v2/util/task"
	crn "gopkg.in/robfig/cron.v2"
)

// CronService gopkg.in/robfig/crn.v2 as a Kernel Service
type CronService struct {
	daemon *kernel.Daemon `kernel:"inject"`
	worker task.Queue     `kernel:"worker"`
	cron   *crn.Cron
}

func (s *CronService) PostInit() error {
	s.cron = crn.New()

	// Mark ourselves as a daemon
	s.daemon.SetDaemon()
	return nil
}

func (s *CronService) Start() error {
	s.cron.Start()
	return nil
}

func (s *CronService) Stop() {
	s.cron.Stop()
}

func (s *CronService) AddFunc(spec string, cmd func()) (crn.EntryID, error) {
	id, err := s.cron.AddFunc(spec, cmd)
	return id, err
}

func (s *CronService) AddJob(spec string, cmd crn.Job) (crn.EntryID, error) {
	id, err := s.cron.AddJob(spec, cmd)
	return id, err
}

func (s *CronService) AddTask(spec string, task task.Task) (crn.EntryID, error) {
	return s.AddFunc(spec, func() {
		s.worker.AddTask(task)
	})
}

func (s *CronService) AddPriorityTask(priority int, spec string, task task.Task) (crn.EntryID, error) {
	return s.AddFunc(spec, func() {
		s.worker.AddPriorityTask(priority, task)
	})
}

func (s *CronService) Entries() []crn.Entry {
	return s.cron.Entries()
}

func (s *CronService) Entry(id crn.EntryID) crn.Entry {
	return s.cron.Entry(id)
}

func (s *CronService) Remove(id crn.EntryID) {
	s.cron.Remove(id)
}

func (s *CronService) Schedule(schedule crn.Schedule, cmd crn.Job) crn.EntryID {
	return s.cron.Schedule(schedule, cmd)
}
