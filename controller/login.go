package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"sgin/service"

	"github.com/mileusna/useragent"
)

type LoginController struct {
	UserService        *service.UserService
	SysLoginLogService *service.SysLoginLogService
}

// 用户登录
// @Summary 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqUserLogin true "登录参数"
// @Success 200 {object} model.ResUserLogin
// @Router /api/v1/login [post]
func (c *LoginController) Login(ctx *app.Context) {
	param := &model.ReqUserLogin{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.UserService.GetUserByUsernameOrEmail(ctx, param.Username)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())

		c.CreateSysLoginLog(ctx, model.LoginStatusFail, param.Username, err.Error())
		return
	}

	if utils.CheckPasswordHashWithSalt(param.Password, user.Password, ctx.Config.PasswdKey) == false {
		ctx.JSONError(http.StatusBadRequest, "用户名或密码错误")
		c.CreateSysLoginLog(ctx, model.LoginStatusFail, param.Username, "密码错误")
		return
	}

	token, err := utils.GenerateToken(user.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(model.ResUserLogin{
		Token: token,
	})
	c.CreateSysLoginLog(ctx, model.LoginStatusSuccess, param.Username, "登录成功")
}

func (c *LoginController) CreateSysLoginLog(ctx *app.Context, status int, username string, msg string) {
	uaString := ctx.GetHeader("User-Agent")
	ua := useragent.Parse(uaString)
	sysLoginLog := model.SysLoginLog{
		Username:  username,
		Status:    status,
		Ip:        ctx.ClientIP(),
		UserAgent: uaString,
		Browser:   ua.Name + " " + ua.Version,
		Os:        ua.OS,
		Device:    ua.Device,
		Message:   msg,
	}
	if err := c.SysLoginLogService.CreateLoginLog(ctx, &sysLoginLog); err != nil {
		ctx.Logger.Error(err)
	}
}
