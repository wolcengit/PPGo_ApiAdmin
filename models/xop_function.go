/**********************************************
** @Des: xop_function
** @Author: wolcen
***********************************************/
package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type XopFunction struct {
	Id          int    `orm:"column(id);pk;auto;unique" json:"id"`
	ProdId      int    `orm:"column(prod_id);type(int)" json:"prod_id"`              //产品ID
	ModId       int    `orm:"column(mod_id);type(int)" json:"mod_id"`                //模块ID
	CatId       int    `orm:"column(cat_id);type(int)" json:"cat_id"`                //类别ID
	GrpId       int    `orm:"column(grp_id);type(int)" json:"grp_id"`                //分组ID
	Code        string `orm:"column(code);size(16)" json:"code"`                     //编码
	Name        string `orm:"column(name);size(64)" json:"name"`                     //名称
	XopName     string `orm:"column(xop_name);size(128)" json:"xop_name"`            //XOP名称
	Detail      string `orm:"column(detail);type(text)" json:"detail"`               //说明
	Detail2     string `orm:"column(detail2);type(text)" json:"detail2"`             //说明2
	Status      int    `orm:"column(status);type(int)" json:"status"`                //状态：1-正常
	CreateId    int    `orm:"column(create_id);type(int)" json:"create_id"`          //创建者
	UpdateId    int    `orm:"column(update_id);type(int)" json:"update_id"`          //修改者
	CreateTime  int64  `orm:"column(create_time);type(bigint)" json:"create_time"`   //创建时间
	UpdateTime  int64  `orm:"column(update_time);type(bigint)" json:"update_time"`   //修改时间
	UpdateTime2 int64  `orm:"column(update_time2);type(bigint)" json:"update_time2"` //修改时间2
	RefId       int    `orm:"column(ref_id);type(int)" json:"ref_id"`                //Reference ID
}

func (a *XopFunction) TableName() string {
	return TableName("xop_function")
}

func (a *XopFunction) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func XopFunctionAdd(a *XopFunction) (int64, error) {
	a.UpdateId = a.CreateId
	if a.CreateTime == 0 {
		a.CreateTime = time.Now().Unix()
		a.UpdateTime = a.CreateTime
	}
	return orm.NewOrm().Insert(a)
}

func XopFunctionDelete(id int) (int64, error) {
	o, _ := XopFunctionGetById(id)
	return orm.NewOrm().Delete(&o)
}

func XopFunctionGetById(id int) (*XopFunction, error) {
	r := new(XopFunction)
	err := orm.NewOrm().QueryTable(new(XopFunction)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func XopFunctionGetByRefId(id int) (*XopFunction, error) {
	r := new(XopFunction)
	err := orm.NewOrm().QueryTable(new(XopFunction)).Filter("ref_id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func XopFunctionGetByCode(prodId int, code string) (*XopFunction, error) {
	a := new(XopFunction)
	err := orm.NewOrm().QueryTable(new(XopFunction)).Filter("prod_id", prodId).Filter("code", code).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func XopFunctionGetByName(prodId int, name string) (*XopFunction, error) {
	a := new(XopFunction)
	err := orm.NewOrm().QueryTable(new(XopFunction)).Filter("prod_id", prodId).Filter("name", name).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func XopFunctionGetList(page, pageSize int, filters ...interface{}) ([]*XopFunction, int64) {
	offset := (page - 1) * pageSize
	list := make([]*XopFunction, 0)
	query := orm.NewOrm().QueryTable(new(XopFunction))
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

func XopFunctionGetListForBrowse(page, pageSize int, filters ...interface{}) ([]map[string]interface{}, int64) {
	result, total := XopFunctionGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["code"] = v.Code
		row["name"] = v.Name
		row["xop_name"] = v.XopName
		//row["detail"] = v.Detail
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		row["mod_id"] = v.ModId
		m, _ := XopModuleGetById(v.ModId)
		row["mod_code"] = m.Code
		row["mod_name"] = m.Name
		row["cat_id"] = v.CatId
		c, _ := XopCategoryGetById(v.CatId)
		row["cat_code"] = c.Code
		row["cat_name"] = c.Name
		row["grp_id"] = v.CatId
		g, _ := XopGroupGetById(v.GrpId)
		row["grp_code"] = g.Code
		row["grp_name"] = g.Name

		list[k] = row
	}
	return list, total
}
