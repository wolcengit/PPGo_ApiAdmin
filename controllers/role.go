/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-14 14:23:32
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:31:13
***********************************************/
package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/george518/PPGo_ApiAdmin/models"
)

type RoleController struct {
	BaseController
}

func (self *RoleController) List() {
	self.Data["pageTitle"] = "角色管理"
	self.display()
}

func (self *RoleController) Add() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "新增角色"
	self.display()
}
func (self *RoleController) Edit() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "编辑角色"

	id, _ := self.GetInt("id", 0)
	role, _ := models.RoleGetById(id)
	row := make(map[string]interface{})
	row["id"] = role.Id
	row["role_name"] = role.RoleName
	row["detail"] = role.Detail
	self.Data["role"] = row

	//获取选择的树节点
	roleAuth, _ := models.RoleAuthGetById(id)
	authId := make([]int, 0)
	for _, v := range roleAuth {
		authId = append(authId, v.AuthId)
	}
	self.Data["auth"] = authId
	fmt.Println(authId)
	self.display()
}
func (self *RoleController) Detail() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "查看角色"

	id, _ := self.GetInt("id", 0)
	role, _ := models.RoleGetById(id)
	row := make(map[string]interface{})
	row["id"] = role.Id
	row["role_name"] = role.RoleName
	row["detail"] = role.Detail
	self.Data["role"] = row

	//获取选择的树节点
	roleAuth, _ := models.RoleAuthGetById(id)
	authId := make([]int, 0)
	for _, v := range roleAuth {
		authId = append(authId, v.AuthId)
	}
	self.Data["auth"] = authId
	fmt.Println(authId)
	self.display()
}

func (self *RoleController) AjaxSave() {
	auths := strings.TrimSpace(self.GetString("nodes_data"))
	role_id, _ := self.GetInt("id")
	if role_id == 0 {
		//新增
		role := new(models.Role)
		role.RoleName = strings.TrimSpace(self.GetString("role_name"))
		role.Detail = strings.TrimSpace(self.GetString("detail"))
		role.CreateTime = time.Now().Unix()
		role.UpdateTime = time.Now().Unix()
		role.Status = 1
		role.CreateId = self.userId
		role.UpdateId = self.userId
		if id, err := models.RoleAdd(role); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		} else {
			models.RoleAuthInsertMult(id, auths)
		}
		self.ajaxMsg("", MSG_OK)
	}
	//修改
	role, _ := models.RoleGetById(role_id)
	if role.Status == 2 {
		self.ajaxMsg("系统角色不可以编辑", MSG_ERR)
	}
	role.RoleName = strings.TrimSpace(self.GetString("role_name"))
	role.Detail = strings.TrimSpace(self.GetString("detail"))
	role.UpdateTime = time.Now().Unix()
	role.UpdateTime = time.Now().Unix()
	role.UpdateId = self.userId
	if err := role.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	} else {
		// 删除该角色权限
		models.RoleAuthInsertMult(int64(role.Id), auths)
	}
	self.ajaxMsg("", MSG_OK)
}

func (self *RoleController) AjaxDel() {

	role_id, _ := self.GetInt("id")
	role, _ := models.RoleGetById(role_id)
	if role.Status == 2 {
		self.ajaxMsg("系统角色不可以删除", MSG_ERR)
	}
	role.Status = 0
	role.Id = role_id
	role.UpdateTime = time.Now().Unix()

	if err := role.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}

func (self *RoleController) AjaxTable() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	roleName := strings.TrimSpace(self.GetString("roleName"))
	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status__gte", 1)
	if roleName != "" {
		filters = append(filters, "role_name__icontains", roleName)
	}
	result, count := models.RoleGetList(page, self.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["status"] = v.Status
		row["role_name"] = v.RoleName
		row["detail"] = v.Detail
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		list[k] = row
	}
	self.ajaxList("成功", MSG_OK, count, list)
}
