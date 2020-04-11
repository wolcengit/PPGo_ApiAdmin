/**********************************************
** @Des: xop function
** @Author: wolcen
***********************************************/
package controllers

import (
	"strings"
	"time"

	"github.com/george518/PPGo_ApiAdmin/models"
)

type XopFunctionController struct {
	BaseController
}

func (self *XopFunctionController) List() {
	self.Data["pageTitle"] = "XOP函数"
	list, _ := models.XopModuleGetListForSelect(self.prodId)
	self.Data["Modules"] = list
	self.display()
}

func (self *XopFunctionController) AjaxTable() {
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
	grpId, _ := self.GetInt("grpId")

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
		filters = append(filters, "cat_id", catId)
	}
	if grpId > 0 {
		filters = append(filters, "grp_id", grpId)
	}
	list, count := models.XopFunctionGetListForBrowse(page, self.pageSize, filters...)
	self.ajaxList("成功", MSG_OK, count, list)
}

func (self *XopFunctionController) Add() {
	self.Data["pageTitle"] = "新增函数"
	list, _ := models.XopModuleGetListForSelect(self.prodId)
	self.Data["Modules"] = list
	self.display()
}

func (self *XopFunctionController) Edit() {
	self.Data["pageTitle"] = "编辑函数"

	id, _ := self.GetInt("id", 0)
	v, err := models.XopFunctionGetById(id)
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
	row["mod_id"] = int(v.ModId)
	row["cat_id"] = int(v.CatId)
	row["grp_id"] = int(v.GrpId)
	self.Data["entity"] = row

	list, _ := models.XopModuleGetListForSelect(self.prodId)
	self.Data["Modules"] = list
	list1, _ := models.XopCategoryGetListForSelect(v.ModId)
	self.Data["Categorys"] = list1
	list2, _ := models.XopGroupGetListForSelect(v.CatId)
	self.Data["Groups"] = list2
	self.display()
}

func (self *XopFunctionController) AjaxSave() {
	pkid, _ := self.GetInt("id")
	var entity *models.XopFunction
	if pkid == 0 {
		entity = new(models.XopFunction)
	} else {
		entity, _ = models.XopFunctionGetById(pkid)
	}

	entity.Code = strings.TrimSpace(self.GetString("code"))
	entity.Name = strings.TrimSpace(self.GetString("name"))
	entity.XopName = strings.TrimSpace(self.GetString("xop_name"))
	entity.Detail = strings.TrimSpace(self.GetString("detail"))
	detail2 := strings.TrimSpace(self.GetString("detail2"))
	entity.Status = 1
	entity.ModId, _ = self.GetInt("mod_id")
	entity.CatId, _ = self.GetInt("cat_id")
	entity.GrpId, _ = self.GetInt("grp_id")
	refid, _ := self.GetInt("ref_id")
	if refid > 0 {
		entity.RefId = refid
	}

	// 检查是否已经存在
	entity1, err := models.XopFunctionGetByCode(self.prodId, entity.Code)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("函数编码已经存在", MSG_ERR)
	}

	entity1, err = models.XopFunctionGetByName(self.prodId, entity.Name)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("函数名称已经存在", MSG_ERR)
	}

	if pkid == 0 {
		entity.ProdId = self.prodId
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if detail2 != "" {
			entity.Detail2 = detail2
			entity.UpdateTime2 = time.Now().Unix()
		}
		if _, err := models.XopFunctionAdd(entity); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	} else {
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if detail2 != "" {
			entity.Detail2 = detail2
			entity.UpdateTime2 = time.Now().Unix()
		}
		if err := entity.Update(); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	}
	self.ajaxMsg("", MSG_OK)
}

func (self *XopFunctionController) AjaxDel() {

	pkid, _ := self.GetInt("id")
	_, err := models.XopFunctionDelete(pkid)
	if err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}
