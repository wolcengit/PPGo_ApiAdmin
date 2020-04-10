/**********************************************
** @Des: book category
** @Author: wolcen
***********************************************/
package controllers

import (
	"strings"
	"time"

	"github.com/george518/PPGo_ApiAdmin/models"
)

type BookCategoryController struct {
	BaseController
}

func (self *BookCategoryController) List() {
	self.Data["pageTitle"] = "书籍类别"
	self.display()
}

func (self *BookCategoryController) Add() {
	self.Data["pageTitle"] = "新增类别"

	self.display()
}

func (self *BookCategoryController) Edit() {
	self.Data["pageTitle"] = "编辑类别"

	id, _ := self.GetInt("id", 0)
	product, _ := models.BookCategoryGetById(id)
	row := make(map[string]interface{})
	row["id"] = product.Id
	row["name"] = product.Name
	row["detail"] = product.Detail
	self.Data["entity"] = row

	self.display()
}

func (self *BookCategoryController) AjaxTable() {
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
	list, count := models.BookCategoryGetListForBrowse(page, self.pageSize, filters...)
	self.ajaxList("成功", MSG_OK, count, list)
}
func (self *BookCategoryController) AjaxSave() {
	pkid, _ := self.GetInt("id")
	var entity *models.BookCategory
	if pkid == 0 {
		entity = new(models.BookCategory)
	} else {
		entity, _ = models.BookCategoryGetById(pkid)
	}

	entity.Name = strings.TrimSpace(self.GetString("name"))
	entity.Detail = strings.TrimSpace(self.GetString("detail"))
	entity.Status = 1

	// 检查是否已经存在
	entity1, err := models.BookCategoryGetByName(self.prodId, entity.Name)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("类别名称已经存在", MSG_ERR)
	}

	if pkid == 0 {
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if _, err := models.BookCategoryAdd(entity); err != nil {
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

func (self *BookCategoryController) AjaxDel() {

	pkid, _ := self.GetInt("id")
	_, err := models.BookCategoryDelete(pkid)
	if err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}
