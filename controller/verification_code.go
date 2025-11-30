package controller

import (
	"fmt"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/pkg/mail"
	"sgin/service"
)

type VerificationCodeController struct {
	VerificationCodeService *service.VerificationCodeService
}

var registerMailContent = `
<html>
<body>
    <h2>注册验证码</h2>
    <p>尊敬的用户，您的注册验证码为：<strong>%s</strong></p>
    <p>请在注册页面输入该验证码完成注册。</p>
    <p>请勿将验证码透露给他人。</p>
</body>
</html>
`

// CreateVerificationCode 创建验证码
// @Summary 创建验证码
// @Description 创建验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param params body model.ReqVerificationCodeParam true "验证码信息"
// @Success 200 {string} string "Successfully fetched user data"
// @Router /api/v1/verification_code/create [post]
func (v *VerificationCodeController) CreateVerificationCode(ctx *app.Context) {
	param := &model.ReqVerificationCodeParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind create verification code params failed")
		return
	}

	if param.Email == "" && param.Phone == "" {
		ctx.JSONErrLog(ecode.BadRequest("邮箱和手机号码不能同时为空"), "missing email and phone")
		return
	}

	code, err := v.VerificationCodeService.CreateVerificationCode(ctx, param.Email, param.Phone)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "create verification code failed")
		return
	}

	if param.Email != "" {
		// 发送邮件

		err = mail.Send(&mail.Options{
			MailHost: ctx.Config.MailConfig.Host,
			MailPort: ctx.Config.MailConfig.Port,
			MailUser: ctx.Config.MailConfig.Username,
			MailPass: ctx.Config.MailConfig.Password,
			MailTo:   param.Email,
			Subject:  ctx.Config.MailConfig.RegisterTile,
			Body:     fmt.Sprintf(registerMailContent, code),
		})
		if err != nil {
			ctx.JSONErrLog(ecode.InternalError(err.Error()), "send verification mail failed", "email", param.Email)
			return
		}
	}

	if param.Phone != "" {
		// 发送短信
	}

	ctx.JSONSuccess("验证码发送成功")
}

// CheckVerificationCode 检查验证码
// @Summary 检查验证码
// @Description 检查验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param params body string true "验证码信息"
// @Success 200 {string} string "Successfully fetched user data"
// @Router /api/v1/verification_code/check [post]
func (v *VerificationCodeController) CheckVerificationCode(ctx *app.Context) {
	param := &model.ReqVerificationCodeParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONErrLog(ecode.BadRequest(err.Error()), "bind check verification code params failed")
		return
	}

	if param.Email == "" && param.Phone == "" {
		ctx.JSONErrLog(ecode.BadRequest("邮箱和手机号码不能同时为空"), "missing email and phone")
		return
	}

	if param.Code == "" {
		ctx.JSONErrLog(ecode.BadRequest("验证码不能为空"), "code is empty")
		return
	}

	ok, err := v.VerificationCodeService.CheckVerificationCode(ctx, param.Code, param.Email, param.Phone)
	if err != nil {
		ctx.JSONErrLog(ecode.InternalError(err.Error()), "check verification code failed")
		return
	}

	if ok == false {
		ctx.JSONErrLog(ecode.BadRequest("验证码错误"), "verification code mismatch")
		return
	}

	ctx.JSONSuccess("验证码正确")
}
