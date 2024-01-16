package tasks

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var DefaultQueue = NewQueue(10)

type Queue struct {
	concurrencyLimit int
	limiter          chan struct{}
	lastTaskRun      time.Time
}

func (q *Queue) ConcurrencyLimit() int {
	return q.concurrencyLimit
}

func (q *Queue) ActiveWorkers() int {
	return len(q.limiter)
}

func (q *Queue) LastTaskRun() time.Time {
	return q.lastTaskRun
}

func NewQueue(concurrencyLimit int) *Queue {
	return &Queue{
		concurrencyLimit: concurrencyLimit,
		limiter:          make(chan struct{}, concurrencyLimit),
	}
}

func (q *Queue) Process(tasks ...Task[any]) error {
	if len(tasks) == 0 {
		log.Debug().Msg("no tasks to process")
		return nil
	}

	log.Debug().Msgf("processing %d tasks", len(tasks))

	var wg sync.WaitGroup
	q.lastTaskRun = time.Now()
	processedTasks := make(chan Task[any], len(tasks))
	for _, task := range tasks {
		wg.Add(1)
		go func(t Task[any]) {
			q.limiter <- struct{}{}
			defer func() {
				processedTasks <- t
				wg.Done()
				<-q.limiter
				log.Debug().Msgf("finished processing task %s", t.ID())
			}()
			log.Debug().Msgf("processing task %s", t.ID())

			comment, err := t.Process()
			if err != nil {
				t.SetStatus(TaskStatusFailed)
			} else {
				t.SetStatus(TaskStatusComplete)
			}

			attempt := AttemptLog{
				Timestamp: time.Now(),
				Targets:   t.Targets(),
				Comment:   comment,
				Error:     err,
			}
			t.LogAttempt(attempt)
		}(task)
	}

	wg.Wait()
	close(processedTasks)

	rescheduledCount := 0
	processedSlice := make([]Task[any], 0, len(processedTasks))
	for task := range processedTasks {
		if task.Status() == TaskStatusFailed && task.RetryOnFail() {
			rescheduledCount++
			task.SetStatus(TaskStatusScheduled)
		}
		processedSlice = append(processedSlice, task)
	}

	err := UpdateTasks(processedSlice...)
	if err != nil {
		return err
	}

	log.Debug().Msgf("processed %d tasks, %d rescheduled", len(processedSlice), rescheduledCount)
	return nil
}