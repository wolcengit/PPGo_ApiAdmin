/**********************************************
** @Des: 权限因子
** @Author: haodaquan
** @Date:   2017-09-09 20:50:36
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 21:42:08
***********************************************/
package models

import (
	"bytes"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type Auth struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	Pid        int    `orm:"column(pid);type(int)" json:"pid"`                    //上级ID，0为顶级
	UserId     int    `orm:"column(user_id);type(int)" json:"user_id"`            //操作者
	AuthName   string `orm:"column(auth_name);size(64)" json:"auth_name"`         //权限名称
	AuthUrl    string `orm:"column(auth_url);size(254)" json:"auth_url"`          //URL地址
	Sort       int    `orm:"column(sort);type(int)" json:"sort"`                  //排序，越小越前
	Icon       string `orm:"column(icon);size(254)" json:"icon"`                  //icon
	IsShow     int    `orm:"column(is_show);type(int)" json:"is_show"`            //是否显示，0-隐藏，1-显示
	Opened     int    `orm:"column(opened);type(int)" json:"opened"`              //是否管控：0-公开 1-管控
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常 0禁用
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
}

func (a *Auth) TableName() string {
	return TableName("uc_auth")
}

func AuthGetList(page, pageSize int, filters ...interface{}) ([]*Auth, int64) {
	offset := (page - 1) * pageSize
	list := make([]*Auth, 0)
	query := orm.NewOrm().QueryTable(new(Auth))
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

func AuthAdd(auth *Auth) (int64, error) {
	return orm.NewOrm().Insert(auth)
}

func AuthGetById(id int) (*Auth, error) {
	a := new(Auth)

	err := orm.NewOrm().QueryTable(new(Auth)).Filter("id", id).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Auth) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func AuthGetListForMenu(uid int, rids string) ([]*Auth, int64) {
	var nodes []*Auth
	var count int64
	var buf bytes.Buffer
	if uid == 1 {
		buf.WriteString("SELECT * FROM " + new(Auth).TableName() + " WHERE status = 1 order by pid,sort ")
		sql := buf.String()
		count, err := orm.NewOrm().Raw(sql).QueryRows(&nodes)
		if err != nil {
			return nil, count
		}
	} else {
		buf.WriteString("SELECT * FROM " + new(Auth).TableName() + " WHERE status = 1 AND (opened=0 OR id IN(SELECT auth_id FROM  ")
		buf.WriteString(new(RoleAuth).TableName())
		buf.WriteString(" WHERE role_id IN(?))) order by pid,sort ")
		sql := buf.String()
		count, err := orm.NewOrm().Raw(sql, rids).QueryRows(&nodes)
		if err != nil {
			return nil, count
		}
	}
	return nodes, count
}

func AuthGetListForRole(rid int) ([]*Auth, int64) {
	rids := strconv.Itoa(rid)
	var nodes []*Auth
	var count int64
	var buf bytes.Buffer
	buf.WriteString("SELECT * FROM " + new(Auth).TableName() + " WHERE status = 1 AND (id IN(SELECT auth_id FROM  ")
	buf.WriteString(new(RoleAuth).TableName())
	buf.WriteString(" WHERE role_id IN(?))) order by pid,sort ")
	sql := buf.String()
	count, err := orm.NewOrm().Raw(sql, rids).QueryRows(&nodes)
	if err != nil {
		return nil, count
	}
	return nodes, count
}
