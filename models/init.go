/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-08 00:18:02
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-16 17:26:48
***********************************************/

package models

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpass := beego.AppConfig.String("db.pass")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	fmt.Println(dsn)

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(Admin), new(Auth), new(Role), new(RoleAuth),
		new(XopProduct), new(XopProductAuth),
		new(XopDocument), new(XopModule), new(XopCategory), new(XopGroup), new(XopFunction),
		new(BookCategory), new(BookLibrary), new(BookDetail), new(BookReader))

	// 自动创建表 参数二为是否开启创建表   参数三是否更新表
	orm.RunSyncdb("default", false, true)

	u, err := AdminGetByName("admin")
	if err != nil || u == nil {
		// admin == admin
		AdminAdd(&Admin{Id: 1, LoginName: "admin", RealName: "超级管理员", Password: "7fa7c707f8dfa222772cd71735c01e51", RoleIds: "0", Salt: "Gfm3", Status: 1, CreateId: 1})
	}

	a, err := AuthGetById(1)
	if err != nil || a == nil {
		initAuth(1, 0, "所有权限", "/", 1, "", 0, 1)

		initAuth(2, 1, "个人中心", "/", 1, "fa-user-circle-o", 1, 0)
		initAuth(3, 2, "我的书籍", "/user/books", 1, "fa-tree", 1, 0)
		initAuth(4, 2, "我的资料", "/user/edit", 2, "fa-edit", 1, 0)

		initAuth(10, 1, "用户权限", "/", 999, "fa-id-card", 1, 2)
		initAuthPart(11, 10, "用户管理", "/admin", 1, "fa-user-o", 1, 2)
		initAuthPart(15, 10, "角色管理", "/role", 2, "fa-user-circle-o", 1, 2)

		initAuth(20, 1, "XOP平台", "/", 998, "fa-cloud", 1, 1)
		initAuthPart(21, 20, "XOP产品", "/xopproduct", 1, "fa-tree", 1, 1)

		initAuth(30, 1, "XOP文档", "/", 12, "fa-medium", 1, 1)
		initAuthPart(31, 30, "XOP文档", "/xopdocument", 2, "fa-wikipedia-w", 1, 1)

		initAuth(40, 1, "XOP接口", "/", 11, "fa-anchor", 1, 1)
		initAuthPart(53, 40, "XOP函数", "/xopfunction", 1, "fa-plug", 1, 1)
		initAuthPart(49, 40, "XOP分组", "/xopgroup", 2, "fa-object-group", 1, 1)
		initAuthPart(45, 40, "XOP类别", "/xopcategory", 3, "fa-th-large", 1, 1)
		initAuthPart(41, 40, "XOP模块", "/xopmodule", 4, "fa-puzzle-piece", 1, 1)

		initAuth(60, 1, "XOP书籍", "/", 13, "fa-book", 1, 1)
		initAuthPart(65, 60, "XOP书籍", "/booklibrary", 1, "fa-book", 1, 1)
		initAuthPart(61, 60, "书籍类别", "/bookcategory", 2, "fa-th-large", 1, 1)

	}

	r, err := RoleGetById(1)
	if err != nil || r == nil {
		initRole(1, "系统管理员", 2)
		initRole(2, "XOP平台管理员", 2)
		initRole(3, "XOP文档管理员", 2)
		initRole(4, "XOP接口管理员", 2)
		initRole(5, "XOP函数管理员", 2)
		initRole(6, "XOP书籍管理员", 2)

		RoleAuthInsertMult(1, "1,10,11,11,13,14,15,16,17,18,")
		RoleAuthInsertMult(2, "1,20,21,22,23,24,25,26,27,28,")
		RoleAuthInsertMult(3, "1,30,31,32,33,34,35,36,37,38,")
		RoleAuthInsertMult(4, "1,40,41,42,43,44,45,46,47,48,49,50,51,52,")
		RoleAuthInsertMult(5, "1,40,53,54,55,56,")
		RoleAuthInsertMult(6, "1,60,61,62,63,64,65,66,67,68,")

	}

	b, err := XopProductGetById(1)
	if err != nil || b == nil {
		XopProductAdd(&XopProduct{Id: 1, Code: "XOP", Name: "XOP平台", Detail: "", CreateId: 1, Status: 1})
		XopDocumentAdd(&XopDocument{Id: 1, Pid: 0, ProdId: 1, Name: "XOP平台说明", Detail: "", Sort: 1, CreateId: 1, Status: 1})
	}

}

func initAuth(id int, pid int, name string, url string, sort int, icon string, show int, open int) {
	AuthAdd(&Auth{Id: id, Pid: pid, AuthName: name, AuthUrl: url, Sort: sort, Icon: icon, IsShow: show, CreateId: 1, Opened: open, Status: 1})
}
func initAuthPart(id int, pid int, name string, url string, sort int, icon string, show int, open int) {
	AuthAdd(&Auth{Id: id, Pid: pid, AuthName: name, AuthUrl: url + "/list", Sort: sort, Icon: icon, IsShow: show, CreateId: 1, Opened: open, Status: 1})
	AuthAdd(&Auth{Id: id + 1, Pid: id, AuthName: "新增", AuthUrl: url + "/add", Sort: 1, Icon: "fa", IsShow: 0, CreateId: 1, Opened: open, Status: 1})
	AuthAdd(&Auth{Id: id + 2, Pid: id, AuthName: "修改", AuthUrl: url + "/edit", Sort: 2, Icon: "fa", IsShow: 0, CreateId: 1, Opened: open, Status: 1})
	AuthAdd(&Auth{Id: id + 3, Pid: id, AuthName: "删除", AuthUrl: url + "/ajaxdel", Sort: 3, Icon: "fa", IsShow: 0, CreateId: 1, Opened: open, Status: 1})
}
func initRole(id int, name string, status int) {
	RoleAdd(&Role{Id: id, RoleName: name, Status: status})
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
