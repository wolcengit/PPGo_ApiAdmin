/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-14 15:24:51
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:48:52
***********************************************/
package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type Role struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	RoleName   string `orm:"column(role_name);size(32)" json:"role_name"`         //角色名称
	Detail     string `orm:"column(detail);size(255)" json:"detail"`              //备注
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常 0禁用 2系统
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
}

func (a *Role) TableName() string {
	return TableName("uc_role")
}

func RoleGetList(page, pageSize int, filters ...interface{}) ([]*Role, int64) {
	offset := (page - 1) * pageSize
	list := make([]*Role, 0)
	query := orm.NewOrm().QueryTable(new(Role))
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

func RoleAdd(role *Role) (int64, error) {
	id, err := orm.NewOrm().Insert(role)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func RoleGetById(id int) (*Role, error) {
	r := new(Role)
	err := orm.NewOrm().QueryTable(new(Role)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Role) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(r, fields...); err != nil {
		return err
	}
	return nil
}

func RoleGetListForSelect(roleIds string) ([]map[string]interface{}, int64) {
	filters := make([]interface{}, 0)
	filters = append(filters, "status__gte", 1)
	result, total := RoleGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))

	role_ids := strings.Split(roleIds, ",")
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["role_name"] = v.RoleName
		row["checked"] = 0
		for i := 0; i < len(role_ids); i++ {
			role_id, _ := strconv.Atoi(role_ids[i])
			if role_id == v.Id {
				row["checked"] = 1
			}
		}
		list[k] = row
	}
	return list, total
}
