/**********************************************
** @Des: xop_document
** @XopDocumentor: wolcen
***********************************************/
package models

import (
	"github.com/astaxie/beego/orm"
)

type XopDocument struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	Pid        int    `orm:"column(pid);type(int)" json:"pid"`
	ProdId     int    `orm:"column(prod_id);type(int)" json:"prod_id"` //产品ID
	Name       string `orm:"column(name);size(64)" json:"name"`
	Sort       int    `orm:"column(sort);type(int)" json:"sort"`
	Detail     string `orm:"column(detail);type(text)" json:"detail"`
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
}

func (a *XopDocument) TableName() string {
	return TableName("xop_document")
}

func (a *XopDocument) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func XopDocumentAdd(auth *XopDocument) (int64, error) {
	return orm.NewOrm().Insert(auth)
}

func XopDocumentDelete(id int) (int64, error) {
	o, _ := XopDocumentGetById(id)
	return orm.NewOrm().Delete(&o)
}

func XopDocumentGetById(id int) (*XopDocument, error) {
	a := new(XopDocument)

	err := orm.NewOrm().QueryTable(new(XopDocument)).Filter("id", id).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func XopDocumentGetList(page, pageSize int, filters ...interface{}) ([]*XopDocument, int64) {
	offset := (page - 1) * pageSize
	list := make([]*XopDocument, 0)
	query := orm.NewOrm().QueryTable(new(XopDocument))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("pid", "sort").Limit(pageSize, offset).All(&list)

	return list, total
}

func XopDocumentTopRecord(prodId int) error {
	a := new(XopDocument)

	err := orm.NewOrm().QueryTable(new(XopDocument)).Filter("pid", 0).Filter("prod_id", prodId).One(a)
	if err != nil {
		return err
	}
	XopDocumentAdd(&XopDocument{Id: 1, Pid: 0, ProdId: prodId, Name: "XOP平台说明", Detail: "", Sort: 1, CreateId: 1, Status: 1})
	return nil
}
