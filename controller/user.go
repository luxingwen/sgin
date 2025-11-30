package controller

import (
	"path/filepath"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/ecode"
	"sgin/service"
)

// UserController handles the operations related to User.
type UserController struct {
	Service *service.UserService
}

// CreateUser creates a new User.
// @Summary 创建用户
// @Description Create a new user with the input payload
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user body model.User true "Create user"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/create [post]
func (uc *UserController) CreateUser(c *app.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind create user failed")
		return
	}

	err := uc.Service.CreateUser(c, &user)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "create user failed", "username", user.Username, "email", user.Email)
		return
	}
	c.Logger.Infow("user created",
		"path", c.FullPath(),
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"user_uuid", user.Uuid,
		"username", user.Username,
	)
	c.JSONSuccess(user)
}

// GetUserByUUID gets a User by UUID.
// @Summary 获取用户信息
// @Description Get a user by its UUID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqUserQueryParam false "查询参数"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/info [post]
func (uc *UserController) GetUserByUUID(c *app.Context) {
	param := &model.ReqUserQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind get user param failed")
		return
	}

	user, err := uc.Service.GetUserByUUID(c, param.Uuid)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "get user by uuid failed", "uuid", param.Uuid)
		return
	}
	c.JSONSuccess(user)
}

// UpdateUser updates an existing User.
// @Summary 更新用户
// @Description Update a user with the input payload
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param user body model.User true "Update user"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/update [post]
func (uc *UserController) UpdateUser(c *app.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind update user failed")
		return
	}

	err := uc.Service.UpdateUser(c, &user)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "update user failed", "uuid", user.Uuid)
		return
	}
	c.Logger.Infow("user updated",
		"path", c.FullPath(),
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"user_uuid", user.Uuid,
	)
	c.JSONSuccess(user)
}

// DeleteUser deletes a User by UUID.
// @Summary 删除用户
// @Description Delete a user by its UUID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqUserDeleteParam true "Delete user"
// @Success 200 {object} app.Response
// @Router /api/v1/user/delete [post]
func (uc *UserController) DeleteUser(c *app.Context) {
	params := &model.ReqUserDeleteParam{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind delete user failed")
		return
	}

	if params.Uuid == c.GetString("user_id") {
		c.JSONErrLog(ecode.BadRequest("You can't delete yourself"), "self delete not allowed", "uuid", params.Uuid)
		return
	}

	err := uc.Service.DeleteUser(c, params.Uuid)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "delete user failed", "uuid", params.Uuid)
		return
	}
	c.Logger.Infow("user deleted",
		"path", c.FullPath(),
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"user_uuid", params.Uuid,
	)
	c.JSONSuccess("User deleted successfully")
}

// 获取用户列表
// @Summary 获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqUserQueryParam true "获取用户列表参数"
// @Success 200 {object} model.UserQueryResponse
// @Router /api/v1/user/list [post]
func (uc *UserController) GetUserList(c *app.Context) {
	param := &model.ReqUserQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "bind list users param failed")
		return
	}

	users, err := uc.Service.GetUserList(c, param)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "list users failed")
		return
	}

	c.JSONSuccess(users)
}

// 获取自己的信息
// @Summary 获取自己的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/myinfo [get]
func (uc *UserController) GetMyInfo(c *app.Context) {
	userId := c.GetString("user_id")
	user, err := uc.Service.GetUserByUUID(c, userId)
	if err != nil {
		c.JSONErr(ecode.InternalError(err.Error()))
		return
	}
	c.JSONSuccess(user)
}

// 修改头像
// @Summary 修改头像
// @Tags 用户
// @Accept json
// @Produce json
// @Param file formData file true "头像文件"
// @Success 200 {object} model.UserInfoResponse
// @Router /api/v1/user/avatar [post]
func (uc *UserController) UpdateAvatar(c *app.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSONErrLog(ecode.BadRequest(err.Error()), "parse avatar file failed")
		return
	}

	// 保存头像

	userid := c.GetString("user_id")
	if userid == "" {
		c.JSONErrLog(ecode.BadRequest("user_id is required"), "missing user_id in context")
		return
	}

	// 获取文件后缀

	extfile := filepath.Ext(file.Filename)

	filename := "/avatar/" + userid + extfile

	err = c.SaveUploadedFile(file, c.Config.Upload.Dir+filename)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "save avatar failed")
		return
	}

	user := model.User{
		Avatar: filename,
		Uuid:   userid,
	}

	err = uc.Service.UpdateUser(c, &user)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "update user avatar failed", "uuid", userid)
		return
	}
	c.Logger.Infow("user avatar updated",
		"path", c.FullPath(),
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"user_uuid", userid,
		"avatar", filename,
	)
	c.JSONSuccess(user)
}

// 获取所有用户

func (uc *UserController) GetAllUsers(c *app.Context) {
	users, err := uc.Service.GetAllUsers(c)
	if err != nil {
		c.JSONErrLog(ecode.InternalError(err.Error()), "get all users failed")
		return
	}
	c.JSONSuccess(users)
}
