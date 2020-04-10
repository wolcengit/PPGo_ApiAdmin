package controllers

import (
	"encoding/json"
	"github.com/george518/PPGo_ApiAdmin/models"
	"time"
)

type ApiController struct {
	BaseController
}

type XopNode struct {
	Level   int    `json:"level"` // 1-module 2-category 3-group 4-function
	Mod     int    `json:"mod"`
	Cat     int    `json:"cat"`
	Grp     int    `json:"grp"`
	Func    int    `json:"func"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Detail  string `json:"detail"`
	Detail2 string `json:"detail2"`
}

func (self *ApiController) XopNode() {
	var req XopNode
	json.Unmarshal(self.Ctx.Input.RequestBody, &req)

	if req.Level == 1 {
		self.xopModule(req)
	}else if req.Level == 2 {
		self.xopCategory(req);
	}else if req.Level == 3 {
		self.xopGroup(req);
	}else if req.Level == 4 {
		self.xopFunction(req);
	}

	self.ajaxMsg("Error Level", MSG_ERR)
}

func (self *ApiController) xopModule(req XopNode) {
	entity, err := models.XopModuleGetByRefId(req.Mod)
	if err != nil {
		entity = new(models.XopModule)
		entity.ProdId = self.prodId
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.Status = 1
		entity.RefId = req.Mod

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		id, err := models.XopModuleAdd(entity)
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(id, MSG_OK)
	} else {
		if entity.ProdId != self.prodId {
			self.ajaxMsg("产品不匹配", MSG_ERR)
		}

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		err = entity.Update()
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(entity.Id, MSG_OK)
	}

}

func (self *ApiController) xopCategory(req XopNode) {
	entity, err := models.XopCategoryGetByRefId(req.Cat)
	if err != nil {
		module, err := models.XopModuleGetByRefId(req.Mod)
		if err != nil {
			self.ajaxMsg("模块不存在", MSG_ERR)
		}

		entity = new(models.XopCategory)
		entity.ProdId = self.prodId
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.Status = 1
		entity.RefId = req.Cat
		entity.ModId = module.Id

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		id, err := models.XopCategoryAdd(entity)
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(id, MSG_OK)
	} else {
		if entity.ProdId != self.prodId {
			self.ajaxMsg("产品不匹配", MSG_ERR)
		}

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		err = entity.Update()
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(entity.Id, MSG_OK)
	}

}

func (self *ApiController) xopGroup(req XopNode) {
	entity, err := models.XopGroupGetByRefId(req.Grp)
	if err != nil {
		module, err := models.XopModuleGetByRefId(req.Mod)
		if err != nil {
			self.ajaxMsg("模块不存在", MSG_ERR)
		}
		category, err := models.XopCategoryGetByRefId(req.Cat)
		if err != nil {
			self.ajaxMsg("类别不存在", MSG_ERR)
		}

		entity = new(models.XopGroup)
		entity.ProdId = self.prodId
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.Status = 1
		entity.RefId = req.Grp
		entity.ModId = module.Id
		entity.CatId = category.Id

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		id, err := models.XopGroupAdd(entity)
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(id, MSG_OK)
	} else {
		if entity.ProdId != self.prodId {
			self.ajaxMsg("产品不匹配", MSG_ERR)
		}

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		err = entity.Update()
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(entity.Id, MSG_OK)
	}

}

func (self *ApiController) xopFunction(req XopNode) {
	entity, err := models.XopFunctionGetByRefId(req.Func)
	if err != nil {
		module, err := models.XopModuleGetByRefId(req.Mod)
		if err != nil {
			self.ajaxMsg("模块不存在", MSG_ERR)
		}

		entity = new(models.XopFunction)
		entity.ProdId = self.prodId
		entity.CreateId = self.userId
		entity.CreateTime = time.Now().Unix()
		entity.Status = 1
		entity.RefId = req.Func
		entity.ModId = module.Id
		category, err := models.XopCategoryGetByRefId(req.Cat)
		if err == nil {
			entity.CatId = category.Id
		}
		group, err := models.XopGroupGetByRefId(req.Grp)
		if err == nil {
			entity.GrpId = group.Id
		}

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail
		entity.Detail2 = req.Detail2

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		id, err := models.XopFunctionAdd(entity)
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(id, MSG_OK)
	} else {
		if entity.ProdId != self.prodId {
			self.ajaxMsg("产品不匹配", MSG_ERR)
		}

		entity.Code = req.Code
		entity.Name = req.Name
		entity.Detail = req.Detail
		entity.Detail2 = req.Detail2

		entity.UpdateId = self.userId
		entity.UpdateTime = time.Now().Unix()

		err = entity.Update()
		if err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg(entity.Id, MSG_OK)
	}

}
