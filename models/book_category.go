/**********************************************
** @Des: book_category
** @Author: wolcen
***********************************************/
package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type BookCategory struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	ProdId     int    `orm:"column(prod_id);type(int)" json:"prod_id"`            //产品ID
	Name       string `orm:"column(name);size(64)" json:"name"`                   //名称
	Detail     string `orm:"column(detail);size(1024)" json:"detail"`             //说明
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
}

func (a *BookCategory) TableName() string {
	return TableName("book_category")
}

func (a *BookCategory) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func BookCategoryAdd(a *BookCategory) (int64, error) {
	a.UpdateId = a.CreateId
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = a.CreateTime
	return orm.NewOrm().Insert(a)
}

func BookCategoryDelete(id int) (int64, error) {
	o, _ := BookCategoryGetById(id)
	return orm.NewOrm().Delete(&o)
}

func BookCategoryGetById(id int) (*BookCategory, error) {
	r := new(BookCategory)
	err := orm.NewOrm().QueryTable(new(BookCategory)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func BookCategoryGetByName(prodId int, name string) (*BookCategory, error) {
	a := new(BookCategory)
	err := orm.NewOrm().QueryTable(new(BookCategory)).Filter("prod_id", prodId).Filter("name", name).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func BookCategoryGetList(page, pageSize int, filters ...interface{}) ([]*BookCategory, int64) {
	offset := (page - 1) * pageSize
	list := make([]*BookCategory, 0)
	query := orm.NewOrm().QueryTable(new(BookCategory))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

func BookCategoryGetListForBrowse(page, pageSize int, filters ...interface{}) ([]map[string]interface{}, int64) {
	result, total := BookCategoryGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		row["detail"] = v.Detail
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		list[k] = row
	}
	return list, total
}

func BookCategoryGetListForSelect(prodId int) ([]map[string]interface{}, int64) {
	filters := make([]interface{}, 0)
	filters = append(filters, "prod_id", prodId)
	filters = append(filters, "status", 1)
	result, total := BookCategoryGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))

	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		list[k] = row
	}
	return list, total
}
