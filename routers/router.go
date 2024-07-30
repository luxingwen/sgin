package routers

import (
	"sgin/controller"
	"sgin/middleware"
	"sgin/pkg/app"

	"io/ioutil"
	"net/http"
	"sgin/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(ctx *app.App) {
	InitSwaggerRouter(ctx)
	InitUserRouter(ctx)
	InitMenuRouter(ctx)
	InitAppRouter(ctx)
	InitVerificationCodeRouter(ctx)
	InitRegisterRouter(ctx)
	InitLoginRouter(ctx)
	InitServerRouter(ctx)
	InitTeamRouter(ctx)
	InitSysLoginLogRouter(ctx)
	InitSysOpLogRouter(ctx)
	InitSysApiRouter(ctx)
	InitPermissionRouter(ctx)
	InitPermissionMenuRouter(ctx)
	InitPermissionUserRouter(ctx)
	InitMenuAPIRouter(ctx)
	InitTeamMemberRouter(ctx)
}

func InitUserRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		userController := &controller.UserController{
			Service: &service.UserService{},
		}

		v1.POST("/user/create", userController.CreateUser)
		v1.POST("/user/info", userController.GetUserByUUID)
		v1.POST("/user/list", userController.GetUserList)
		v1.POST("/user/update", userController.UpdateUser)
		v1.POST("/user/delete", userController.DeleteUser)
		v1.GET("/user/myinfo", userController.GetMyInfo)
		v1.POST("/user/avatar", userController.UpdateAvatar)
		v1.POST("/user/all", userController.GetAllUsers)

	}

	{
		roleController := &controller.RoleController{
			RoleService: &service.RoleService{},
		}

		v1.POST("/role/create", roleController.CreateRole)
		v1.POST("/role/list", roleController.GetRoleList)
		v1.POST("/role/update", roleController.UpdateRole)
		v1.POST("/role/delete", roleController.DeleteRole)

	}
}

func InitMenuRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		menuController := &controller.MenuController{
			MenuService: &service.MenuService{},
		}
		v1.POST("/menu/create", menuController.CreateMenu)
		v1.POST("/menu/list", menuController.GetMenuList)
		v1.POST("/menu/update", menuController.UpdateMenu)
		v1.POST("/menu/delete", menuController.DeleteMenu)
		v1.POST("/menu/info", menuController.GetMenuInfo)
	}
}

func InitAppRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		appController := &controller.AppController{
			AppService: &service.AppService{},
		}
		v1.POST("/app/list", appController.GetAppList)
		v1.POST("/app/create", appController.CreateApp)
		v1.POST("/app/update", appController.UpdateApp)
		v1.POST("/app/delete", appController.DeleteApp)

	}
}

func InitVerificationCodeRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	{
		verificationCodeController := &controller.VerificationCodeController{
			VerificationCodeService: &service.VerificationCodeService{},
		}
		v1.POST("/verification_code/create", verificationCodeController.CreateVerificationCode)
	}
}

// 注册的路由
func InitRegisterRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	{
		registerController := &controller.RegisterController{
			UserService:             &service.UserService{},
			VerificationCodeService: &service.VerificationCodeService{},
		}
		v1.POST("/register", registerController.Register)
	}
}

func InitLoginRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	{
		loginController := &controller.LoginController{
			UserService: &service.UserService{},
		}
		v1.POST("/login", loginController.Login)
	}
}

// 服务的路由
func InitServerRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		serverController := &controller.ServerController{
			ServerService: &service.ServerService{},
		}
		v1.POST("/server/create", serverController.CreateServer)
		v1.POST("/server/update", serverController.UpdateServer)
		v1.POST("/server/delete", serverController.DeleteServer)
		v1.POST("/server/info", serverController.GetServerInfo)
		v1.POST("/server/list", serverController.GetServerList)
	}
}

// 团队的路由
func InitTeamRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		teamController := &controller.TeamController{
			TeamService: &service.TeamService{},
		}
		v1.POST("/team/create", teamController.CreateTeam)
		v1.POST("/team/update", teamController.UpdateTeam)
		v1.POST("/team/delete", teamController.DeleteTeam)
		v1.POST("/team/info", teamController.GetTeamInfo)
		v1.POST("/team/list", teamController.GetTeamList)
	}
}

func InitTeamMemberRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		teamMemberController := &controller.TeamMemberController{
			TeamMemberService: &service.TeamMemberService{},
		}
		v1.POST("/team_member/create", teamMemberController.CreateTeamMember)
		v1.POST("/team_member/delete", teamMemberController.DeleteTeamMember)
		v1.POST("/team_member/list", teamMemberController.GetTeamMemberList)
	}
}

func InitSysApiRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		apiController := &controller.APIController{
			APIService: &service.APIService{},
		}
		v1.POST("/sys_api/create", apiController.CreateAPI)
		v1.POST("/sys_api/update", apiController.UpdateAPI)
		v1.POST("/sys_api/delete", apiController.DeleteAPI)
		v1.POST("/sys_api/list", apiController.GetAPIList)
		v1.POST("/sys_api/info", apiController.GetAPIInfo)

	}
}

func InitSysOpLogRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		sysOpLogController := &controller.SysOpLogController{
			SysOpLogService: &service.SysOpLogService{},
		}

		v1.POST("/sysoplog/delete", sysOpLogController.DeleteSysOpLog)
		v1.POST("/sysoplog/info", sysOpLogController.GetSysOpLogInfo)
		v1.POST("/sysoplog/list", sysOpLogController.GetSysOpLogList)
	}
}

func InitSysLoginLogRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		sysLoginLogController := &controller.SysLoginLogController{
			LoginLogService: &service.SysLoginLogService{},
		}

		v1.POST("/sys_login_log/info", sysLoginLogController.GetLoginLog)
		v1.POST("/sys_login_log/list", sysLoginLogController.GetLoginLogList)
	}
}

func InitPermissionRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		permissionController := &controller.PermissionController{
			PermissionService: &service.PermissionService{},
		}
		v1.POST("/permission/create", permissionController.CreatePermission)
		v1.POST("/permission/update", permissionController.UpdatePermission)
		v1.POST("/permission/delete", permissionController.DeletePermission)
		v1.POST("/permission/info", permissionController.GetPermissionInfo)
		v1.POST("/permission/list", permissionController.GetPermissionList)
	}
}

func InitPermissionMenuRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		permissionMenuController := &controller.PermissionMenuController{
			PermissionMenuService: &service.PermissionMenuService{},
		}
		v1.POST("/permission_menu/create", permissionMenuController.CreatePermissionMenu)
		v1.POST("/permission_menu/update", permissionMenuController.UpdatePermissionMenu)
		v1.POST("/permission_menu/delete", permissionMenuController.DeletePermissionMenu)
		v1.POST("/permission_menu/info", permissionMenuController.GetPermissionMenuInfo)
		v1.POST("/permission_menu/info_menu", permissionMenuController.GetPermissionMenuListByPermissionUUID)
		v1.POST("/permission_menu/list", permissionMenuController.GetPermissionMenuList)
	}
}

func InitPermissionUserRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		permissionUserController := &controller.UserPermissionController{
			UserPermissionService: &service.UserPermissionService{},
		}
		v1.POST("/permission_user/create", permissionUserController.CreateUserPermission)
		v1.POST("/permission_user/update", permissionUserController.UpdateUserPermission)
		v1.POST("/permission_user/delete", permissionUserController.DeleteUserPermission)
		v1.POST("/permission_user/info", permissionUserController.GetUserPermissionInfo)
		v1.POST("/permission_user/list", permissionUserController.GetUserPermissionList)
	}
}

func InitMenuAPIRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		menuAPIController := &controller.MenuAPIController{
			MenuAPIService: &service.MenuAPIService{},
		}
		v1.POST("/menu_api/create", menuAPIController.CreateMenuAPI)
		v1.POST("/menu_api/update", menuAPIController.UpdateMenuAPI)
		v1.POST("/menu_api/delete", menuAPIController.DeleteMenuAPI)
		v1.POST("/menu_api/info", menuAPIController.GetMenuAPIInfo)
		v1.POST("/menu_api/info_menu", menuAPIController.GetMenuAPIListByMenuUUID)
		v1.POST("/menu_api/info_api", menuAPIController.GetMenuAPIListByAPIUUID)
		v1.POST("/menu_api/list", menuAPIController.GetMenuAPIList)
	}
}

func InitSwaggerRouter(ctx *app.App) {
	ctx.GET("/swagger/doc.json", func(c *app.Context) {
		jsonFile, err := ioutil.ReadFile("./docs/swagger.json") // Replace with your actual json file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", jsonFile)
	})

	ctx.GET("/swagger/redoc.standalone.js", func(c *app.Context) {
		b, err := ioutil.ReadFile("./swagger/redoc.standalone.js") // Replace with your actual json file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})

	ctx.GET("/swagger/index.html", func(c *app.Context) {
		b, err := ioutil.ReadFile("./swagger/swagger.html") // Replace with your actual json file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})
}
