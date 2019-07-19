package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/zmlog"
)

func AdminIndex(c *gin.Context) {
	v := &map[string]string{
		"code": "200",
		"msg":  "成功了",
		"id":   c.Query("id"),
		"uid":  c.GetString("LOGIN-TOKEN"),
	}
	zmlog.Info("%v", v)
	c.JSON(200, v)
}
