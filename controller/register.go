package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"
	"sgin/service"

	"github.com/google/uuid"
)

type RegisterController struct {
	UserService             *service.UserService
	VerificationCodeService *service.VerificationCodeService
}

// @Tags Register
// @Summary 注册
// @Description 注册
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqRegisterParam true "Register"
// @Success 200 {object} string "Successfully fetched user data"
// @Router /register [post]
func (rc *RegisterController) Register(c *app.Context) {

	params := &model.ReqRegisterParam{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 验证验证码
	ok, err := rc.VerificationCodeService.CheckVerificationCode(c, params.Code, params.Email, params.Phone)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	if ok == false {
		c.JSONError(http.StatusBadRequest, "验证码错误")
		return
	}
	// 更新验证码状态
	err = rc.VerificationCodeService.UpdateVerificationCode(c, params.Code, params.Email, params.Phone)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
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
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = ""
	c.JSONSuccess(user)
}
