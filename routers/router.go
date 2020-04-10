package routers

import (
	"github.com/astaxie/beego"
	"github.com/george518/PPGo_ApiAdmin/controllers"
)

func init() {

	beego.Router("/", &controllers.LoginController{}, "*:Login")

	beego.AutoRouter(&controllers.LoginController{})
	beego.AutoRouter(&controllers.HomeController{})

	beego.AutoRouter(&controllers.AuthController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.AdminController{})
	beego.AutoRouter(&controllers.UserController{})

	beego.AutoRouter(&controllers.XopProductController{})
	beego.AutoRouter(&controllers.XopDocumentController{})

	beego.AutoRouter(&controllers.XopModuleController{})
	beego.AutoRouter(&controllers.XopCategoryController{})
	beego.AutoRouter(&controllers.XopGroupController{})
	beego.AutoRouter(&controllers.XopFunctionController{})

	beego.AutoRouter(&controllers.BookCategoryController{})
	beego.AutoRouter(&controllers.BookLibraryController{})

}
