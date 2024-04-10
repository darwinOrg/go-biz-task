package tests

import (
	task_api "github.com/darwinOrg/go-biz-task/api"
	"github.com/darwinOrg/go-biz-task/model"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/darwinOrg/go-web/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/rolandhe/daog"
	"os"
	"testing"
)

func TestBizTaskApi(t *testing.T) {
	e := wrapper.DefaultEngine()
	task_api.RegisterApi(e)
	task_api.RegisterAuthHook(task_api.DefaultAuthFunc(os.Getenv("TASK_AUTH_TOKEN")))
	task_api.RegisterPullTaskHook(func(c *gin.Context, ctx *dgctx.DgContext, _ *daog.TransContext, req *task_model.PullTaskRequest) error {
		dglogger.Infof(ctx, "pull task req: %v", req)
		return nil
	})
	task_api.RegisterPushTaskResultHook(func(c *gin.Context, ctx *dgctx.DgContext, _ *daog.TransContext, req *task_model.PushTaskResultRequest) error {
		dglogger.Infof(ctx, "push task result req: %v", req)
		return nil
	})
	_ = e.Run(":8080")
}
