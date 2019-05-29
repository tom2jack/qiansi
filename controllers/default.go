package controllers

import (
	"github.com/gin-gonic/gin"
	"tools-server/conf"
)

func Index(c *gin.Context) {
	v := &map[string]int{"111": conf.App.MustInt("app", "page_size")}
	c.JSON(200, v)
}
