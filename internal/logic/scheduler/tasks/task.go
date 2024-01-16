package tasks

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Task types
type taskType string

const TaskUpdateClans taskType = "UPDATE_CLANS"

const TaskRecordSessions taskType = "RECORD_ACCOUNT_SESSIONS"

const TaskUpdateAccountWN8 taskType = "UPDATE_ACCOUNT_WN8"
const TaskRecordPlayerAchievements taskType = "UPDATE_ACCOUNT_ACHIEVEMENTS"

// Task statuses
type taskStatus string

const TaskStatusScheduled taskStatus = "TASK_SCHEDULED"
const TaskStatusInProgress taskStatus = "TASK_IN_PROGRESS"
const TaskStatusComplete taskStatus = "TASK_COMPLETE"
const TaskStatusFailed taskStatus = "TASK_FAILED"

type TaskData[T any] interface {
	Process(*Task[T]) (string, error)
	RetryOnFail(*Task[T]) bool
}

type Task[T any] struct {
	id        primitive.ObjectID `bson:"_id,omitempty"`
	Kind      taskType           `bson:"kind"`
	createdAt time.Time          `bson:"created_at"`
	updatedAt time.Time          `bson:"updated_at"`

	targets []int `bson:"targets"`

	logs []AttemptLog `bson:"logs"`

	status         taskStatus `bson:"status"`
	scheduledAfter time.Time  `bson:"scheduled_after"`
	lastAttempt    time.Time  `bson:"last_attempt"`

	Data TaskData[T] `bson:"data"`
}

func (t *Task[T]) ID() primitive.ObjectID {
	return t.id
}
func (t *Task[T]) Type() taskType {
	return t.Kind
}

func (t *Task[T]) Targets() []int {
	return t.targets
}
func (t *Task[T]) SetTargets(targets []int) {
	t.targets = targets
}

func (t *Task[T]) Status() taskStatus {
	return t.status
}
func (t *Task[T]) SetStatus(status taskStatus) {
	t.status = status
}

func (t *Task[T]) LogAttempt(log AttemptLog) {
	t.logs = append(t.logs, log)
}

func (t *Task[T]) OnCreated() {
	t.lastAttempt = time.Now()
	t.createdAt = time.Now()
}
func (t *Task[T]) OnUpdated() {
	t.updatedAt = time.Now()
}

func (t *Task[T]) SetScheduledAfter(scheduledAfter time.Time) {
	t.scheduledAfter = scheduledAfter
}

func (t *Task[T]) RetryOnFail() bool {
	t.lastAttempt = time.Now()
	return t.Data.RetryOnFail(t)
}

func (t *Task[T]) Process() (string, error) {
	t.lastAttempt = time.Now()
	return t.Data.Process(t)
}

type AttemptLog struct {
	Targets   []int     `json:"targets" bson:"targets"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Comment   string    `json:"result" bson:"result"`
	Error     error     `json:"error" bson:"error"`
}

func NewAttemptLog(task Task[any], comment string, err error) AttemptLog {
	return AttemptLog{
		Targets:   task.Targets(),
		Timestamp: time.Now(),
		Comment:   comment,
		Error:     err,
	}
}

/*
Retrieves all target IDs from the database based on task.Type() and filter, creates a new task in queue.
  - If splitTaskFn is provided, it will split the task into subtasks.
*/
func CreateBulkTask(filter bson.M, task Task[any], splitTaskFn func(Task[any]) []Task[any]) (err error) {
	if len(task.Targets()) != 0 {
		return errors.New("target IDs already set")
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	switch task.Type() {
	case TaskUpdateClans:
		// Get all clan IDs on the realm
		result, err := database.DefaultClient.Collection(database.CollectionClans).Distinct(ctx, "externalID", filter)
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

		task.SetTargets(targetIDs)

	case TaskRecordPlayerAchievements:
		// All players on the realm
		fallthrough
	case TaskUpdateAccountWN8:
		// All players on the realm
		fallthrough
	case TaskRecordSessions:
		// Get all player IDs on the realm
		result, err := database.DefaultClient.Collection(database.CollectionAccounts).Distinct(ctx, "externalID", filter)
		if err != nil {
			return err
		}

		var targetIDs []int
		for _, id := range result {
			if idInt, ok := id.(int); ok {
				targetIDs = append(targetIDs, idInt)
			} else {
				log.Error().Msgf("invalid player ID type: %T", id)
			}
		}

		task.SetTargets(targetIDs)

	default:
		return errors.New("invalid task type")
	}

	if len(task.Targets()) == 0 {
		return errors.New("no target IDs found on realm")
	}

	if splitTaskFn != nil {
		return CreateTasks(splitTaskFn(task)...)
	}
	return CreateTasks(task)
}

func CreateTasks(tasks ...Task[any]) error {
	var writes []mongo.WriteModel
	for _, task := range tasks {
		if len(task.Targets()) == 0 {
			return errors.New("task targets not set")
		}

		task.OnCreated()
		task.SetStatus(TaskStatusScheduled)
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

func UpdateTasks(tasks ...Task[any]) error {
	var writes []mongo.WriteModel
	for _, task := range tasks {
		if task.ID().IsZero() {
			return errors.New("task ID not set")
		}
		task.OnUpdated()
		writes = append(writes, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": task.ID()}).SetUpdate(bson.M{"$set": task}))
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
func StartScheduledTasks(filter bson.M, limit int, target *[]Task[any]) error {
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

	cur, err := database.DefaultClient.Collection(database.CollectionTasks).Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	return cur.All(ctx, target)
}
