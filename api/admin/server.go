/**
 * 服务器模块
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/12 19:02
 */

package admin

import (
	"github.com/gin-gonic/gin"
	"tools-server/common/utils"
	"tools-server/models"
)

// @Summary 获取服务器(客户端)列表
// @Produce  json
// @Accept  json
// @Success 200 {string} json "{"code": 1,"msg": "读取成功","data": [{"ApiSecret": "123456","CreateTime": "2019-03-02T16:10:10+08:00","DeviceId": "","Domain": "127.0.0.1","Id": 1,"ServerName": "纸喵5号机","ServerRuleId": 0,"ServerStatus": 0,"Uid": 2,"UpdateTime": "2019-05-22T17:40:18+08:00"}]}"
// @Router /admin/ServerLists [post]
func ServerLists(c *gin.Context) {
	s := &[]models.Server{}
	models.ZM_Mysql.Find(s).Select("CreateTime, DeviceId, Domain, Id, ServerName, ServerRuleId, ServerStatus, Uid, UpdateTime").Where("uid = ?")
	utils.Show(c, 1, "读取成功", s)
}
