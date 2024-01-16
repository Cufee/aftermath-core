package tasks

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type SessionsUpdateTaskData struct {
	realm     string `bson:"realm"`
	triesLeft int    `bson:"tries_left"`
}

func (data *SessionsUpdateTaskData) RetryOnFail(task *Task[any]) bool {
	task.SetScheduledAfter(time.Now().Add(5 * time.Minute)) // Backoff for 5 minutes to avoid spamming
	data.triesLeft -= 1
	return data.triesLeft > 0
}

func (data SessionsUpdateTaskData) Process(task *Task[any]) (string, error) {
	time.Sleep(5 * time.Second) // Simulate processing time

	if task.Targets()[0] == 1 {
		return "simulate error", errors.New("fake error")
	}

	return "this was a dummy task", nil
}

func CreateSessionUpdateTasks(realm string) error {
	task := Task[any]{
		Kind: TaskUpdateClans,
		Data: &SessionsUpdateTaskData{
			realm:     realm,
			triesLeft: 3,
		},
	}
	// This update requires (2 + n) requests per n players
	return CreateBulkTask(bson.M{"realm": realm}, task, splitTasksByTargets(50))
}
