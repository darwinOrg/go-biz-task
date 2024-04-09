package dal

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
	"time"
)

var ExtBizTaskResultDao = &extBizTaskResultDao{}

type extBizTaskResultDao struct {
}

func (d *extBizTaskResultDao) PushTaskResult(tc *daog.TransContext, taskId int64, content string) error {
	matcher := daog.NewMatcher().Eq(BizTaskResultFields.TaskId, taskId)
	rt, err := BizTaskResultDao.QueryOneMatcher(tc, matcher)
	if err != nil {
		return err
	}

	if rt == nil {
		rt = &BizTaskResult{
			TaskId:    taskId,
			Content:   content,
			CreatedAt: ttypes.NormalDatetime(time.Now()),
		}

		_, err = BizTaskResultDao.Insert(tc, rt)
		if err != nil {
			return err
		}
	} else {
		rt.Content = content
		_, err = BizTaskResultDao.Update(tc, rt)
		if err != nil {
			return err
		}
	}

	return nil
}
