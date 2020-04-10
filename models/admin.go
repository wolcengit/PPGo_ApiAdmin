/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-16 15:42:43
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:48:17
***********************************************/
package models

import (
	"github.com/astaxie/beego/orm"
)

type Admin struct {
	Id         int    `orm:"column(id);pk;auto;unique" json:"id"`
	LoginName  string `orm:"column(login_name);size(128)" json:"login_name"`      //用户名
	RealName   string `orm:"column(real_name);size(32)" json:"real_name"`         //真实姓名
	Password   string `orm:"column(password);size(32)" json:"password"`           //密码
	RoleIds    string `orm:"column(role_ids);size(254)" json:"role_ids"`          //角色id字符串，如：2,3,4
	Phone      string `orm:"column(phone);size(32)" json:"phone"`                 //手机号码
	Email      string `orm:"column(email);size(128)" json:"email"`                //邮箱
	Salt       string `orm:"column(salt);size(10)" json:"salt"`                   //密码盐
	LastLogin  int64  `orm:"column(last_login);type(bigint)" json:"last_login"`   //最后登录时间
	LastIp     string `orm:"column(last_ip);size(16)" json:"last_ip"`             //最后登录IP
	LastProd   int    `orm:"column(last_prod);type(int)" json:"last_prod"`        //最后登录产品
	Status     int    `orm:"column(status);type(int)" json:"status"`              //状态：1-正常，0-禁用
	Appkey     string `orm:"column(appkey);size(16)" json:"appkey"`               //APP1访问码
	CreateId   int    `orm:"column(create_id);type(int)" json:"create_id"`        //创建者
	UpdateId   int    `orm:"column(update_id);type(int)" json:"update_id"`        //修改者
	CreateTime int64  `orm:"column(create_time);type(bigint)" json:"create_time"` //创建时间
	UpdateTime int64  `orm:"column(update_time);type(bigint)" json:"update_time"` //修改时间
}

func (a *Admin) TableName() string {
	return TableName("uc_admin")
}

func AdminAdd(a *Admin) (int64, error) {
	return orm.NewOrm().Insert(a)
}

func AdminGetByName(loginName string) (*Admin, error) {
	a := new(Admin)
	err := orm.NewOrm().QueryTable(new(Admin)).Filter("login_name", loginName).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func AdminGetList(page, pageSize int, filters ...interface{}) ([]*Admin, int64) {
	offset := (page - 1) * pageSize
	list := make([]*Admin, 0)
	query := orm.NewOrm().QueryTable(new(Admin))
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

func AdminGetById(id int) (*Admin, error) {
	r := new(Admin)
	err := orm.NewOrm().QueryTable(new(Admin)).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a *Admin) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}

func AdminGetByAppkey(appkey string) (*Admin, error) {
	r := new(Admin)
	err := orm.NewOrm().QueryTable(new(Admin)).Filter("appkey", appkey).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
