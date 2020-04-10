/**********************************************
** @Des: book document
** @XopDocumentor: wolcen
***********************************************/

package controllers

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/george518/PPGo_ApiAdmin/utils"

	"strconv"

	"github.com/george518/PPGo_ApiAdmin/models"
	"github.com/patrickmn/go-cache"
)

type XopDocumentController struct {
	BaseController
}

func (self *XopDocumentController) Index() {

	self.Data["pageTitle"] = "XOP文档"
	self.display()
}

func (self *XopDocumentController) List() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "XOP文档"
	self.display()
}
func (self *XopDocumentController) Edit() {
	id, _ := self.GetInt("id")
	document, err := models.XopDocumentGetById(id)
	if err != nil {
		self.ajaxMsg("查找文档失败！", MSG_ERR)
	}
	if document.Detail == ""{
		document.Detail = "[TOC]"
	}
	self.Data["entity"] = document
	self.display()
}

//获取全部节点
func (self *XopDocumentController) AjaxNodes() {
	filters := make([]interface{}, 0)
	filters = append(filters, "prod_id", self.prodId)
	filters = append(filters, "status", 1)
	result, count := models.XopDocumentGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["pId"] = v.Pid
		row["name"] = v.Name
		row["open"] = true
		list[k] = row
	}

	self.ajaxList("成功", MSG_OK, count, list)
}

//获取一个节点
func (self *XopDocumentController) AjaxNode() {
	id, _ := self.GetInt("id")
	result, _ := models.XopDocumentGetById(id)
	// if err == nil {
	// 	self.ajaxMsg(err.Error(), MSG_ERR)
	// }
	row := make(map[string]interface{})
	row["id"] = result.Id
	row["pid"] = result.Pid
	row["name"] = result.Name
	row["detail"] = result.Detail
	row["sort"] = result.Sort

	fmt.Println(row)

	self.ajaxList("成功", MSG_OK, 0, row)
}

//新增或修改
func (self *XopDocumentController) AjaxSave() {
	id, _ := self.GetInt("id")
	if id == 1 {
		self.ajaxMsg("系统根节点，无法修改", MSG_ERR)
	}
	var doc *models.XopDocument
	if id == 0 {
		doc = new(models.XopDocument)
	} else {
		doc, _ = models.XopDocumentGetById(id)
	}

	doc.Pid, _ = self.GetInt("pid")
	doc.Name = strings.TrimSpace(self.GetString("name"))
	doc.Sort, _ = self.GetInt("sort")
	doc.UpdateTime = time.Now().Unix()

	if doc.Pid == 0 {
		doc.Pid = 1
	}
	if id == 0 {
		//新增
		doc.Status = 1
		doc.ProdId = self.prodId
		doc.CreateTime = time.Now().Unix()
		doc.CreateId = self.userId
		doc.UpdateId = self.userId
		if _, err := models.XopDocumentAdd(doc); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	} else {
		doc.Id = id
		doc.UpdateId = self.userId
		if err := doc.Update(); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	}
	utils.Che.Set("menu"+strconv.Itoa(self.user.Id), nil, cache.DefaultExpiration)
	self.ajaxMsg("", MSG_OK)
}

func (self *XopDocumentController) AjaxSaveDetail() {
	id, _ := self.GetInt("id")
	if id == 1 {
		self.ajaxMsg("系统根节点，无法修改", MSG_ERR)
	}
	doc, err := models.XopDocumentGetById(id)
	if err != nil {
		self.ajaxMsg("查找文档失败！", MSG_ERR)
	}
	doc.Detail = strings.TrimSpace(self.GetString("detail"))
	if err := doc.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}

//删除
func (self *XopDocumentController) AjaxDel() {
	id, _ := self.GetInt("id")
	document, _ := models.XopDocumentGetById(id)
	document.Id = id
	document.Status = 0
	if err := document.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	utils.Che.Set("menu"+strconv.Itoa(self.user.Id), nil, cache.DefaultExpiration)
	self.ajaxMsg("", MSG_OK)
}


func (self *XopDocumentController) Upload() {
	if !self.isPost() {
		self.ajaxMsg("请求方式有误！", MSG_ERR)
	}
	id, _ := self.GetInt("id")
	if id == 0 {
		self.ajaxMsg("请求参数错误！", MSG_ERR)
	}
	document, err := models.XopDocumentGetById(id)
	if err != nil {
		self.ajaxMsg("查找文档失败！", MSG_ERR)
	}

	// handle upload
	f, h, err := self.GetFile("editormd-image-file")
	if err != nil {
		self.ajaxMsg("上传图片数据错误！", MSG_ERR)
	}
	if h == nil || f == nil {
		self.ajaxMsg("上传图片错误！", MSG_ERR)
	}
	_ = f.Close()

	// file save dir
	saveDir := fmt.Sprintf("/uploads/%d", document.Id)
	ok, _ := utils.File.PathIsExists(saveDir)
	if !ok {
		err := os.MkdirAll(saveDir, 0777)
		if err != nil {
			self.ajaxMsg("上传图片失败！", MSG_ERR)
		}
	}
	filename := strings.ReplaceAll(h.Filename," ","_")
	// check file is exists
	imageFile := path.Join(saveDir, filename)
	ok, _ = utils.File.PathIsExists(imageFile)
	if ok {
		self.ajaxMsg("该图片已经上传过！", MSG_ERR)
	}
	// save file
	err = self.SaveToFile("editormd-image-file", imageFile)
	if err != nil {
		self.ajaxMsg("图片保存失败！", MSG_ERR)
	}
	saveUrl := fmt.Sprintf("/uploads/%d/&s", document.Id,filename)


	self.ajaxMsg(saveUrl, MSG_OK)
}

