package dal

import (
	task_enum "github.com/darwinOrg/go-biz-task/enum"
	"github.com/darwinOrg/go-biz-task/model"
	dgcoll "github.com/darwinOrg/go-common/collection"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
	"time"
)

var ExtBizTaskDao = &extBizTaskDao{}

type extBizTaskDao struct {
}

func (d *extBizTaskDao) InsertInitTask(tc *daog.TransContext, req *task_model.InitTaskRequest) (int64, error) {
	now := time.Now()
	task := &BizTask{
		Type:       int32(req.TaskType),
		Channel:    int32(req.Channel),
		Content:    *ttypes.FromString(req.Content),
		Status:     task_enum.TaskStatus.INIT,
		CreatedAt:  ttypes.NormalDatetime(now),
		ModifiedAt: ttypes.NormalDatetime(now),
	}

	if req.ScheduledStartAt != "" {
		scheduledStartAt, err := time.Parse(req.ScheduledStartAt, ttypes.DatetimeFormat)
		if err != nil {
			return 0, err
		}

		task.ScheduledStartAt = *ttypes.FromDatetime(scheduledStartAt)
	}

	if req.ScheduledEndAt != "" {
		scheduledEndAt, err := time.Parse(req.ScheduledEndAt, ttypes.DatetimeFormat)
		if err != nil {
			return 0, err
		}

		task.ScheduledEndAt = *ttypes.FromDatetime(scheduledEndAt)
	}

	_, err := BizTaskDao.Insert(tc, task)
	if err != nil {
		return 0, err
	}
	if task.Id == 0 {
		return 0, dgerr.SYSTEM_ERROR
	}

	return task.Id, nil
}

func (d *extBizTaskDao) FindToHandleTasks(tc *daog.TransContext, req *task_model.PullTaskRequest) ([]*BizTask, error) {
	matcher := daog.NewMatcher().
		In(BizTaskFields.Status, daog.ConvertToAnySlice(task_enum.ToHandleStatuses))
	if req.TaskType > 0 {
		matcher.Eq(BizTaskFields.Type, req.TaskType)
	}
	if req.Channel > 0 {
		matcher.Eq(BizTaskFields.Channel, req.Channel)
	}
	if req.FixedLockedBy {
		matcher.Add(daog.NewOrMatcher().Null(BizTaskFields.LockedBy, false).Eq(BizTaskFields.LockedBy, req.LockedBy))
	}

	var order *daog.Order
	if req.FollowScheduledTime {
		now := time.Now()
		matcher.Lte(BizTaskFields.ScheduledStartAt, now)
		matcher.Add(daog.NewOrMatcher().Null(BizTaskFields.ScheduledEndAt, false).Gte(BizTaskFields.ScheduledEndAt, now))
		order = daog.NewOrder(BizTaskFields.ScheduledStartAt)
	} else {
		order = daog.NewOrder(BizTaskFields.CreatedAt)
	}

	if req.PageSize == 0 {
		req.PageSize = 1
	}
	pager := daog.NewPager(req.PageSize, 1)

	return BizTaskDao.QueryPageListMatcher(tc, matcher, pager, order)
}

func (d *extBizTaskDao) RandomLockForProcessing(tc *daog.TransContext, req *task_model.PullTaskRequest) (*BizTask, error) {
	tasks, err := d.FindToHandleTasks(tc, req)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	if len(tasks) > 1 {
		dgcoll.Shuffle(tasks)
	}

	for _, task := range tasks {
		if req.LockMilli == 0 {
			req.LockMilli = 5000
		}
		ok, err := d.LockForProcessing(tc, task.Id, req.LockMilli, req.LockedBy)
		if err != nil {
			return nil, err
		}
		if ok {
			return task, nil
		}
	}

	return nil, nil
}

func (d *extBizTaskDao) LockForProcessing(tc *daog.TransContext, taskId int64, lockMilli int64, lockedBy string) (bool, error) {
	now := time.Now()
	lockUntil := now.Add(time.Millisecond * time.Duration(lockMilli))

	modifier := daog.NewModifier().
		Add(BizTaskFields.Status, task_enum.TaskStatus.PROCESSING).
		Add(BizTaskFields.LockedAt, ttypes.NormalDatetime(now)).
		Add(BizTaskFields.LockUntil, lockUntil).
		Add(BizTaskFields.LockedBy, lockedBy).
		Add(BizTaskFields.ModifiedAt, ttypes.NormalDatetime(now))
	matcher := daog.NewMatcher().Eq(BizTaskFields.Id, taskId).In(BizTaskFields.Status, daog.ConvertToAnySlice(task_enum.ToHandleStatuses)).
		Add(daog.NewOrMatcher().Null(BizTaskFields.LockUntil, false).Lte(BizTaskFields.LockUntil, now))

	count, err := BizTaskDao.UpdateByModifier(tc, modifier, matcher)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *extBizTaskDao) EndAsSuccess(tc *daog.TransContext, taskId int64) error {
	now := time.Now()
	modifier := daog.NewModifier().
		Add(BizTaskFields.Status, task_enum.TaskStatus.SUCCESS).
		Add(BizTaskFields.EndAt, ttypes.NormalDatetime(now)).
		SelfAdd(BizTaskFields.ProcessedCount, 1).
		Add(BizTaskFields.ModifiedAt, ttypes.NormalDatetime(now))
	_, err := BizTaskDao.UpdateById(tc, modifier, taskId)
	return err
}

