package api

import (
	"github.com/gin-gonic/gin"
	"tools-server/models"
)

func ApiIndex(c *gin.Context) {
	v := &map[string]string{
		"code": "200",
		"msg":  "成功了",
		"id":   c.Query("id"),
	}
	models.RedisSet("zhimiao_test", v, 0)
	c.JSON(200, v)
}
