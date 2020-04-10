/**********************************************
** @Des: login
** @Author: haodaquan
** @Date:   2017-09-07 16:30:10
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:55:21
***********************************************/
package controllers

import (
	"errors"
	"fmt"
	"gopkg.in/ldap.v2"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/george518/PPGo_ApiAdmin/libs"
	"github.com/george518/PPGo_ApiAdmin/models"
	"github.com/george518/PPGo_ApiAdmin/utils"
	"github.com/patrickmn/go-cache"
)

type LoginController struct {
	BaseController
}

//登录 TODO:XSRF过滤
func (self *LoginController) Login() {
	if self.userId > 0 {
		self.redirect(beego.URLFor("HomeController.Index"))
	}
	list, _ := models.XopProductGetListForSelect()
	self.Data["Products"] = list
	beego.ReadFromRequest(&self.Controller)
	if self.isPost() {

		username := strings.TrimSpace(self.GetString("username"))
		password := strings.TrimSpace(self.GetString("password"))
		productId, _ := self.GetInt("prod_id")

		if username != "" && password != "" && productId > 0 {
			user, err := checkLogin(username, password)
			fmt.Println(user)
			flash := beego.NewFlash()
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			} else if user.Status == 0 {
				errorMsg = "该帐号已禁用"
			} else {
				if user.Id > 1 {
					if !models.XopProductAuthCheckForLogin(productId, user.Id) {
						errorMsg = "该产品没有权限"
					}
				}
				if errorMsg == "" {
					user.LastIp = self.getClientIp()
					user.LastProd = productId
					user.LastLogin = time.Now().Unix()
					user.Update()
					utils.Che.Set("uid"+strconv.Itoa(user.Id), user, cache.DefaultExpiration)
					authkey := libs.Md5([]byte(self.getClientIp() + "|" + user.Password + user.Salt))
					self.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)

					self.redirect(beego.URLFor("HomeController.Index"))
				}
			}
			flash.Error(errorMsg)
			flash.Store(&self.Controller)
			self.redirect(beego.URLFor("LoginController.Login"))
		}
	}
	self.TplName = "login/login.html"
}

//登出
func (self *LoginController) Logout() {
	self.Ctx.SetCookie("auth", "")
	self.redirect(beego.URLFor("LoginController.Login"))
}

func checkLogin(username string, password string) (*models.Admin, error) {
	if username == "" || password == "" {
		return nil, errors.New("用户密码错误")
	}
	domain := beego.AppConfig.String("ldap.domain")
	if !strings.HasSuffix(username, domain) {
		// Local
		user, err := models.AdminGetByName(username)
		if err != nil {
			return nil, errors.New("用户密码错误")
		}
		if user.Password != libs.Md5([]byte(password+user.Salt)) {
			return nil, errors.New("用户密码错误")
		}
		return user, nil
	}

	//LDAP
	if beego.AppConfig.DefaultBool("ldap.enable", false) == false {
		return nil, errors.New("用户密码错误!")
	}

	lc, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", beego.AppConfig.String("ldap.host"), beego.AppConfig.DefaultInt("ldap.port", 3268)))
	if err != nil {
		return nil, errors.New("绑定 LDAP 用户失败")
	}
	defer lc.Close()
	err = lc.Bind(beego.AppConfig.String("ldap.user"), beego.AppConfig.String("ldap.password"))
	if err != nil {
		return nil, errors.New("绑定 LDAP 用户失败")
	}
	searchRequest := ldap.NewSearchRequest(
		beego.AppConfig.String("ldap.base"),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		//修改objectClass通过配置文件获取值
		fmt.Sprintf("(&(%s)(%s=%s))", beego.AppConfig.String("ldap.filter"), beego.AppConfig.String("ldap.attribute"), username),
		[]string{"dn", "mail", "displayName", "telephoneNumber"},
		nil,
	)
	searchResult, err := lc.Search(searchRequest)
	if err != nil {
		return nil, errors.New("绑定 LDAP 用户失败")
	}
	if len(searchResult.Entries) != 1 {
		return nil, errors.New("LDAP用户不存在或者多于一个")
	}
	userdn := searchResult.Entries[0].DN
	err = lc.Bind(userdn, password)
	if err != nil {
		return nil, errors.New("用户密码错误")
	}
	user, err := models.AdminGetByName(username)
	if err != nil {
		user = new(models.Admin)
		user.LoginName = username
		user.RealName = searchResult.Entries[0].GetAttributeValue("displayName")
		user.Salt = "XyZ0"
		user.Password = "!@#!$%^$^#"
		user.RoleIds = ""
		user.Phone = searchResult.Entries[0].GetAttributeValue("telephoneNumber")
		user.Email = searchResult.Entries[0].GetAttributeValue("mail")
		user.Status = 1
		user.CreateId = 1
		models.AdminAdd(user)
	} else {
		user.RealName = searchResult.Entries[0].GetAttributeValue("displayName")
		user.Phone = searchResult.Entries[0].GetAttributeValue("telephoneNumber")
		user.Email = searchResult.Entries[0].GetAttributeValue("mail")
		user.Update()
	}
	user, err = models.AdminGetByName(username)
	return user, err
}
