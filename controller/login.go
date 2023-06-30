package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"sgin/service"
)

type LoginController struct {
	UserService *service.UserService
}

// 用户登录
// @Summary 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqUserLogin true "登录参数"
// @Success 200 {object} model.ResUserLogin
// @Router /login [post]
func (c *LoginController) Login(ctx *app.Context) {
	param := &model.ReqUserLogin{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.UserService.GetUserByUsernameOrEmail(ctx, param.Username)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	if utils.CheckPasswordHashWithSalt(param.Password, user.Password, ctx.Config.PasswdKey) == false {
		ctx.JSONError(http.StatusBadRequest, "用户名或密码错误")
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
}
