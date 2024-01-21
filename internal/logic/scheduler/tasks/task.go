package tasks

import (
	"errors"
	"fmt"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TaskUpdateClans string = "UPDATE_CLANS"

	TaskRecordSessions           = "RECORD_ACCOUNT_SESSIONS"
	TaskUpdateAccountWN8         = "UPDATE_ACCOUNT_WN8"
	TaskRecordPlayerAchievements = "UPDATE_ACCOUNT_ACHIEVEMENTS"
)

var taskHandlers = make(map[string]TaskHandler)

type TaskHandler struct {
	RetryOnFail func(*Task) bool
	Process     func(*Task) (string, error)
}

func registerTaskHandler(kind string, handler TaskHandler) {
	if _, ok := taskHandlers[kind]; ok {
		panic(fmt.Sprintf("task handler for %s already registered", kind))
	}
	taskHandlers[kind] = handler
}

// Task statuses
type taskStatus string

const TaskStatusScheduled taskStatus = "TASK_SCHEDULED"
const TaskStatusInProgress taskStatus = "TASK_IN_PROGRESS"
const TaskStatusComplete taskStatus = "TASK_COMPLETE"
const TaskStatusFailed taskStatus = "TASK_FAILED"

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Type      string             `bson:"kind"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`

	Targets []int `bson:"targets"`

	Logs []AttemptLog `bson:"logs"`

	Status         taskStatus `bson:"status"`
	ScheduledAfter time.Time  `bson:"scheduled_after"`
	LastAttempt    time.Time  `bson:"last_attempt"`

	Data map[string]any `bson:"data"`
}

func (t *Task) LogAttempt(log AttemptLog) {
	t.Logs = append(t.Logs, log)
}

func (t *Task) OnCreated() {
	t.LastAttempt = time.Now()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}
func (t *Task) OnUpdated() {
	t.UpdatedAt = time.Now()
}

func (t *Task) RetryOnFail() bool {
	handlers, ok := taskHandlers[t.Type]
	if !ok {
		return false
	}

	t.LastAttempt = time.Now()
	return handlers.RetryOnFail(t)
}

func (t *Task) Process() (string, error) {
	handlers, ok := taskHandlers[t.Type]
	if !ok {
		return "", fmt.Errorf("no handler for task type %s", t.Type)
	}

	t.LastAttempt = time.Now()
	return handlers.Process(t)
}

type AttemptLog struct {
	Targets   []int     `json:"targets" bson:"targets"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Comment   string    `json:"result" bson:"result"`
	Error     string    `json:"error" bson:"error"`
}

func NewAttemptLog(task Task, comment, err string) AttemptLog {
	return AttemptLog{
		Targets:   task.Targets,
		Timestamp: time.Now(),
		Comment:   comment,
		Error:     err,
	}
}

/*
Retrieves all target IDs from the database based on task.Type() and filter, creates a new task in queue.
  - If splitTaskFn is provided, it will split the task into subtasks.
*/
func CreateBulkTask(filter bson.M, task Task, splitTaskFn func(Task) []Task) (err error) {
	if len(task.Targets) != 0 {
		return errors.New("target IDs already set")
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	switch task.Type {
	case TaskUpdateClans:
		// Get all clan IDs on the realm
		result, err := database.DefaultClient.Collection(database.CollectionClans).Distinct(ctx, "_id", filter)
		if err != nil {
			return err
		}

		var targetIDs []int
		for _, id := range result {
			if idInt, ok := id.(int); ok {
				targetIDs = append(targetIDs, idInt)
			} else {
				log.Error().Msgf("invalid clan ID type: %T", id)
			}
		}

		task.Targets = targetIDs

	case TaskRecordPlayerAchievements:
		// All players on the realm
		fallthrough
	case TaskUpdateAccountWN8:
		// All players on the realm
		fallthrough
	case TaskRecordSessions:
		// Get all player IDs on the realm
		result, err := database.DefaultClient.Collection(database.CollectionAccounts).Distinct(ctx, "_id", filter)
		if err != nil {
			return err
		}

		var targetIDs []int
		for _, id := range result {
			switch cast := id.(type) {
			case int:
				targetIDs = append(targetIDs, cast)
			case int32:
				targetIDs = append(targetIDs, int(cast))
			case int64:
				targetIDs = append(targetIDs, int(cast))
			default:
				log.Error().Msgf("invalid player ID %v type: %T", cast, id)
			}
		}

		task.Targets = targetIDs

	default:
		return errors.New("invalid task type")
	}

	if len(task.Targets) == 0 {
		return fmt.Errorf("no targets found for task type %s and filter %+v", task.Type, filter)
	}

	if splitTaskFn != nil {
		return CreateTasks(splitTaskFn(task)...)
	}
	return CreateTasks(task)
}

func CreateTasks(tasks ...Task) error {
	var writes []mongo.WriteModel
	for _, task := range tasks {
		if len(task.Targets) == 0 {
			return errors.New("task targets not set")
		}

		task.OnCreated()
		task.Status = TaskStatusScheduled
		writes = append(writes, mongo.NewInsertOneModel().SetDocument(task))
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionTasks).BulkWrite(ctx, writes)
	if err != nil {
		return err
	}

	return err
}

func UpdateTasks(tasks ...Task) error {
	var writes []mongo.WriteModel
	for _, task := range tasks {
		if task.ID.IsZero() {
			return errors.New("task ID not set")
		}
		task.OnUpdated()
		writes = append(writes, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": task.ID}).SetUpdate(bson.M{"$set": task}))
	}

	if len(writes) == 0 {
		return nil
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionTasks).BulkWrite(ctx, writes)
	if err != nil {
		return err
	}

	return nil
}

/*
Retrieves all tasks with status TaskStatusScheduled that match filter and updates their status to TaskStatusInProgress.
*/
func StartScheduledTasks(filter bson.M, limit int) ([]Task, error) {
	if filter == nil {
		filter = bson.M{}
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var pipeline mongo.Pipeline
	filter["status"] = TaskStatusScheduled
	filter["scheduled_after"] = bson.M{"$lte": time.Now()}
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit}})
	pipeline = append(pipeline, bson.D{{Key: "$set", Value: bson.M{"status": TaskStatusInProgress}}})

	var tasks []Task
	cur, err := database.DefaultClient.Collection(database.CollectionTasks).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return tasks, cur.All(ctx, &tasks)
}
