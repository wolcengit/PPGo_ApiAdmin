/**********************************************
** @Des: book_library
** @Author: wolcen
***********************************************/
package models

import (
	"bytes"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type BookLibrary struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	ProdId     int    `orm:"column(prod_id);type(int)" json:"prod_id"`            //产品ID
	CatId      int    `orm:"column(cat_id);type(int)" json:"cat_id"`              //类别ID
	Name       string `orm:"column(name);size(64)" json:"name"`                   //名称
	Detail     string `orm:"column(detail);size(1024)" json:"detail"`             //说明
	Opened     int    `orm:"column(opened);type(int)" json:"opened"`              //公开：0-公开 1-私有
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
}

type BookCatNode struct {
	Id       int            `json:"id"`
	Title    string         `json:"title"`
	Children []*BookCatNode `json:"children"`
}

func (a *BookLibrary) TableName() string {
	return TableName("book_library")
}

func (a *BookLibrary) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func BookLibraryAdd(a *BookLibrary) (int64, error) {
	a.UpdateId = a.CreateId
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = a.CreateTime
	return orm.NewOrm().Insert(a)
}

func BookLibraryDelete(id int) (int64, error) {
	o, _ := BookLibraryGetById(id)
	return orm.NewOrm().Delete(&o)
}

func BookLibraryGetById(id int) (*BookLibrary, error) {
	r := new(BookLibrary)
	err := orm.NewOrm().QueryTable(new(BookLibrary)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func BookLibraryGetByName(prodId int, name string) (*BookLibrary, error) {
	a := new(BookLibrary)
	err := orm.NewOrm().QueryTable(new(BookLibrary)).Filter("prod_id", prodId).Filter("name", name).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func BookLibraryGetList(page, pageSize int, filters ...interface{}) ([]*BookLibrary, int64) {
	offset := (page - 1) * pageSize
	list := make([]*BookLibrary, 0)
	query := orm.NewOrm().QueryTable(new(BookLibrary))
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

func BookLibraryGetListForBrowse(page, pageSize int, filters ...interface{}) ([]map[string]interface{}, int64) {
	result, total := BookLibraryGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		row["detail"] = v.Detail
		row["opened"] = v.Opened
		if v.Opened == 0 {
			row["opened_name"] = "公开"
		} else {
			row["opened_name"] = "私有"
		}
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		row["cat_id"] = v.CatId
		c, _ := BookCategoryGetById(v.CatId)
		row["cat_name"] = c.Name

		list[k] = row
	}
	return list, total
}

func BookLibraryGetListForUser(page, pageSize int, prodId int, searchName string, usrid int) ([]map[string]interface{}, int64) {
	var buf bytes.Buffer
	buf.WriteString("SELECT * FROM ")
	buf.WriteString(new(BookLibrary).TableName())
	if searchName != "" {
		buf.WriteString(" WHERE prod_id = ?  ")
		buf.WriteString(" AND name like '%" + searchName + "%'  ")
	} else {
		buf.WriteString(" WHERE prod_id = ?  ")
	}
	if usrid > 1 {
		buf.WriteString(" AND id IN(SELECT book_id FROM " + new(BookReader).TableName() + " WHERE reader_id = ?)")
	}
	sql := buf.String()
	var result []BookLibrary
	if usrid > 1 {
		_, err := orm.NewOrm().Raw(sql, prodId, usrid).QueryRows(&result)
		if err != nil {
			return nil, 0
		}
	} else {
		_, err := orm.NewOrm().Raw(sql, prodId).QueryRows(&result)
		if err != nil {
			return nil, 0
		}

	}
	total := len(result)
	list := make([]map[string]interface{}, total)
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		row["detail"] = v.Detail
		row["opened"] = v.Opened
		if v.Opened == 0 {
			row["opened_name"] = "公开"
		} else {
			row["opened_name"] = "私有"
		}
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		row["cat_id"] = v.CatId

		list[k] = row
	}
	return list, int64(total)
}

type BookNode struct {
	Id      int    `json:"id"`
	Pid     int    `json:"pId"`
	Name    string `json:"name"`
	Checked int    `json:"checked"`
}

func BookLibraryGetNodes() ([]BookNode, error) {
	var buf bytes.Buffer
	buf.WriteString("SELECT id,pid,name FROM ")
	buf.WriteString(new(XopDocument).TableName())
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT 10000,10000,'XOP函数清单' ")
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT id+10000,10000,name FROM ")
	buf.WriteString(new(XopModule).TableName())
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT id+20000,mod_id+10000,name FROM ")
	buf.WriteString(new(XopCategory).TableName())
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT id+30000,cat_id+20000,name FROM ")
	buf.WriteString(new(XopGroup).TableName())
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT id+40000,grp_id+30000,name FROM ")
	buf.WriteString(new(XopFunction).TableName())
	sql := buf.String()
	var nodes []BookNode
	_, err := orm.NewOrm().Raw(sql).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func BookLibraryGetBookNodes(bookid int) ([]BookNode, error) {
	var buf bytes.Buffer
	buf.WriteString("SELECT a.id,a.pid,a.name,(SELECT COUNT(1) FROM " + new(BookDetail).TableName() + " WHERE book_id = ? AND detail_id = a.id ) as checked FROM ")
	buf.WriteString(new(XopDocument).TableName())
	buf.WriteString(" as a UNION ")
	buf.WriteString("SELECT 10000,10000,'XOP函数清单', (SELECT COUNT(1) FROM " + new(BookDetail).TableName() + " WHERE book_id = ? AND detail_id = 10000 ) as checked")
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT a.id+10000,10000,a.name,(SELECT COUNT(1) FROM " + new(BookDetail).TableName() + " WHERE book_id = ? AND detail_id = a.id+10000 ) as checked FROM ")
	buf.WriteString(new(XopModule).TableName())
	buf.WriteString(" as a UNION ")
	buf.WriteString("SELECT a.id+20000,a.mod_id+10000,a.name,(SELECT COUNT(1) FROM " + new(BookDetail).TableName() + " WHERE book_id = ? AND detail_id = a.id+20000 ) as checked FROM ")
	buf.WriteString(new(XopCategory).TableName())
	buf.WriteString(" as a UNION ")
	buf.WriteString("SELECT a.id+30000,a.cat_id+20000,a.name,(SELECT COUNT(1) FROM " + new(BookDetail).TableName() + " WHERE book_id = ? AND detail_id = a.id+30000 ) as checked FROM ")
	buf.WriteString(new(XopGroup).TableName())
	buf.WriteString(" as a UNION ")
	buf.WriteString("SELECT a.id+40000,a.grp_id+30000,a.name,(SELECT COUNT(1) FROM " + new(BookDetail).TableName() + " WHERE book_id = ? AND detail_id = a.id+40000 ) as checked FROM ")
	buf.WriteString(new(XopFunction).TableName())
	buf.WriteString(" as a ")
	sql := buf.String()
	var nodes []BookNode
	_, err := orm.NewOrm().Raw(sql, bookid, bookid, bookid, bookid, bookid, bookid).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func BookLibraryGetNodeLists() ([]BookNode, error) {
	var buf bytes.Buffer
	buf.WriteString("SELECT id,id as pid ,name FROM ")
	buf.WriteString(new(BookCategory).TableName())
	buf.WriteString(" UNION ")
	buf.WriteString("SELECT id+10000,cat_id,name FROM ")
	buf.WriteString(new(BookLibrary).TableName())
	sql := buf.String()
	var nodes []BookNode
	_, err := orm.NewOrm().Raw(sql).QueryRows(&nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
