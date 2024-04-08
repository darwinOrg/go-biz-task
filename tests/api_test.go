package tests

import (
	task_api "github.com/darwinOrg/go-biz-task/api"
	"github.com/darwinOrg/go-web/wrapper"
	"testing"
)

func TestBizTaskApi(t *testing.T) {
	e := wrapper.DefaultEngine()
	task_api.Register(e)
	_ = e.Run(":8080")
}
