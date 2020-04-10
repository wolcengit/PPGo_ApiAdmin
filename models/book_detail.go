/**********************************************
** @Des: book_detail
** @Author: wolcen
***********************************************/
package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type BookDetail struct {
	Id         int   `orm:"column(id);pk;auto;unique" json:"id"`
	BookId     int   `orm:"column(book_id);type(int)" json:"book_id"`
	DetailId   int   `orm:"column(detail_id);type(int)" json:"detail_id"`
	CreateId   int   `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	CreateTime int64 `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
}

func (a *BookDetail) TableName() string {
	return TableName("book_detail")
}

func (a *BookDetail) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func BookDetailAdd(a *BookDetail) (int64, error) {
	a.CreateTime = time.Now().Unix()
	return orm.NewOrm().Insert(a)
}

func BookDetailDelete(id int) (int64, error) {
	o, _ := BookDetailGetById(id)
	return orm.NewOrm().Delete(&o)
}

func BookDetailDeleteByBook(bid int) error {
	_, err := orm.NewOrm().Raw("DELETE FROM "+new(BookDetail).TableName()+" WHERE book_id = ?", bid).Exec()
	return err
}

func BookDetailCreteOrUpdateByBook(bid int, uid int, ids string) {
	BookDetailDeleteByBook(bid)
	adid := strings.Split(ids, ",")
	for _, did := range adid {
		if len(did) == 0 {
			continue
		}
		detil := new(BookDetail)
		detil.BookId = bid
		detil.DetailId, _ = strconv.Atoi(did)
		detil.CreateId = uid
		orm.NewOrm().Insert(detil)
	}
	return
}

func BookDetailGetById(id int) (*BookDetail, error) {
	r := new(BookDetail)
	err := orm.NewOrm().QueryTable(new(BookDetail)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func BookDetailGetList(page, pageSize int, filters ...interface{}) ([]*BookDetail, int64) {
	offset := (page - 1) * pageSize
	list := make([]*BookDetail, 0)
	query := orm.NewOrm().QueryTable(new(BookDetail))
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

func BookDetailGetListForSelect(detail int) ([]map[string]interface{}, int64) {
	filters := make([]interface{}, 0)
	filters = append(filters, "detail_id", detail)
	result, total := BookDetailGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))

	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["detail_id"] = v.DetailId
		row["book_id"] = v.BookId
		v, _ := BookLibraryGetById(v.BookId)
		row["name"] = v.Name
		list[k] = row
	}
	return list, total
}
