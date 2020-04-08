/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-08 00:18:02
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-16 17:26:48
***********************************************/

package models

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	// fmt.Println(dsn)

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(Auth), new(Role), new(RoleAuth), new(Admin))

	// 自动创建表 参数二为是否开启创建表   参数三是否更新表
	orm.RunSyncdb("default", false, true)

	a, err := AdminGetByName("admin")
	if err != nil || a == nil {
		// admin == admin
		AdminAdd(&Admin{Id: 1, LoginName: "admin", RealName: "超级管理员", Password: "7fa7c707f8dfa222772cd71735c01e51", RoleIds: "0", Salt: "Gfm3", Status: 1, CreateId: 1})
	}

	t, err := AuthGetById(1)
	if err != nil || t == nil {
		AuthAdd(&Auth{Id: 1, Pid: 0, AuthName: "所有权限", AuthUrl: "/", Sort: 1, Icon: "", IsShow: 1, CreateId: 1, Status: 1})

		AuthAdd(&Auth{Id: 2, Pid: 1, AuthName: "权限管理", AuthUrl: "/", Sort: 999, Icon: "fa-id-card", IsShow: 1, CreateId: 1, Status: 1})

		AuthAdd(&Auth{Id: 3, Pid: 2, AuthName: "用户管理", AuthUrl: "/admin/list", Sort: 1, Icon: "fa-user-o", IsShow: 1, CreateId: 1, Status: 1})
		AuthAdd(&Auth{Id: 4, Pid: 3, AuthName: "新增", AuthUrl: "/admin/add", Sort: 1, Icon: "", IsShow: 1, CreateId: 1, Status: 1})
		AuthAdd(&Auth{Id: 5, Pid: 3, AuthName: "修改", AuthUrl: "/admin/edit", Sort: 2, Icon: "", IsShow: 1, CreateId: 1, Status: 1})
		AuthAdd(&Auth{Id: 6, Pid: 3, AuthName: "删除", AuthUrl: "/admin/ajaxdel", Sort: 3, Icon: "", IsShow: 1, CreateId: 1, Status: 1})

		AuthAdd(&Auth{Id: 7, Pid: 2, AuthName: "角色管理", AuthUrl: "/role/list", Sort: 2, Icon: "fa-user-circle-o", IsShow: 1, CreateId: 1, Status: 1})
		AuthAdd(&Auth{Id: 8, Pid: 7, AuthName: "新增", AuthUrl: "/role/add", Sort: 1, Icon: "", IsShow: 1, CreateId: 1, Status: 1})
		AuthAdd(&Auth{Id: 9, Pid: 7, AuthName: "修改", AuthUrl: "/role/edit", Sort: 2, Icon: "", IsShow: 1, CreateId: 1, Status: 1})
		AuthAdd(&Auth{Id: 10, Pid: 7, AuthName: "删除", AuthUrl: "/role/ajaxdel", Sort: 3, Icon: "", IsShow: 1, CreateId: 1, Status: 1})

		AuthAdd(&Auth{Id: 11, Pid: 1, AuthName: "个人中心", AuthUrl: "profile/edit", Sort: 1001, Icon: "fa-user-circle-o", IsShow: 1, CreateId: 1, Status: 1})
		AuthAdd(&Auth{Id: 12, Pid: 11, AuthName: "资料修改", AuthUrl: "/user/edit", Sort: 1, Icon: "fa-edit", IsShow: 1, CreateId: 1, Status: 1})

	}

}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
