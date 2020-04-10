/**********************************************
** @Des: xop_category
** @Author: wolcen
***********************************************/
package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type XopCategory struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	ProdId     int    `orm:"column(prod_id);type(int)" json:"prod_id"`            //产品ID
	ModId      int    `orm:"column(mod_id);type(int)" json:"mod_id"`              //模块ID
	Code       string `orm:"column(code);size(16)" json:"code"`                   //编码
	Name       string `orm:"column(name);size(64)" json:"name"`                   //名称
	Detail     string `orm:"column(detail);size(1024)" json:"detail"`             //说明
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
	RefId      int    `orm:"column(ref_id);type(int)" json:"ref_id"`              //Reference ID
}

func (a *XopCategory) TableName() string {
	return TableName("xop_category")
}

func (a *XopCategory) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func XopCategoryAdd(a *XopCategory) (int64, error) {
	a.UpdateId = a.CreateId
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = a.CreateTime
	return orm.NewOrm().Insert(a)
}

func XopCategoryDelete(id int) (int64, error) {
	o, _ := XopCategoryGetById(id)
	return orm.NewOrm().Delete(&o)
}

func XopCategoryGetById(id int) (*XopCategory, error) {
	r := new(XopCategory)
	err := orm.NewOrm().QueryTable(new(XopCategory)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func XopCategoryGetByRefId(id int) (*XopCategory, error) {
	r := new(XopCategory)
	err := orm.NewOrm().QueryTable(new(XopCategory)).Filter("ref_id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func XopCategoryGetByCode(prodId int, code string) (*XopCategory, error) {
	a := new(XopCategory)
	err := orm.NewOrm().QueryTable(new(XopCategory)).Filter("prod_id", prodId).Filter("code", code).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func XopCategoryGetByName(prodId int, name string) (*XopCategory, error) {
	a := new(XopCategory)
	err := orm.NewOrm().QueryTable(new(XopCategory)).Filter("prod_id", prodId).Filter("name", name).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func XopCategoryGetList(page, pageSize int, filters ...interface{}) ([]*XopCategory, int64) {
	offset := (page - 1) * pageSize
	list := make([]*XopCategory, 0)
	query := orm.NewOrm().QueryTable(new(XopCategory))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	// total := int64(12)
	return list, total
}

func XopCategoryGetListForBrowse(page, pageSize int, filters ...interface{}) ([]map[string]interface{}, int64) {
	result, total := XopCategoryGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["code"] = v.Code
		row["name"] = v.Name
		row["detail"] = v.Detail
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		row["mod_id"] = v.ModId
		m, _ := XopModuleGetById(v.ModId)
		row["mod_code"] = m.Code
		row["mod_name"] = m.Name

		list[k] = row
	}
	return list, total
}

func XopCategoryGetListForSelect(modId int) ([]map[string]interface{}, int64) {
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	filters = append(filters, "mod_id", modId)
	result, total := XopCategoryGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))

	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["code"] = v.Code
		row["name"] = v.Name
		row["mod_id"] = v.ModId
		list[k] = row
	}
	return list, total
}
