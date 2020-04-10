/**********************************************
** @Des: xop_product
** @Author: wolcen
***********************************************/
package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type XopProduct struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	Code       string `orm:"column(code);size(16)" json:"code"`                   //编码
	Name       string `orm:"column(name);size(64)" json:"name"`                   //名称
	Detail     string `orm:"column(detail);size(1024)" json:"detail"`             //说明
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
}

func (a *XopProduct) TableName() string {
	return TableName("xop_product")
}

func (a *XopProduct) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func XopProductAdd(a *XopProduct) (int64, error) {
	a.UpdateId = a.CreateId
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = a.CreateTime
	return orm.NewOrm().Insert(a)
}

func XopProductDelete(id int) (int64, error) {
	o, _ := XopProductGetById(id)
	return orm.NewOrm().Delete(&o)
}

func XopProductGetById(id int) (*XopProduct, error) {
	r := new(XopProduct)
	err := orm.NewOrm().QueryTable(new(XopProduct)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func XopProductGetByCode(code string) (*XopProduct, error) {
	a := new(XopProduct)
	err := orm.NewOrm().QueryTable(new(XopProduct)).Filter("code", code).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func XopProductGetByName(name string) (*XopProduct, error) {
	a := new(XopProduct)
	err := orm.NewOrm().QueryTable(new(XopProduct)).Filter("name", name).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func XopProductGetList(page, pageSize int, filters ...interface{}) ([]*XopProduct, int64) {
	offset := (page - 1) * pageSize
	list := make([]*XopProduct, 0)
	query := orm.NewOrm().QueryTable(new(XopProduct))
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

func XopProductGetListForBrowse(page, pageSize int, filters ...interface{}) ([]map[string]interface{}, int64) {
	result, total := XopProductGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["code"] = v.Code
		row["name"] = v.Name
		row["detail"] = v.Detail
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")

		list[k] = row
	}
	return list, total
}

func XopProductGetListForSelect() ([]map[string]interface{}, int64) {
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, total := XopProductGetList(0, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["code"] = v.Code
		row["name"] = v.Name
		list[k] = row
	}
	return list, total
}
