/**********************************************
** @Des: xop module
** @Author: wolcen
***********************************************/
package controllers

import (
	"strings"
	"time"

	"github.com/george518/PPGo_ApiAdmin/models"
)

type XopModuleController struct {
	BaseController
}

func (self *XopModuleController) List() {
	self.Data["pageTitle"] = "XOP模块"
	self.display()
}

func (self *XopModuleController) AjaxTable() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	searchName := strings.TrimSpace(self.GetString("searchName"))

	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "prod_id", self.prodId)
	filters = append(filters, "status", 1)
	if searchName != "" {
		filters = append(filters, "name__icontains", searchName)
	}
	list, count := models.XopModuleGetListForBrowse(page, self.pageSize, filters...)
	self.ajaxList("成功", MSG_OK, count, list)
}

func (self *XopModuleController) Add() {
	self.Data["pageTitle"] = "新增模块"
	self.display()
}

func (self *XopModuleController) Edit() {
	self.Data["pageTitle"] = "编辑模块"

	id, _ := self.GetInt("id", 0)
	v, err := models.XopModuleGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = v.Id
	row["code"] = v.Code
	row["name"] = v.Name
	row["xop_name"] = v.XopName
	row["detail"] = v.Detail
	self.Data["entity"] = row

	self.display()
}

func (self *XopModuleController) AjaxSave() {
	pkid, _ := self.GetInt("id")
	var entity *models.XopModule
	if pkid == 0 {
		entity = new(models.XopModule)
	} else {
		entity, _ = models.XopModuleGetById(pkid)
	}
	entity.Code = strings.TrimSpace(self.GetString("code"))
	entity.Name = strings.TrimSpace(self.GetString("name"))
	entity.XopName = strings.TrimSpace(self.GetString("xop_name"))
	entity.Detail = strings.TrimSpace(self.GetString("detail"))
	entity.Status = 1
	refid, _ := self.GetInt("ref_id")
	if refid > 0 {
		entity.RefId = refid
	}

	// 检查是否已经存在
	entity1, err := models.XopModuleGetByCode(self.prodId, entity.Code)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("模块编码已经存在", MSG_ERR)
	}

	entity1, err = models.XopModuleGetByName(self.prodId, entity.Name)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("模块名称已经存在", MSG_ERR)
	}

	if pkid == 0 {
		entity.ProdId = self.prodId
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if _, err := models.XopModuleAdd(entity); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	} else {
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if err := entity.Update(); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	}
	self.ajaxMsg("", MSG_OK)
}

func (self *XopModuleController) AjaxDel() {

	pkid, _ := self.GetInt("id")
	_, err := models.XopModuleDelete(pkid)
	if err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}
