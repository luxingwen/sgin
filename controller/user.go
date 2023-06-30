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
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body model.User true "Create user"
// @Success 200 {object} model.User "Successfully created user"
// @Router /users [post]
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
// @Summary Get a user by UUID
// @Description Get a user by its UUID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param uuid path string true "User's UUID"
// @Success 200 {object} model.User "Successfully fetched user data"
// @Router /users/{uuid} [get]
func (uc *UserController) GetUserByUUID(c *app.Context) {
	uuid := c.Param("uuid")

	user, err := uc.Service.GetUserByUUID(c, uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(user)
}

// UpdateUser updates an existing User.
// @Summary Update a user
// @Description Update a user with the input payload
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body model.User true "Update user"
// @Success 200 {object} model.User "Successfully updated user"
// @Router /users [put]
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
// @Summary Delete a user by UUID
// @Description Delete a user by its UUID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param uuid path string true "User's UUID"
// @Success 200 {string} string "Successfully deleted user"
// @Router /users/{uuid} [delete]
func (uc *UserController) DeleteUser(c *app.Context) {
	uuid := c.Param("uuid")

	err := uc.Service.DeleteUser(c, uuid)
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
