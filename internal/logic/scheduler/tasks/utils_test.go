package tasks

import (
	"testing"
)

func TestSplitAccountUpdateTasks(t *testing.T) {
	var input Task[any]
	batchSize := 100

	input.SetTargets(make([]int, 420))
	subTasks := splitTasksByTargets(batchSize)(input)
	if len(subTasks) != 5 {
		t.Errorf("expected 2 subtasks, got %d", len(subTasks))
	}
	if len(subTasks[0].Targets()) != 100 {
		t.Errorf("expected 100 targets, got %d", len(subTasks[0].Targets()))
	}
	if len(subTasks[4].Targets()) != 20 {
		t.Errorf("expected 20 targets in last task, got %d", len(subTasks[4].Targets()))
	}

	input.SetTargets(make([]int, 100))
	subTasks = splitTasksByTargets(batchSize)(input)
	if len(subTasks) != 1 {
		t.Errorf("expected 1 subtask, got %d", len(subTasks))
	}
	if len(subTasks[0].Targets()) != 100 {
		t.Errorf("expected 100 targets, got %d", len(subTasks[0].Targets()))
	}

	input.SetTargets(make([]int, 69))
	subTasks = splitTasksByTargets(batchSize)(input)
	if len(subTasks) != 1 {
		t.Errorf("expected 2 subtasks, got %d", len(subTasks))
	}
	if len(subTasks[0].Targets()) != 69 {
		t.Errorf("expected 69 targets, got %d", len(subTasks[1].Targets()))
	}
}
