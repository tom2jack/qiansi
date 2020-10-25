/**
 * 部署服务
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/13 16:10
 */

package open

import (
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/mqttbroker"

	"github.com/gin-gonic/gin"
)

type deployApi struct{}

var Deploy = &deployApi{}

// @Summary 启动部署
// @Produce  json
// @Accept  json
// @Param key query string true "入参集合"
// @Success 200 {object} string "操作结果"
// @Router /open/DeployRun [GET]
func (r *deployApi) DeployRun(c *gin.Context) {
	openId := c.Query("Key")
	if len(openId) != 32 {
		c.String(404, "服务不存在")
		return
	}
	deployInfo, err := models.GetDeployModels().GetByOpenId(openId)
	if err != nil {
		c.String(404, "服务不存在")
		return
	}
	err = mqttbroker.SendDeployTask(deployInfo.UId, deployInfo.ID)
	if err != nil {
		c.String(500, "启动异常")
		return
	}
	c.String(200, "Successful")
}
