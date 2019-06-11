package admin

import (
	"github.com/gin-gonic/gin"
)

func AdminIndex(c *gin.Context) {
	v := &map[string]string{
		"code": "200",
		"msg":  "成功了",
		"id":   c.Query("id"),
		"uid":  c.GetString("LOGIN-TOKEN"),
	}
	c.JSON(200, v)
}
