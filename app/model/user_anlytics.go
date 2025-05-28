package model

type WeekOfTask struct {
	Date  string
	Count int
}

type CollectionCount struct {
	Count int
	Name  string
}

type UserAnalytics struct {
	CollectionTotal int
	TasksTotal      int
	TasksSuccess    int
	TasksFailed     int
	WeekOfTasks     []WeekOfTask
	Collections     []CollectionCount
}
