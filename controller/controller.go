package task_controller

import (
	daogext "github.com/darwinOrg/daog-ext"
	"github.com/darwinOrg/go-biz-task/model"
	"github.com/darwinOrg/go-biz-task/permission"
	common_provider "github.com/darwinOrg/go-biz-task/provider"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/result"
	"github.com/darwinOrg/go-web/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/rolandhe/daog"
)

func Register(rg *gin.RouterGroup) {
	taskGroup := rg.Group("/public/v1/task", permission.Check)

	wrapper.Post(&wrapper.RequestHolder[task_model.GetTaskRequest, *result.Result[*task_model.CommonTaskVo]]{
		Remark:       "获取任务",
		RouterGroup:  taskGroup,
		RelativePath: "/pull",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, ctx *dgctx.DgContext, req *task_model.GetTaskRequest) *result.Result[*task_model.CommonTaskVo] {
			task, err := daogext.WriteWithResult(ctx, func(tc *daog.TransContext) (*task_model.CommonTaskVo, error) {
				return common_provider.RandomLockForProcessing(ctx, tc, req)
			})
			if err != nil {
				return result.FailByError[*task_model.CommonTaskVo](err)
			}

			return result.Success(task)
		},
	})

	wrapper.Post(&wrapper.RequestHolder[task_model.EndTaskAsFailRequest, *result.Result[*result.Void]]{
		Remark:       "失败任务",
		RouterGroup:  taskGroup,
		RelativePath: "/end-as-fail",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, ctx *dgctx.DgContext, req *task_model.EndTaskAsFailRequest) *result.Result[*result.Void] {
			err := daogext.Write(ctx, func(tc *daog.TransContext) error {
				return common_provider.EndAsFail(ctx, tc, req.Id, req.Reason)
			})
			if err != nil {
				return result.FailByError[*result.Void](err)
			}

			return result.SimpleSuccess()
		},
	})

	wrapper.Post(&wrapper.RequestHolder[task_model.EndTaskAsCanceledRequest, *result.Result[*result.Void]]{
		Remark:       "取消任务",
		RouterGroup:  taskGroup,
		RelativePath: "/end-as-canceled",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, ctx *dgctx.DgContext, req *task_model.EndTaskAsCanceledRequest) *result.Result[*result.Void] {
			err := daogext.Write(ctx, func(tc *daog.TransContext) error {
				return common_provider.EndAsCanceled(ctx, tc, req.Id, req.Reason)
			})
			if err != nil {
				return result.FailByError[*result.Void](err)
			}

			return result.SimpleSuccess()
		},
	})

}
