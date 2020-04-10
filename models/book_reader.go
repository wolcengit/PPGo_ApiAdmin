/**********************************************
** @Des: book_reader
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

type BookReader struct {
	Id         int   `orm:"column(id);pk;auto;unique" json:"id"`
	BookId     int   `orm:"column(book_id);type(int)" json:"book_id"`
	ReaderId   int   `orm:"column(reader_id);type(int)" json:"reader_id"`
	CreateId   int   `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	CreateTime int64 `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
}

func (a *BookReader) TableName() string {
	return TableName("book_reader")
}

func (a *BookReader) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func BookReaderAdd(a *BookReader) (int64, error) {
	a.CreateTime = time.Now().Unix()
	return orm.NewOrm().Insert(a)
}

func BookReaderDeleteByBook(bid int) error {
	_, err := orm.NewOrm().Raw("DELETE FROM "+new(BookReader).TableName()+" WHERE book_id = ?", bid).Exec()
	return err
}

func BookReaderCreteOrUpdateByBook(bid int, uid int, ids string) {
	BookReaderDeleteByBook(bid)
	adid := strings.Split(ids, ",")
	for _, did := range adid {
		if len(did) == 0 {
			continue
		}
		detil := new(BookReader)
		detil.BookId = bid
		detil.ReaderId, _ = strconv.Atoi(did)
		detil.CreateId = uid
		orm.NewOrm().Insert(detil)
	}
	return
}

func BookReaderDelete(id int) (int64, error) {
	o, _ := BookReaderGetById(id)
	return orm.NewOrm().Delete(&o)
}

func BookReaderGetById(id int) (*BookReader, error) {
	r := new(BookReader)
	err := orm.NewOrm().QueryTable(new(BookReader)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func BookReaderGetList(page, pageSize int, filters ...interface{}) ([]*BookReader, int64) {
	offset := (page - 1) * pageSize
	list := make([]*BookReader, 0)
	query := orm.NewOrm().QueryTable(new(BookReader))
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

func BookReaderGetListForSelect(reader int) ([]map[string]interface{}, int64) {
	filters := make([]interface{}, 0)
	filters = append(filters, "reader_id", reader)
	result, total := BookReaderGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))

	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["reader_id"] = v.ReaderId
		row["book_id"] = v.BookId
		v, _ := BookLibraryGetById(v.BookId)
		row["name"] = v.Name
		list[k] = row
	}
	return list, total
}

func BookReaderGetListForBook(bookId int) ([]BookNode, error) {
	var buf bytes.Buffer
	buf.WriteString("SELECT 1 as id,0 as pid,'所有人员' as name,(SELECT COUNT(1) FROM " + new(BookReader).TableName() + " WHERE book_id = ? AND reader_id = 1 ) as checked ")
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT a.id,1,a.real_name,(SELECT COUNT(1) FROM " + new(BookReader).TableName() + " WHERE book_id = ? AND reader_id = a.id ) as checked FROM ")
	buf.WriteString(new(Admin).TableName())
	buf.WriteString(" as a WHERE a.id > 1  ")
	sql := buf.String()
	var nodes []BookNode
	_, err := orm.NewOrm().Raw(sql, bookId, bookId).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
