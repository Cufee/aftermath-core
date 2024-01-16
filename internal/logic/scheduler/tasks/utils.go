package tasks

func splitTasksByTargets(batchSize int) func(task Task[any]) []Task[any] {
	return func(task Task[any]) []Task[any] {
		if len(task.Targets()) <= batchSize {
			return []Task[any]{task}
		}

		var tasks []Task[any]
		subTasks := len(task.Targets()) / batchSize

		for i := 0; i <= subTasks; i++ {
			subTask := task
			if len(task.Targets()) > batchSize*(i+1) {
				subTask.SetTargets(task.Targets()[batchSize*i : batchSize*(i+1)])
			} else {
				subTask.SetTargets(task.Targets()[batchSize*i:])
			}
			tasks = append(tasks, subTask)
		}

		return tasks
	}
}
