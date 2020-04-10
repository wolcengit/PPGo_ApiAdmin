/**********************************************
** @Des: xop_product_auth
** @Author: wolcen
***********************************************/
package models

import (
	"bytes"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type XopProductAuth struct {
	Id         int   `orm:"column(id);pk;auto;unique" json:"id"`
	ProdId     int   `orm:"column(prod_id);type(int)" json:"prod_id"`
	UserId     int   `orm:"column(user_id);type(int)" json:"user_id"`
	CreateId   int   `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	CreateTime int64 `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
}

func (a *XopProductAuth) TableName() string {
	return TableName("xop_product_auth")
}

func (a *XopProductAuth) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func XopProductAuthAdd(a *XopProductAuth) (int64, error) {
	a.CreateTime = time.Now().Unix()
	return orm.NewOrm().Insert(a)
}

func XopProductAuthDeleteByBook(bid int) error {
	_, err := orm.NewOrm().Raw("DELETE FROM "+new(XopProductAuth).TableName()+" WHERE prod_id = ?", bid).Exec()
	return err
}

func XopProductAuthCreteOrUpdateByProduct(bid int, uid int, ids string) {
	XopProductAuthDeleteByBook(bid)
	adid := strings.Split(ids, ",")
	for _, did := range adid {
		if len(did) == 0 {
			continue
		}
		detil := new(XopProductAuth)
		detil.ProdId = bid
		detil.UserId, _ = strconv.Atoi(did)
		detil.CreateId = uid
		orm.NewOrm().Insert(detil)
	}
	return
}

func XopProductAuthDelete(id int) (int64, error) {
	o, _ := XopProductAuthGetById(id)
	return orm.NewOrm().Delete(&o)
}

func XopProductAuthGetById(id int) (*XopProductAuth, error) {
	r := new(XopProductAuth)
	err := orm.NewOrm().QueryTable(new(XopProductAuth)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func XopProductAuthGetList(page, pageSize int, filters ...interface{}) ([]*XopProductAuth, int64) {
	offset := (page - 1) * pageSize
	list := make([]*XopProductAuth, 0)
	query := orm.NewOrm().QueryTable(new(XopProductAuth))
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

func XopProductAuthCheckForLogin(prodId int, userId int) bool {
	filters := make([]interface{}, 0)
	filters = append(filters, "prod_id", prodId)
	filters = append(filters, "user_id", userId)
	result, _ := XopProductAuthGetList(1, 1000, filters...)
	return len(result) > 0
}

func XopProductAuthGetListForProduct(prodId int) ([]BookNode, error) {
	var buf bytes.Buffer
	buf.WriteString("SELECT 1 as id,0 as pid,'所有人员' as name,(SELECT COUNT(1) FROM " + new(XopProductAuth).TableName() + " WHERE prod_id = ? AND user_id = 1 ) as checked ")
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT a.id,1,a.real_name,(SELECT COUNT(1) FROM " + new(XopProductAuth).TableName() + " WHERE prod_id = ? AND user_id = a.id ) as checked FROM ")
	buf.WriteString(new(Admin).TableName())
	buf.WriteString(" as a WHERE a.id > 1  ")
	sql := buf.String()
	var nodes []BookNode
	_, err := orm.NewOrm().Raw(sql, prodId, prodId).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
