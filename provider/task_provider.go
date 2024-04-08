package task_provider

import (
	"github.com/darwinOrg/go-biz-task/dal"
	task_model "github.com/darwinOrg/go-biz-task/model"
	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
)

func InsertInitTask(ctx *dgctx.DgContext, tc *daog.TransContext, taskType string, channel string, content string) (int64, error) {
	taskId, err := dal.ExtBizTaskDao.InsertInitTask(tc, taskType, channel, content)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.InsertInitTask error: %v", err)
		return 0, err
	}
	return taskId, nil
}

func FindToHandleTasks(ctx *dgctx.DgContext, tc *daog.TransContext, req *task_model.GetTaskRequest) ([]*task_model.CommonTaskVo, error) {
	tasks, err := dal.ExtBizTaskDao.FindToHandleTasks(tc, req)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.FindToHandleTasks error: %v", err)
		return nil, err
	}
	if len(tasks) == 0 {
		return []*task_model.CommonTaskVo{}, nil
	}
	return dgcoll.MapToList(tasks, convertTaskVo), nil
}

func RandomLockForProcessing(ctx *dgctx.DgContext, tc *daog.TransContext, req *task_model.GetTaskRequest) (*task_model.CommonTaskVo, error) {
	task, err := dal.ExtBizTaskDao.RandomLockForProcessing(tc, req)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.RandomLockForProcessing error: %v", err)
		return nil, err
	}
	if task == nil {
		return nil, nil
	}
	return convertTaskVo(task), nil
}

func LockForProcessing(tc *daog.TransContext, taskId int64, lockMilli int64, lockedBy string) (bool, error) {
	return dal.ExtBizTaskDao.LockForProcessing(tc, taskId, lockMilli, lockedBy)
}

func EndAsSuccess(ctx *dgctx.DgContext, tc *daog.TransContext, taskId int64) error {
	err := dal.ExtBizTaskDao.EndAsSuccess(tc, taskId)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.EndAsSuccess error: %v", err)
	}
	return err
}

func EndAsFail(ctx *dgctx.DgContext, tc *daog.TransContext, taskId int64, reason string) error {
	err := dal.ExtBizTaskDao.EndAsFail(tc, taskId, reason)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.EndAsFail error: %v", err)
	}
	return err
}

func EndAsCanceled(ctx *dgctx.DgContext, tc *daog.TransContext, taskId int64, reason string) error {
	err := dal.ExtBizTaskDao.EndAsCanceled(tc, taskId, reason)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.EndAsCanceled error: %v", err)
	}
	return err
}

func UpdateContent(ctx *dgctx.DgContext, tc *daog.TransContext, taskId int64, content string) error {
	err := dal.ExtBizTaskDao.UpdateContent(tc, taskId, content)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.UpdateContent error: %v", err)
	}
	return err
}

func ReInit(ctx *dgctx.DgContext, tc *daog.TransContext, taskId int64) error {
	err := dal.ExtBizTaskDao.ReInit(tc, taskId)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.ReInit error: %v", err)
	}
	return err
}

func GetByIdsAndLockedBy(ctx *dgctx.DgContext, tc *daog.TransContext, taskIds []int64, lockedBy string) ([]*task_model.CommonTaskVo, error) {
	tasks, err := dal.ExtBizTaskDao.GetByIdsAndLockedBy(tc, taskIds, lockedBy)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.GetByIdsAndLockedBy error: %v", err)
		return nil, err
	}

	return dgcoll.MapToList(tasks, convertTaskVo), nil
}

func ReInitTimeoutProcessingTasks(ctx *dgctx.DgContext, tc *daog.TransContext, taskType string, timeoutMinutes int) (int64, error) {
	count, err := dal.ExtBizTaskDao.ReInitTimeoutProcessingTasks(tc, taskType, timeoutMinutes)
	if err != nil {
		dglogger.Errorf(ctx, "ExtBizTaskDao.ReInitTimeoutProcessingTasks error: %v", err)
		return 0, err
	}
	if count > 0 {
		dglogger.Infof(ctx, "ReInitTimeoutProcessingTasks success, taskType: %s, count: %d", taskType, count)
	}

	return count, nil
}

func convertTaskVo(task *dal.BizTask) *task_model.CommonTaskVo {
	return &task_model.CommonTaskVo{
		Id:             task.Id,
		Content:        task.Content.StringNilAsEmpty(),
		ProcessedCount: task.ProcessedCount,
	}
}
