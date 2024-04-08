package task_model

type GetTaskRequest struct {
	TaskType      string `json:"taskType" binding:"required" remark:"任务类型"`
	Channel       string `json:"channel" remark:"渠道"`
	PageSize      int    `json:"pageSize" binding:"required"`
	LockMilli     int64  `json:"lockMilli" binding:"required" remark:"锁定毫秒数"`
	LockedBy      string `json:"lockedBy" binding:"required" remark:"锁定者"`
	FixedLockedBy bool   `json:"fixedLockedBy" remark:"是否固定锁定者"`
}

type EndTaskAsFailRequest struct {
	Id     int64  `json:"id" binding:"required" remark:"任务id"`
	Reason string `json:"reason" binding:"required" remark:"失败原因"`
}

type EndTaskAsCanceledRequest struct {
	Id     int64  `json:"id" binding:"required" remark:"任务id"`
	Reason string `json:"reason" binding:"required" remark:"取消原因"`
}

type CommonTaskVo struct {
	Id             int64  `json:"id" remark:"任务id"`
	Content        string `json:"content" remark:"任务内容体"`
	ProcessedCount int32  `json:"processedCount" remark:"任务已处理次数"`
}
