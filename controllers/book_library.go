/**********************************************
** @Des: book list
** @Author: wolcen
***********************************************/
package controllers

import (
	"strings"
	"time"

	"github.com/george518/PPGo_ApiAdmin/models"
)

type BookLibraryController struct {
	BaseController
}

func (self *BookLibraryController) List() {
	self.Data["pageTitle"] = "XOP书籍"
	list, _ := models.BookCategoryGetListForSelect(self.prodId)
	self.Data["Categorys"] = list
	self.display()
}

func (self *BookLibraryController) Add() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "新增书籍"
	list, _ := models.BookCategoryGetListForSelect(self.prodId)
	self.Data["Categorys"] = list
	self.display()
}

func (self *BookLibraryController) Edit() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "编辑书籍"

	id, _ := self.GetInt("id", 0)
	v, err := models.BookLibraryGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = v.Id
	row["name"] = v.Name
	row["detail"] = v.Detail
	row["opened"] = v.Opened
	row["cat_id"] = int(v.CatId)
	self.Data["entity"] = row

	list, _ := models.BookCategoryGetListForSelect(self.prodId)
	self.Data["Categorys"] = list
	self.display()
}

func (self *BookLibraryController) Auth() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "授权书籍"

	id, _ := self.GetInt("id", 0)
	v, err := models.BookLibraryGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = v.Id
	row["name"] = v.Name
	row["detail"] = v.Detail
	row["opened"] = v.Opened
	row["cat_id"] = int(v.CatId)
	c, _ := models.BookCategoryGetById(v.CatId)
	row["cat_name"] = c.Name
	self.Data["entity"] = row

	list, _ := models.BookCategoryGetListForSelect(self.prodId)
	self.Data["Categorys"] = list
	self.display()
}

func (self *BookLibraryController) Detail() {
	self.Data["zTree"] = true //引入ztreecss

	id, _ := self.GetInt("id", 0)
	book, err := models.BookLibraryGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	self.Data["pageTitle"] = book.Name
	self.Data["readBook"] = book.Id
	self.Data["readContent"] = book.Detail
	self.display()
}

func (self *BookLibraryController) AjaxBookNodeDetail() {
	id, _ := self.GetInt("id", 0)
	nid, _ := self.GetInt("nid", 0)
	book, err := models.BookLibraryGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	var md string
	md = book.Detail
	if nid < 10000 {
		v, err := models.XopDocumentGetById(nid)
		if err != nil {
			self.Ctx.WriteString("数据不存在")
			return
		}
		md = v.Detail
	} else if nid > 40000 {
		v, err := models.XopFunctionGetById(nid - 40000)
		if err != nil {
			self.Ctx.WriteString("数据不存在")
			return
		}
		if v.Detail2 != "" {
			md = v.Detail2
		} else {
			md = v.Detail
		}
	} else {
		md = ""
	}
	self.ajaxList("成功", MSG_OK, int64(0), md)
}

func (self *BookLibraryController) AjaxDownload() {
	id, _ := self.GetInt("id", 0)
	book, err := models.BookLibraryGetById(id)
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	self.Data["pageTitle"] = book.Name
	self.Data["readBook"] = book.Id
	self.Data["readContent"] = book.Detail
	self.display()
}

func (self *BookLibraryController) AjaxTable() {
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
	catId, _ := self.GetInt("catId")

	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "prod_id", self.prodId)
	filters = append(filters, "status", 1)
	if searchName != "" {
		filters = append(filters, "name__icontains", searchName)
	}
	if catId > 0 {
		filters = append(filters, "cat_id", catId)
	}
	list, count := models.BookLibraryGetListForBrowse(page, self.pageSize, filters...)
	self.ajaxList("成功", MSG_OK, count, list)
}

func (self *BookLibraryController) AjaxNodes() {
	list, err := models.BookLibraryGetNodes()
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	list1 := make([]map[string]interface{}, len(list))
	for k, v := range list {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["pId"] = v.Pid
		row["name"] = v.Name
		row["open"] = true
		list1[k] = row
	}

	self.ajaxList("成功", MSG_OK, int64(len(list1)), list1)
}

func (self *BookLibraryController) AjaxBookNodes() {
	bookid, _ := self.GetInt("id")
	checked, _ := self.GetInt("checked")
	list, err := models.BookLibraryGetBookNodes(bookid)
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

func (self *BookLibraryController) AjaxUserNodes() {
	bookid, _ := self.GetInt("id")
	checked, _ := self.GetInt("checked")
	list, err := models.BookReaderGetListForBook(bookid)
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

func (self *BookLibraryController) AjaxNodeLists() {
	list, err := models.BookLibraryGetNodeLists()
	if err != nil {
		self.Ctx.WriteString("数据不存在")
		return
	}
	list1 := make([]map[string]interface{}, len(list))
	for k, v := range list {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["pId"] = v.Pid
		row["name"] = v.Name
		row["open"] = true
		row["checked"] = false
		list1[k] = row
	}

	self.ajaxList("成功", MSG_OK, int64(len(list1)), list1)
}

func (self *BookLibraryController) AjaxSave() {
	pkid, _ := self.GetInt("id")
	var entity *models.BookLibrary
	if pkid == 0 {
		entity = new(models.BookLibrary)
	} else {
		entity, _ = models.BookLibraryGetById(pkid)
	}

	entity.Name = strings.TrimSpace(self.GetString("name"))
	entity.Detail = strings.TrimSpace(self.GetString("detail"))
	entity.Status = 1
	entity.CatId, _ = self.GetInt("cat_id")
	entity.Opened, _ = self.GetInt("opened")
	nodes_data := strings.TrimSpace(self.GetString("nodes_data"))

	// 检查是否已经存在
	entity1, err := models.BookLibraryGetByName(self.prodId, entity.Name)
	if err == nil && entity1.Id != pkid {
		self.ajaxMsg("书籍名称已经存在", MSG_ERR)
	}

	if pkid == 0 {
		entity.ProdId = self.prodId
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if _, err := models.BookLibraryAdd(entity); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		entity, _ = models.BookLibraryGetByName(self.prodId, entity.Name)
	} else {
		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()
		if err := entity.Update(); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	}
	models.BookDetailCreteOrUpdateByBook(entity.Id, self.userId, nodes_data)
	self.ajaxMsg("", MSG_OK)
}

func (self *BookLibraryController) AjaxAuthSave() {
	pkid, _ := self.GetInt("id")

	nodes_data := strings.TrimSpace(self.GetString("nodes_data"))

	models.BookReaderCreteOrUpdateByBook(pkid, self.userId, nodes_data)
	self.ajaxMsg("", MSG_OK)
}

func (self *BookLibraryController) AjaxDel() {

	pkid, _ := self.GetInt("id")
	_, err := models.BookLibraryDelete(pkid)
	if err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}
