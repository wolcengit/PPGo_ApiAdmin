/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-15 11:44:13
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:49:13
***********************************************/
package models

import (
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

func RoleAuthDeleteByRole(id int64) (int64, error) {
	query := orm.NewOrm().QueryTable(new(RoleAuth))
	return query.Filter("role_id", id).Delete()
}

func RoleAuthInsertMult(role_id int64, auths string) {
	RoleAuthDeleteByRole(role_id)
	authsSlice := strings.Split(auths, ",")
	for _, v := range authsSlice {
		aid, _ := strconv.Atoi(v)
		ra := new(RoleAuth)
		ra.AuthId = aid
		ra.RoleId = role_id
		RoleAuthAdd(ra)
	}
}
