package job

import (
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler *gocron.Scheduler
	jobs      *[]Job
}

func NewScheduler(jobs *[]Job) *Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	s := &Scheduler{
		scheduler: &scheduler,
		jobs:      jobs,
	}

	s.addAllJobs()
	return s
}

func (s *Scheduler) addAllJobs() {
	for _, job := range *s.jobs {
		job.AddJob(s.scheduler)
	}
}

func (s *Scheduler) Shutdown() {
	err := (*s.scheduler).Shutdown()
	if err != nil {
		panic(err)
	}
}

func (s *Scheduler) Start() {
	(*s.scheduler).Start()
}
