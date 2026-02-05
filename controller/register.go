package controller

import (
	"github.com/luxingwen/sgin/service"

	"github.com/luxingwen/sgin/model"
	"github.com/luxingwen/sgin/pkg/app"
	"github.com/luxingwen/sgin/pkg/ecode"
	"github.com/luxingwen/sgin/pkg/utils"

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
	// 获取系统配置 register_email_verify
	needVerify := false
	// 使用 map 接收结果，避免引入 shop-backend/model 导致循环依赖
	var configResult struct {
		Value string
	}
	// 注意：这里假设表名是 system_configs
	if err := c.DB.Table("system_configs").Select("value").Where("`key` = ?", "register_email_verify").Scan(&configResult).Error; err == nil {
		if configResult.Value == "true" {
			needVerify = true
		}
	}

	if needVerify {
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
	}

	// 获取密码加密密钥，如果为空则使用默认值（与登录验证保持一致）
	passwdKey := c.Config.PasswdKey
	if passwdKey == "" {
		passwdKey = "default-secret-key"
	}

	// 创建用户
	user := model.User{
		Uuid:     uuid.New().String(),
		Username: params.Username,
		Password: utils.HashPasswordWithSalt(params.Password, passwdKey),
		Email:    params.Email,
		Phone:    params.Phone,
	}

	err := rc.UserService.CreateUser(c, &user)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "create user failed", "username", user.Username, "email", user.Email)
		return
	}

	user.Password = ""
	c.JSONSuccess(user)
}
