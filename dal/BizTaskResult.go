package dal

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var BizTaskResultFields = struct {
	Id        string
	TaskId    string
	Content   string
	CreatedAt string
}{
	"id",
	"task_id",
	"content",
	"created_at",
}

var BizTaskResultMeta = &daog.TableMeta[BizTaskResult]{
	Table: "biz_task_result",
	Columns: []string{
		"id",
		"task_id",
		"content",
		"created_at",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *BizTaskResult, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "task_id" == columnName {
			if point {
				return &ins.TaskId
			}
			return ins.TaskId
		}
		if "content" == columnName {
			if point {
				return &ins.Content
			}
			return ins.Content
		}
		if "created_at" == columnName {
			if point {
				return &ins.CreatedAt
			}
			return ins.CreatedAt
		}

		return nil
	},
}

var BizTaskResultDao daog.QuickDao[BizTaskResult] = &struct {
	daog.QuickDao[BizTaskResult]
}{
	daog.NewBaseQuickDao(BizTaskResultMeta),
}

type BizTaskResult struct {
	Id        int64
	TaskId    int64
	Content   string
	CreatedAt ttypes.NormalDatetime
}
