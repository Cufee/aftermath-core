package tasks

func splitTasksByTargets(batchSize int) func(task Task) []Task {
	return func(task Task) []Task {
		if len(task.Targets) <= batchSize {
			return []Task{task}
		}

		var tasks []Task
		subTasks := len(task.Targets) / batchSize

		for i := 0; i <= subTasks; i++ {
			subTask := task
			if len(task.Targets) > batchSize*(i+1) {
				subTask.Targets = (task.Targets[batchSize*i : batchSize*(i+1)])
			} else {
				subTask.Targets = (task.Targets[batchSize*i:])
			}
			tasks = append(tasks, subTask)
		}

		return tasks
	}
}
