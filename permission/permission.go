package permission

import (
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/result"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/darwinOrg/go-web/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var AuthToken string

func Check(c *gin.Context) {
	ctx := utils.GetDgContext(c)
	authToken := c.Request.Header["Authorization"]
	if len(authToken) > 0 {
		dglogger.Infof(ctx, "authToken: %s", authToken[0])
	}

	if len(authToken) == 0 || authToken[0] == "" || !strings.EqualFold(authToken[0], "Bearer "+AuthToken) {
		c.AbortWithStatusJSON(http.StatusOK, *result.FailByError[*dgerr.DgError](dgerr.NO_PERMISSION))
		return
	}

	c.Next()
}
