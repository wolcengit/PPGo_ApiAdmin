/**********************************************
** @Des: 管理员
** @Author: haodaquan
** @Date:   2017-09-16 14:17:37
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:14:07
***********************************************/
package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/george518/PPGo_ApiAdmin/libs"
	"github.com/george518/PPGo_ApiAdmin/models"
)

type AdminController struct {
	BaseController
}

func (self *AdminController) List() {
	self.Data["pageTitle"] = "管理员管理"
	self.display()
}

func (self *AdminController) Add() {
	self.Data["pageTitle"] = "新增管理员"

	// 角色
	list, _ := models.RoleGetListForSelect("")
	self.Data["role"] = list

	self.display()
}

func (self *AdminController) Edit() {
	self.Data["pageTitle"] = "编辑管理员"

	id, _ := self.GetInt("id", 0)
	Admin, _ := models.AdminGetById(id)
	row := make(map[string]interface{})
	row["id"] = Admin.Id
	row["login_name"] = Admin.LoginName
	row["real_name"] = Admin.RealName
	row["phone"] = Admin.Phone
	row["email"] = Admin.Email
	row["role_ids"] = Admin.RoleIds
	self.Data["admin"] = row

	list, _ := models.RoleGetListForSelect(Admin.RoleIds)
	self.Data["role"] = list
	self.display()
}

func (self *AdminController) Detail() {
	self.Data["pageTitle"] = "查看管理员"

	id, _ := self.GetInt("id", 0)
	Admin, _ := models.AdminGetById(id)
	row := make(map[string]interface{})
	row["id"] = Admin.Id
	row["login_name"] = Admin.LoginName
	row["real_name"] = Admin.RealName
	row["phone"] = Admin.Phone
	row["email"] = Admin.Email
	row["appkey"] = Admin.Appkey
	row["role_ids"] = Admin.RoleIds
	self.Data["admin"] = row

	list, _ := models.RoleGetListForSelect(Admin.RoleIds)
	self.Data["role"] = list
	self.display()
}

func (self *AdminController) AjaxSave() {
	Admin_id, _ := self.GetInt("id")
	if Admin_id == 0 {
		Admin := new(models.Admin)
		Admin.LoginName = strings.TrimSpace(self.GetString("login_name"))
		Admin.RealName = strings.TrimSpace(self.GetString("real_name"))
		Admin.Phone = strings.TrimSpace(self.GetString("phone"))
		Admin.Email = strings.TrimSpace(self.GetString("email"))
		Admin.RoleIds = strings.TrimSpace(self.GetString("roleids"))
		Admin.UpdateTime = time.Now().Unix()
		Admin.UpdateId = self.userId
		Admin.Status = 1

		// 检查登录名是否已经存在
		_, err := models.AdminGetByName(Admin.LoginName)

		if err == nil {
			self.ajaxMsg("登录名已经存在", MSG_ERR)
		}
		//新增
		pwd, salt := libs.Password(4, "")
		Admin.Password = pwd
		Admin.Salt = salt
		Admin.Appkey = libs.GetRandomString(16)
		Admin.CreateTime = time.Now().Unix()
		Admin.CreateId = self.userId
		if _, err := models.AdminAdd(Admin); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg("", MSG_OK)
	}

	Admin, _ := models.AdminGetById(Admin_id)
	//修改
	Admin.Id = Admin_id
	Admin.UpdateTime = time.Now().Unix()
	Admin.UpdateId = self.userId
	Admin.LoginName = strings.TrimSpace(self.GetString("login_name"))
	Admin.RealName = strings.TrimSpace(self.GetString("real_name"))
	Admin.Phone = strings.TrimSpace(self.GetString("phone"))
	Admin.Email = strings.TrimSpace(self.GetString("email"))
	Admin.RoleIds = strings.TrimSpace(self.GetString("roleids"))
	Admin.UpdateTime = time.Now().Unix()
	Admin.UpdateId = self.userId
	Admin.Status = 1

	resetPwd, _ := self.GetInt("reset_pwd")
	if resetPwd == 1 {
		pwd, salt := libs.Password(4, "")
		Admin.Password = pwd
		Admin.Salt = salt
	}
	resetAppkey, _ := self.GetInt("reset_appkey")
	if resetAppkey == 1 {
		Admin.Appkey = libs.GetRandomString(16)
	}
	if err := Admin.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg(strconv.Itoa(resetPwd), MSG_OK)
}

func (self *AdminController) AjaxDel() {

	Admin_id, _ := self.GetInt("id")
	status := strings.TrimSpace(self.GetString("status"))
	if Admin_id == 1 {
		self.ajaxMsg("超级管理员不允许操作", MSG_ERR)
	}

	Admin_status := 0
	if status == "enable" {
		Admin_status = 1
	}
	Admin, _ := models.AdminGetById(Admin_id)
	Admin.UpdateTime = time.Now().Unix()
	Admin.Status = Admin_status
	Admin.Id = Admin_id

	if err := Admin.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("操作成功", MSG_OK)
}

func (self *AdminController) AjaxTable() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	realName := strings.TrimSpace(self.GetString("realName"))

	StatusText := make(map[int]string)
	StatusText[0] = "<font color='red'>禁用</font>"
	StatusText[1] = "正常"

	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	//
	if realName != "" {
		filters = append(filters, "real_name__icontains", realName)
	}
	result, count := models.AdminGetList(page, self.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["login_name"] = v.LoginName
		row["real_name"] = v.RealName
		row["phone"] = v.Phone
		row["email"] = v.Email
		row["role_ids"] = v.RoleIds
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		row["status"] = v.Status
		row["status_text"] = StatusText[v.Status]
		list[k] = row
	}
	self.ajaxList("成功", MSG_OK, count, list)
}
