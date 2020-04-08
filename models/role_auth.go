/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-15 11:44:13
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:49:13
***********************************************/
package models

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

type RoleAuth struct {
	Id     int   `orm:"column(id);pk;auto;unique" json:"id"`
	AuthId int   `orm:"column(auth_id);type(int)" json:"auth_id"`
	RoleId int64 `orm:"column(role_id);type(int)" json:"role_id"`
}

func (ra *RoleAuth) TableName() string {
	return TableName("uc_role_auth")
}

func RoleAuthAdd(ra *RoleAuth) (int64, error) {
	return orm.NewOrm().Insert(ra)
}

func RoleAuthGetById(id int) ([]*RoleAuth, error) {
	list := make([]*RoleAuth, 0)
	query := orm.NewOrm().QueryTable(new(RoleAuth))
	_, err := query.Filter("role_id", id).All(&list, "AuthId")
	if err != nil {
		return nil, err
	}
	return list, nil
}

func RoleAuthDelete(id int) (int64, error) {
	query := orm.NewOrm().QueryTable(new(RoleAuth))
	return query.Filter("role_id", id).Delete()
}

//获取多个
func RoleAuthGetByIds(RoleIds string) (Authids string, err error) {
	list := make([]*RoleAuth, 0)
	query := orm.NewOrm().QueryTable(new(RoleAuth))
	ids := strings.Split(RoleIds, ",")
	_, err = query.Filter("role_id__in", ids).All(&list, "AuthId")
	if err != nil {
		return "", err
	}
	b := bytes.Buffer{}
	for _, v := range list {
		if v.AuthId != 0 && v.AuthId != 1 {
			b.WriteString(strconv.Itoa(v.AuthId))
			b.WriteString(",")
		}
	}
	Authids = strings.TrimRight(b.String(), ",")
	return Authids, nil
}
