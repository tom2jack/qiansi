/**
 * 部署服务
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/13 16:10
 */

package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"tools-server/common/utils"
	"tools-server/models"
)

// @Summary 获取部署服务列表
// @Produce  json
// @Accept  json
// @Success 200 {object} utils.Json "{"code": 1,"msg": "读取成功","data": [{"AfterCommand": "324545","BeforeCommand": "1232132132","Branch": "123213","CreateTime": "2019-02-28T10:24:41+08:00","DeployType": 1,"Id": 491,"LocalPath": "123213","NowVersion": 0,"RemoteUrl": "123213","Title": "491-一号机器的修改241","Uid": 2,"UpdateTime": "2019-02-28T10:25:17+08:00"}]}"
// @Router /admin/DeployLists [get]
func DeployLists(c *gin.Context) {
	d := &[]models.Deploy{}
	models.ZM_Mysql.Order("id desc").Find(d, "uid=?", c.GetInt("UID"))
	utils.Show(c, 1, "读取成功", d)
}

// @Summary 创建部署服务
// @Produce  json
// @Accept  json
// @Param deploy_id formData string true "服务器ID"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/DeployCreate [POST]
func DeployCreate(c *gin.Context) {
	// TODO
}

// @Summary 删除部署服务
// @Produce  json
// @Accept  json
// @Param deploy_id formData string true "服务器ID"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/DeployDel [DELETE]
func DeployDel(c *gin.Context) {
	deploy_id, err := strconv.Atoi(c.PostForm("deploy_id"))
	if err != nil || !(deploy_id > 0) {
		utils.Show(c, -4, "服务器ID读取错误", nil)
		return
	}

	db := models.ZM_Mysql.Delete(models.Deploy{}, "id=? and uid=?", deploy_id, c.GetInt("UID"))
	if db.Error != nil || db.RowsAffected != 1 {
		utils.Show(c, -5, "删除失败", *db)
		return
	}
	utils.Show(c, 1, "操作成功", *db)
}
