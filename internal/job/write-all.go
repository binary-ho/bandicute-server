package job

import (
	"bandicute-server/pkg/logger"
	"context"
	"github.com/go-co-op/gocron/v2"
	"time"
)

type Job interface {
	AddJob(scheduler *gocron.Scheduler)
	Cancel()
}

const writeJobDuration = 4 * time.Hour

type WriteAllMemberPostJob struct {
	duration      time.Duration
	writeFunction func(ctx context.Context)
}

func NewWriteAllMemberPostJob(writeFunction func(ctx context.Context)) *WriteAllMemberPostJob {
	return &WriteAllMemberPostJob{
		duration:      writeJobDuration,
		writeFunction: writeFunction,
	}
}

func (j *WriteAllMemberPostJob) AddJob(scheduler *gocron.Scheduler) {
	_, err := (*scheduler).NewJob(
		gocron.DurationJob(j.duration),
		gocron.NewTask(func() {
			logger.Info("start WriteAllMemberPostJob by batch", nil)
			ctx, cancelFunction := context.WithTimeout(context.Background(), j.duration)
			j.writeFunction(ctx)
			defer cancelFunction()
			logger.Info("end WriteAllMemberPostJob by batch", nil)
		}),
	)

	if err != nil {
		panic(err)
	}
}
