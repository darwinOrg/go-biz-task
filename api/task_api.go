package task_api

import (
	daogext "github.com/darwinOrg/daog-ext"
	"github.com/darwinOrg/go-biz-task/model"
	task_provider "github.com/darwinOrg/go-biz-task/provider"
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/result"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/darwinOrg/go-web/utils"
	"github.com/darwinOrg/go-web/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/rolandhe/daog"
	"net/http"
	"strings"
)

var authHook gin.HandlerFunc

func RegisterAuthHook(myAuthHook gin.HandlerFunc) {
	authHook = myAuthHook
}
func DefaultAuthFunc(myAuthToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if myAuthToken == "" {
			c.Next()
			return
		}

		ctx := utils.GetDgContext(c)
		authToken := c.Request.Header["Authorization"]
		if len(authToken) > 0 {
			dglogger.Infof(ctx, "authToken: %s", authToken[0])
		}

		if len(authToken) == 0 || authToken[0] == "" || !strings.EqualFold(authToken[0], "Bearer "+myAuthToken) {
			c.AbortWithStatusJSON(http.StatusOK, *result.FailByError[*dgerr.DgError](dgerr.NO_PERMISSION))
			return
		}

		c.Next()
	}
}

type PullTaskHook func(c *gin.Context, ctx *dgctx.DgContext, tc *daog.TransContext, req *task_model.PullTaskRequest) error

var pullTaskHook PullTaskHook

func RegisterPullTaskHook(hook PullTaskHook) {
	pullTaskHook = hook
}

type PushTaskResultHook func(c *gin.Context, ctx *dgctx.DgContext, tc *daog.TransContext, req *task_model.PushTaskResultRequest) error

var pushTaskResultHook PushTaskResultHook

func RegisterPushTaskResultHook(hook PushTaskResultHook) {
	pushTaskResultHook = hook
}

func RegisterApi(e *gin.Engine) {
	rg := e.Group("/public/v1/task", authHook)

	wrapper.Post(&wrapper.RequestHolder[task_model.PullTaskRequest, *result.Result[*task_model.CommonTaskVo]]{
		Remark:       "拉取任务",
		RouterGroup:  rg,
		RelativePath: "/pull",
		NonLogin:     true,
		BizHandler: func(c *gin.Context, ctx *dgctx.DgContext, req *task_model.PullTaskRequest) *result.Result[*task_model.CommonTaskVo] {
			task, err := daogext.WriteWithResult(ctx, func(tc *daog.TransContext) (*task_model.CommonTaskVo, error) {
				if pullTaskHook != nil {
					err := pullTaskHook(c, ctx, tc, req)
					if err != nil {
						return nil, err
					}
				}

				return task_provider.RandomLockForProcessing(ctx, tc, req)
			})
			if err != nil {
				return result.FailByError[*task_model.CommonTaskVo](err)
			}

			return result.Success(task)
		},
	})

	wrapper.Post(&wrapper.RequestHolder[task_model.PushTaskResultRequest, *result.Result[*result.Void]]{
		Remark:       "推送任务结果",
		RouterGroup:  rg,
		RelativePath: "push",
		NonLogin:     true,
		BizHandler: func(c *gin.Context, ctx *dgctx.DgContext, req *task_model.PushTaskResultRequest) *result.Result[*result.Void] {
			err := daogext.Write(ctx, func(tc *daog.TransContext) error {
				if pushTaskResultHook != nil {
					err := pushTaskResultHook(c, ctx, tc, req)
					if err != nil {
						return err
					}
				}

				err := task_provider.EndAsSuccess(ctx, tc, req.Id)
				if err != nil {
					return err
				}

				return task_provider.PushTaskResult(ctx, tc, req.Id, req.Content)
			})
			if err != nil {
				return result.FailByError[*result.Void](err)
			}

			return result.SimpleSuccess()
		},
	})

	wrapper.Post(&wrapper.RequestHolder[task_model.EndTaskAsFailRequest, *result.Result[*result.Void]]{
		Remark:       "失败任务",
		RouterGroup:  rg,
		RelativePath: "/end-as-fail",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, ctx *dgctx.DgContext, req *task_model.EndTaskAsFailRequest) *result.Result[*result.Void] {
			err := daogext.Write(ctx, func(tc *daog.TransContext) error {
				return task_provider.EndAsFail(ctx, tc, req.Id, req.Reason)
			})
			if err != nil {
				return result.FailByError[*result.Void](err)
			}

			return result.SimpleSuccess()
		},
	})

	wrapper.Post(&wrapper.RequestHolder[task_model.EndTaskAsCanceledRequest, *result.Result[*result.Void]]{
		Remark:       "取消任务",
		RouterGroup:  rg,
		RelativePath: "/end-as-canceled",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, ctx *dgctx.DgContext, req *task_model.EndTaskAsCanceledRequest) *result.Result[*result.Void] {
			err := daogext.Write(ctx, func(tc *daog.TransContext) error {
				return task_provider.EndAsCanceled(ctx, tc, req.Id, req.Reason)
			})
			if err != nil {
				return result.FailByError[*result.Void](err)
			}

			return result.SimpleSuccess()
		},
	})

}

func RegisterInitTaskApi(rg *gin.RouterGroup) {
	wrapper.Post(&wrapper.RequestHolder[task_model.InitTaskRequest, *result.Result[int64]]{
		Remark:       "初始化任务",
		RouterGroup:  rg,
		RelativePath: "/init",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, ctx *dgctx.DgContext, req *task_model.InitTaskRequest) *result.Result[int64] {
			taskId, err := daogext.WriteWithResult(ctx, func(tc *daog.TransContext) (int64, error) {
				return task_provider.InsertInitTask(ctx, tc, req)
			})
			if err != nil {
				return result.FailByError[int64](err)
			}

			return result.Success(taskId)
		},
	})
}
