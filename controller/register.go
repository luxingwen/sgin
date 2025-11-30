package controller

import (
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/pkg/utils"
	"sgin/service"

	"github.com/google/uuid"
)

type RegisterController struct {
	UserService             *service.UserService
	VerificationCodeService *service.VerificationCodeService
}

// @Summary 注册
// @Description 注册
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqRegisterParam true "Register"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/register [post]
func (rc *RegisterController) Register(c *app.Context) {

	params := &model.ReqRegisterParam{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind register params failed")
		return
	}

	// 验证验证码
	ok, err := rc.VerificationCodeService.CheckVerificationCode(c, params.Code, params.Email, params.Phone)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "check verification code failed")
		return
	}

	if ok == false {
		c.JSONErrLog(ecode.BadRequest("验证码错误"), "verification code mismatch", "email", params.Email, "phone", params.Phone)
		return
	}
	// 更新验证码状态
	err = rc.VerificationCodeService.UpdateVerificationCode(c, params.Code, params.Email, params.Phone)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "update verification code failed")
		return
	}

	// 创建用户
	user := model.User{
		Uuid:     uuid.New().String(),
		Username: params.Username,
		Password: utils.HashPasswordWithSalt(params.Password, c.Config.PasswdKey),
		Email:    params.Email,
		Phone:    params.Phone,
	}

	err = rc.UserService.CreateUser(c, &user)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "create user failed", "username", user.Username, "email", user.Email)
		return
	}

	user.Password = ""
	c.JSONSuccess(user)
}
