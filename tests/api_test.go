package tests

import (
	task_api "github.com/darwinOrg/go-biz-task/api"
	"github.com/darwinOrg/go-biz-task/model"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/darwinOrg/go-web/wrapper"
	"testing"
)

func TestBizTaskApi(t *testing.T) {
	e := wrapper.DefaultEngine()
	task_api.RegisterApi(e)
	task_api.RegisterPushTaskResultHook(func(ctx *dgctx.DgContext, req *task_model.PushTaskResultRequest) error {
		dglogger.Infof(ctx, "req: %v", req)
		return nil
	})
	_ = e.Run(":8080")
}