func (d *extBizTaskDao) EndAsFail(tc *daog.TransContext, taskId int64, reason string) error {
	now := time.Now()
	modifier := daog.NewModifier().
		Add(BizTaskFields.Status, task_enum.TaskStatus.FAIL).
		Add(BizTaskFields.EndAt, ttypes.NormalDatetime(now)).
		Add(BizTaskFields.Reason, reason).
		SelfAdd(BizTaskFields.ProcessedCount, 1).
		Add(BizTaskFields.ModifiedAt, ttypes.NormalDatetime(now))
	_, err := BizTaskDao.UpdateById(tc, modifier, taskId)
	return err
}

func (d *extBizTaskDao) EndAsCanceled(tc *daog.TransContext, taskId int64, reason string) error {
	now := time.Now()
	modifier := daog.NewModifier().
		Add(BizTaskFields.Status, task_enum.TaskStatus.CANCELED).
		Add(BizTaskFields.EndAt, ttypes.NormalDatetime(now)).
		Add(BizTaskFields.Reason, reason).
		SelfAdd(BizTaskFields.ProcessedCount, 1).
		Add(BizTaskFields.ModifiedAt, ttypes.NormalDatetime(now))
	_, err := BizTaskDao.UpdateById(tc, modifier, taskId)
	return err
}

func (d *extBizTaskDao) UpdateContent(tc *daog.TransContext, taskId int64, content string) error {
	now := time.Now()
	modifier := daog.NewModifier().
		Add(BizTaskFields.Content, content).
		Add(BizTaskFields.ModifiedAt, ttypes.NormalDatetime(now))
	_, err := BizTaskDao.UpdateById(tc, modifier, taskId)
	return err
}

func (d *extBizTaskDao) ReInit(tc *daog.TransContext, taskId int64) error {
	now := time.Now()
	modifier := daog.NewModifier().
		Add(BizTaskFields.Status, task_enum.TaskStatus.INIT).
		Add(BizTaskFields.Reason, nil).
		Add(BizTaskFields.EndAt, nil).
		Add(BizTaskFields.LockedBy, nil).
		Add(BizTaskFields.LockedAt, nil).
		Add(BizTaskFields.LockUntil, nil).
		Add(BizTaskFields.ModifiedAt, ttypes.NormalDatetime(now))
	_, err := BizTaskDao.UpdateById(tc, modifier, taskId)
	return err
}

func (d *extBizTaskDao) GetByIdsAndLockedBy(tc *daog.TransContext, taskIds []int64, lockedBy string) ([]*BizTask, error) {
	return BizTaskDao.QueryListMatcher(tc, daog.NewMatcher().
		In(BizTaskFields.Id, daog.ConvertToAnySlice(taskIds)).
		Eq(BizTaskFields.LockedBy, lockedBy))
}

func (d *extBizTaskDao) ReInitTimeoutProcessingTasks(tc *daog.TransContext, taskType int, timeoutMinutes int) (int64, error) {
	now := time.Now()
	modifier := daog.NewModifier().
		Add(BizTaskFields.Status, task_enum.TaskStatus.INIT).
		Add(BizTaskFields.Reason, nil).
		Add(BizTaskFields.EndAt, nil).
		Add(BizTaskFields.LockedBy, nil).
		Add(BizTaskFields.LockedAt, nil).
		Add(BizTaskFields.LockUntil, nil).
		Add(BizTaskFields.ModifiedAt, ttypes.NormalDatetime(now))
	matcher := daog.NewMatcher().
		Eq(BizTaskFields.Type, taskType).
		Eq(BizTaskFields.Status, task_enum.TaskStatus.PROCESSING).
		Lt(BizTaskFields.LockedAt, ttypes.NormalDatetime(now.Add(time.Minute*time.Duration(-timeoutMinutes))))
	return BizTaskDao.UpdateByModifier(tc, modifier, matcher)
}
