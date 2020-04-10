/**********************************************
** @Des: xop group
** @Author: wolcen
***********************************************/
package controllers

import (
	"strings"
	"time"

	"github.com/george518/PPGo_ApiAdmin/models"
)

type XopGroupController struct {
	BaseController
}

func (self *XopGroupController) List() {
	self.Data["pageTitle"] = "XOP分组"
	list, _ := models.XopModuleGetListForSelect(self.prodId)
	self.Data["Modules"] = list
	self.display()
}

func (self *XopGroupController) AjaxTable() {
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
	modId, _ := self.GetInt("modId")
	catId, _ := self.GetInt("catId")

	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "prod_id", self.prodId)
	filters = append(filters, "status", 1)
	if searchName != "" {
		filters = append(filters, "name__icontains", searchName)
	}
	if modId > 0 {
		filters = append(filters, "mod_id", modId)
	}
	if catId > 0 {
		filters = append(filters, "cat _id", catId)
	}
	list, count := models.XopGroupGetListForBrowse(page, self.pageSize, filters...)
	self.ajaxList("成功", MSG_OK, count, list)
}

func (self *XopGroupController) Add() {
	self.Data["pageTitle"] = "新增分组"
	list, _ := models.XopModuleGetListForSelect(self.prodId)
	self.Data["Modules"] = list
	self.display()
}

func (self *XopGroupController) Edit() {
	self.Data["pageTitle"] = "编辑分组"

	id, _ := self.GetInt("id", 0)
	v, err := models.XopGroupGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = v.Id
	row["code"] = v.Code
	row["name"] = v.Name
	row["detail"] = v.Detail
	row["mod_id"] = int(v.ModId)
	row["cat_id"] = int(v.CatId)
	self.Data["entity"] = row

	list, _ := models.XopModuleGetListForSelect(self.prodId)
	self.Data["Modules"] = list
	list1, _ := models.XopCategoryGetListForSelect(v.ModId)
	self.Data["Categorys"] = list1
	self.display()
}

func (self *XopGroupController) AjaxSave() {
	pkid, _ := self.GetInt("id")
	var entity *models.XopGroup
	if pkid == 0 {
		entity = new(models.XopGroup)
		entity.ProdId = self.prodId
	} else {
		entity, _ = models.XopGroupGetById(pkid)
	}

	entity.Code = strings.TrimSpace(self.GetString("code"))
	entity.Name = strings.TrimSpace(self.GetString("name"))
	entity.Detail = strings.TrimSpace(self.GetString("detail"))
	entity.Status = 1
	entity.ModId, _ = self.GetInt("mod_id")
	entity.CatId, _ = self.GetInt("cat_id")
	refid, _ := self.GetInt("ref_id")
	if(refid > 0){
		entity.RefId = refid
	}

	// 检查是否已经存在
	entity1, err := models.XopGroupGetByCode(self.prodId, entity.Code)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("分组编码已经存在", MSG_ERR)
	}

	entity1, err = models.XopGroupGetByName(self.prodId, entity.Name)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("分组名称已经存在", MSG_ERR)
	}

	if pkid == 0 {
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if _, err := models.XopGroupAdd(entity); err != nil {
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

func (self *XopGroupController) AjaxDel() {

	pkid, _ := self.GetInt("id")
	_, err := models.XopGroupDelete(pkid)
	if err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}

func (self *XopGroupController) AjaxList() {
	pkid, _ := self.GetInt("id")
	list, total := models.XopGroupGetListForSelect(pkid)
	self.ajaxList("", MSG_OK, total, list)
}
