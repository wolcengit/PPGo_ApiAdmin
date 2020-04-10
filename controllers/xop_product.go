/**********************************************
** @Des: xop product
** @Author: wolcen
***********************************************/
package controllers

import (
	"strings"
	"time"

	"github.com/george518/PPGo_ApiAdmin/models"
)

type XopProductController struct {
	BaseController
}

func (self *XopProductController) List() {
	self.Data["pageTitle"] = "XOP产品"
	self.display()
}

func (self *XopProductController) AjaxTable() {
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
	filters = append(filters, "status", 1)
	if searchName != "" {
		filters = append(filters, "name__icontains", searchName)
	}
	list, count := models.XopProductGetListForBrowse(page, self.pageSize, filters...)
	self.ajaxList("成功", MSG_OK, count, list)
}

func (self *XopProductController) Add() {
	self.Data["pageTitle"] = "新增产品"
	self.display()
}

func (self *XopProductController) Edit() {
	self.Data["pageTitle"] = "编辑产品"

	id, _ := self.GetInt("id", 0)
	v, err := models.XopProductGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = v.Id
	row["code"] = v.Code
	row["name"] = v.Name
	row["detail"] = v.Detail
	self.Data["entity"] = row

	self.display()
}

func (self *XopProductController) Auth() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "授权产品"

	id, _ := self.GetInt("id", 0)
	v, err := models.XopProductGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = v.Id
	row["code"] = v.Code
	row["name"] = v.Name
	row["detail"] = v.Detail
	self.Data["entity"] = row

	self.display()
}

func (self *XopProductController) AjaxSave() {
	pkid, _ := self.GetInt("id")
	var entity *models.XopProduct
	if pkid == 0 {
		entity = new(models.XopProduct)
	} else {
		entity, _ = models.XopProductGetById(pkid)
	}

	entity.Code = strings.TrimSpace(self.GetString("code"))
	entity.Name = strings.TrimSpace(self.GetString("name"))
	entity.Detail = strings.TrimSpace(self.GetString("detail"))
	entity.Status = 1

	// 检查是否已经存在
	entity1, err := models.XopProductGetByCode(entity.Code)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("产品编码已经存在", MSG_ERR)
	}

	entity1, err = models.XopProductGetByName(entity.Name)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("产品名称已经存在", MSG_ERR)
	}

	if pkid == 0 {
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		prodId, err := models.XopProductAdd(entity)
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		models.XopDocumentTopRecord(int(prodId))
		entity, _ = models.XopProductGetById(int(prodId))
	} else {
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if err := entity.Update(); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	}
	self.ajaxMsg(entity.Id, MSG_OK)
}

func (self *XopProductController) AjaxAuthSave() {
	pkid, _ := self.GetInt("id")

	nodes_data := strings.TrimSpace(self.GetString("nodes_data"))

	models.XopProductAuthCreteOrUpdateByProduct(pkid, self.userId, nodes_data)
	self.ajaxMsg("", MSG_OK)
}

func (self *XopProductController) AjaxDel() {

	pkid, _ := self.GetInt("id")
	_, err := models.XopProductDelete(pkid)
	if err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}

func (self *XopProductController) AjaxUserNodes() {
	prodId, _ := self.GetInt("id")
	checked, _ := self.GetInt("checked")
	list, err := models.XopProductAuthGetListForProduct(prodId)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	if checked == 1 {
		len := 0
		for _, v := range list {
			if v.Checked > 0 {
				len++
			}
		}
		list1 := make([]map[string]interface{}, len)
		len = 0
		for _, v := range list {
			if v.Checked > 0 {
				row := make(map[string]interface{})
				row["id"] = v.Id
				row["pId"] = v.Pid
				row["name"] = v.Name
				row["open"] = true
				row["checked"] = (v.Checked > 0)
				list1[len] = row
				len++
			}
		}
		self.ajaxList("成功", MSG_OK, int64(len), list1)
	} else {
		list1 := make([]map[string]interface{}, len(list))
		for k, v := range list {
			row := make(map[string]interface{})
			row["id"] = v.Id
			row["pId"] = v.Pid
			row["name"] = v.Name
			row["open"] = true
			row["checked"] = (v.Checked > 0)
			list1[k] = row
		}

		self.ajaxList("成功", MSG_OK, int64(len(list1)), list1)
	}
}
