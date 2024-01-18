package tasks

import (
	"errors"
	"strings"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	registerTaskHandler(TaskRecordSessions, TaskHandler{
		Process: func(task *Task) (string, error) {
			if task.Data == nil {
				return "no data provided", errors.New("no data provided")
			}
			realm, ok := task.Data["realm"].(string)
			if !ok {
				return "invalid realm", errors.New("invalid realm")
			}

			accountErrs, err := cache.RefreshSessionsAndAccounts(models.SessionTypeDaily, realm, task.Targets...)
			if err != nil {
				return "failed to refresh sessions on all account", err
			}

			if len(accountErrs) == 0 {
				return "finished session update on all accounts", nil
			}

			var failedAccounts []int
			for accountId, err := range accountErrs {
				if err != nil {
					failedAccounts = append(failedAccounts, accountId)
				}
			}

			// Retry failed accounts
			task.Targets = failedAccounts
			return "retrying failed accounts", errors.New("some accounts failed")
		},
		RetryOnFail: func(task *Task) bool {
			triesLeft, ok := task.Data["triesLeft"].(int32)
			if !ok {
				return false
			}
			if triesLeft <= 0 {
				return false
			}

			triesLeft -= 1
			task.Data["triesLeft"] = triesLeft
			task.ScheduledAfter = time.Now().Add(5 * time.Minute) // Backoff for 5 minutes to avoid spamming
			return true
		},
	})
}

func CreateSessionUpdateTasks(realm string) error {
	realm = strings.ToUpper(realm)
	task := Task{
		Type: TaskRecordSessions,
		Data: map[string]any{
			"realm":     realm,
			"triesLeft": int32(3),
		},
	}
	// This update requires (2 + n) requests per n players
	return CreateBulkTask(bson.M{"realm": realm}, task, splitTasksByTargets(50))
}
