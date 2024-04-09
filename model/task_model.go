package task_model

type InitTaskRequest struct {
	TaskType  int    `json:"taskType" binding:"required" remark:"任务类型"`
	Channel   int    `json:"channel" remark:"渠道"`
	Content   string `json:"content" remark:"任务内容体"`
	BeginTime string `json:"beginTime" binding:"isDatetime" title:"开始时间" remark:"yyyy-MM-dd HH:mm:ss格式"`
}

type PullTaskRequest struct {
	TaskType        int    `json:"taskType" remark:"任务类型"`
	Channel         int    `json:"channel" remark:"渠道"`
	PageSize        int    `json:"pageSize" title:"页码" remark:"从多少条任务中随机选择一条，默认为1"`
	LockMilli       int64  `json:"lockMilli" title:"锁定毫秒数" remark:"默认为5000"`
	LockedBy        string `json:"lockedBy" title:"锁定者" remark:"最好是唯一标志"`
	FixedLockedBy   bool   `json:"fixedLockedBy" title:"是否固定锁定者" remark:"如果是，则后续任务失败后也只能再由当前锁定者处理"`
	FollowBeginTime bool   `json:"followBeginTime" title:"是否按照开始时间" remark:"如果是，则按照开始时间顺序拉取任务，时间未到不会返回任务"`
}

type PushTaskResultRequest struct {
	Id      int64  `json:"id" binding:"required" remark:"任务id"`
	Content string `json:"content" remark:"结果内容"`
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
