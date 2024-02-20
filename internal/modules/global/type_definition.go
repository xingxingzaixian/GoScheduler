package global

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}
