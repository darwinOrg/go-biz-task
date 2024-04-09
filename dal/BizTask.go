package dal

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var BizTaskFields = struct {
	Id             string
	Type           string
	Channel        string
	Content        string
	Status         string
	Reason         string
	StartAt        string
	EndAt          string
	Priority       string
	ProcessedCount string
	LockedAt       string
	LockUntil      string
	LockedBy       string
	CreatedAt      string
	ModifiedAt     string
}{
	"id",
	"type",
	"channel",
	"content",
	"status",
	"reason",
	"start_at",
	"end_at",
	"priority",
	"processed_count",
	"locked_at",
	"lock_until",
	"locked_by",
	"created_at",
	"modified_at",
}

var BizTaskMeta = &daog.TableMeta[BizTask]{
	Table: "biz_task",
	Columns: []string{
		"id",
		"type",
		"channel",
		"content",
		"status",
		"reason",
		"start_at",
		"end_at",
		"priority",
		"processed_count",
		"locked_at",
		"lock_until",
		"locked_by",
		"created_at",
		"modified_at",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *BizTask, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "type" == columnName {
			if point {
				return &ins.Type
			}
			return ins.Type
		}
		if "channel" == columnName {
			if point {
				return &ins.Channel
			}
			return ins.Channel
		}
		if "content" == columnName {
			if point {
				return &ins.Content
			}
			return ins.Content
		}
		if "status" == columnName {
			if point {
				return &ins.Status
			}
			return ins.Status
		}
		if "reason" == columnName {
			if point {
				return &ins.Reason
			}
			return ins.Reason
		}
		if "start_at" == columnName {
			if point {
				return &ins.StartAt
			}
			return ins.StartAt
		}
		if "end_at" == columnName {
			if point {
				return &ins.EndAt
			}
			return ins.EndAt
		}
		if "priority" == columnName {
			if point {
				return &ins.Priority
			}
			return ins.Priority
		}
		if "processed_count" == columnName {
			if point {
				return &ins.ProcessedCount
			}
			return ins.ProcessedCount
		}
		if "locked_at" == columnName {
			if point {
				return &ins.LockedAt
			}
			return ins.LockedAt
		}
		if "lock_until" == columnName {
			if point {
				return &ins.LockUntil
			}
			return ins.LockUntil
		}
		if "locked_by" == columnName {
			if point {
				return &ins.LockedBy
			}
			return ins.LockedBy
		}
		if "created_at" == columnName {
			if point {
				return &ins.CreatedAt
			}
			return ins.CreatedAt
		}
		if "modified_at" == columnName {
			if point {
				return &ins.ModifiedAt
			}
			return ins.ModifiedAt
		}

		return nil
	},
}

var BizTaskDao daog.QuickDao[BizTask] = &struct {
	daog.QuickDao[BizTask]
}{
	daog.NewBaseQuickDao(BizTaskMeta),
}

type BizTask struct {
	Id             int64
	Type           string
	Channel        string
	Content        ttypes.NilableString
	Status         int8
	Reason         ttypes.NilableString
	StartAt        ttypes.NilableDatetime
	EndAt          ttypes.NilableDatetime
	Priority       int32
	ProcessedCount int32
	LockedAt       ttypes.NilableDatetime
	LockUntil      ttypes.NilableDatetime
	LockedBy       ttypes.NilableString
	CreatedAt      ttypes.NormalDatetime
	ModifiedAt     ttypes.NormalDatetime
}
