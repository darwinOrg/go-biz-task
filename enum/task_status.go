package task_enum

var (
	TaskStatus = struct {
		INIT       int8
		PROCESSING int8
		SUCCESS    int8
		FAIL       int8
		CANCELED   int8
	}{
		0,
		1,
		2,
		3,
		4,
	}
)

var ToHandleStatuses = []int8{TaskStatus.INIT, TaskStatus.FAIL}
