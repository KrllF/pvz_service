package entity

// ListTasksOptions опции для репозитория задач
type ListTasksOptions struct {
	Status string
	Count  int64
}

// ListTaskOption функция, чтобы устанавливать опции
type ListTaskOption func(*ListTasksOptions)

// WithTaskStatus опция со статусом задачи
func WithTaskStatus(status string) ListTaskOption {
	return func(opts *ListTasksOptions) {
		opts.Status = status
	}
}

// WithTaskCount опция с количеством задач
func WithTaskCount(count int64) ListTaskOption {
	return func(opts *ListTasksOptions) {
		opts.Count = count
	}
}
