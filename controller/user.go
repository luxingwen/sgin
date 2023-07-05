package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
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
// @Success 200 {object} model.User "Successfully created user"
// @Router /user/create [post]
func (uc *UserController) CreateUser(c *app.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := uc.Service.CreateUser(c, &user)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// GetUserByUUID gets a User by UUID.
// @Summary 获取用户信息
// @Description Get a user by its UUID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqUserQueryParam false "查询参数"
// @Success 200 {object} model.User "Successfully fetched user data"
// @Router /user/info [post]
func (uc *UserController) GetUserByUUID(c *app.Context) {
	param := &model.ReqUserQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := uc.Service.GetUserByUUID(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
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
// @Success 200 {object} model.User "Successfully updated user"
// @Router /user/update [post]
func (uc *UserController) UpdateUser(c *app.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := uc.Service.UpdateUser(c, &user)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// DeleteUser deletes a User by UUID.
// @Summary 删除用户
// @Description Delete a user by its UUID
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param params body model.ReqUserDeleteParam true "Delete user"
// @Success 200 {string} string "Successfully deleted user"
// @Router /user/delete [post]
func (uc *UserController) DeleteUser(c *app.Context) {
	params := &model.ReqUserDeleteParam{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := uc.Service.DeleteUser(c, params.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess("User deleted successfully")
}

// 获取用户列表
// @Summary 获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param params body model.ReqUserQueryParam true "获取用户列表参数"
// @Success 200 {object} model.PagedResponse
// @Router /user/list [post]
func (uc *UserController) GetUserList(c *app.Context) {
	param := &model.ReqUserQueryParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	users, err := uc.Service.GetUserList(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(users)
}

// 获取自己的信息
// @Summary 获取自己的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /user/myinfo [get]
func (uc *UserController) GetMyInfo(c *app.Context) {
	userId := c.GetString("user_id")
	user, err := uc.Service.GetUserByUUID(c, userId)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}
