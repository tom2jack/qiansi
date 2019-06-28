/**
 * 部署服务
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/13 16:10
 */

package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"tools-server/common/utils"
	"tools-server/models"
	"tools-server/server/service/udp/clinet_task_loop"
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

// @Summary 设置部署应用
// @Produce  json
// @Accept  json
// @Param Id formData string true "应用ID"
// @Param AfterCommand formData string true "后置命令"
// @Param BeforeCommand formData string true "前置命令"
// @Param Branch formData string true "抓取分支"
// @Param DeployType formData int true "部署方式"
// @Param LocalPath formData string true "部署地址"
// @Param RemoteUrl formData string true "资源地址"
// @Param Title formData string true "应用名称"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/DeploySet [POST]
func DeploySet(c *gin.Context) {
	param := &models.Deploy{}
	if c.ShouldBind(param) != nil {
		utils.Show(c, -4, "入参绑定失败", nil)
		return
	}
	if param.Id == 0 {
		param.Uid = c.GetInt("UID")
		if models.ZM_Mysql.Save(param).RowsAffected > 0 {
			utils.Show(c, 1, "创建成功", param)
			return
		}
	}
	if param.Id > 0 {
		if models.ZM_Mysql.Table("deploy").Where("id=? and uid=?", param.Id, c.GetInt("UID")).Updates(param).RowsAffected > 0 {
			utils.Show(c, 1, "更新成功", param)
			return
		}
	}
	utils.Show(c, 0, "系统错误", nil)
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

// @Summary 部署应用关联服务器
// @Produce  json
// @Accept  json
// @Param deploy_id formData string true "部署应用ID"
// @Param server_id formData string true "服务器ID"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "关联成功","data": null}"
// @Router /admin/DeployRelationServer [POST]
func DeployRelationServer(c *gin.Context) {
	deploy_id, err := strconv.Atoi(c.PostForm("deploy_id"))
	if err != nil || !(deploy_id > 0) {
		utils.Show(c, -4, "部署应用ID读取错误", nil)
		return
	}
	server_id, err := strconv.Atoi(c.PostForm("server_id"))
	if err != nil || !(server_id > 0) {
		utils.Show(c, -4, "服务器ID读取错误", nil)
		return
	}
	var (
		num int
		db  *gorm.DB
	)
	db = models.ZM_Mysql.Table("server").Where("id=? and uid=?", server_id, c.GetInt("UID")).Count(&num)
	if db.Error != nil || num == 0 {
		utils.Show(c, -5, "服务器不存在", nil)
		return
	}
	db = models.ZM_Mysql.Table("deploy").Where("id=? and uid=?", server_id, c.GetInt("UID")).Count(&num)
	if db.Error != nil || num == 0 {
		utils.Show(c, -5, "部署服务不存在", nil)
		return
	}
	relation := &models.DeployServerRelation{
		ServerId: server_id,
		DeployId: deploy_id,
	}
	db = models.ZM_Mysql.Save(relation)
	if db.Error != nil || db.RowsAffected != 1 {
		utils.Show(c, -5, "关联失败", nil)
		return
	}
	utils.Show(c, 1, "关联成功", nil)
}

// @Summary 部署应用取消关联服务器
// @Produce  json
// @Accept  json
// @Param deploy_id formData string true "部署应用ID"
// @Param server_id formData string true "服务器ID"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "关联解除成功","data": null}"
// @Router /admin/DeployUnRelationServer [DELETE]
func DeployUnRelationServer(c *gin.Context) {
	deploy_id, err := strconv.Atoi(c.PostForm("deploy_id"))
	if err != nil || !(deploy_id > 0) {
		utils.Show(c, -4, "部署应用ID读取错误", nil)
		return
	}
	server_id, err := strconv.Atoi(c.PostForm("server_id"))
	if err != nil || !(server_id > 0) {
		utils.Show(c, -4, "服务器ID读取错误", nil)
		return
	}
	var (
		num int
		db  *gorm.DB
	)
	db = models.ZM_Mysql.Table("server").Where("id=? and uid=?", server_id, c.GetInt("UID")).Count(&num)
	if db.Error != nil || num == 0 {
		utils.Show(c, -5, "服务器不存在", nil)
		return
	}
	db = models.ZM_Mysql.Table("deploy").Where("id=? and uid=?", deploy_id, c.GetInt("UID")).Count(&num)
	if db.Error != nil || num == 0 {
		utils.Show(c, -5, "部署服务不存在", nil)
		return
	}
	db = models.ZM_Mysql.Delete(models.DeployServerRelation{}, "server_id=? and deploy_id=?", server_id, deploy_id)
	if db.Error != nil || db.RowsAffected != 1 {
		utils.Show(c, -5, "关联解除失败", nil)
		return
	}
	utils.Show(c, 1, "关联解除成功", nil)
}

// @Summary 启动部署
// @Produce  json
// @Accept  json
// @Param deploy_id formData string true "部署应用ID"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "启动成功","data": null}"
// @Router /admin/DeployDo [GET]
func DeployDo(c *gin.Context) {
	deploy_id, err := strconv.Atoi(c.Query("deploy_id"))
	if err != nil || !(deploy_id > 0) {
		utils.Show(c, -4, "部署应用ID读取错误", nil)
		return
	}
	var (
		db *gorm.DB
	)
	db = models.ZM_Mysql.Exec("update deploy set now_version=now_version+1 where id=? and uid=?", deploy_id, c.GetInt("UID"))
	if db.Error != nil || db.RowsAffected != 1 {
		utils.Show(c, -5, "部署服务不存在", nil)
		return
	}
	server := &[]models.Server{}
	models.ZM_Mysql.Select("id").Where("id in (select server_id from deploy_server_relation where deploy_id=?)", deploy_id).Find(server)
	for _, v := range *server {
		clinet_task_loop.Task.SET(strconv.Itoa(v.Id), "1")
	}
	utils.Show(c, 1, "启动成功", server)
}
