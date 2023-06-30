package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"gorm.io/gorm"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(ctx *app.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	err := ctx.DB.Create(user).Error
	if err != nil {
		ctx.Logger.Error("Failed to create user", err)
		return errors.New("failed to create user")
	}
	return nil
}

func (s *UserService) GetUserByUUID(ctx *app.Context, uuid string) (*model.User, error) {
	user := &model.User{}
	err := ctx.DB.Where("uuid = ?", uuid).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		ctx.Logger.Error("Failed to get user by UUID", err)
		return nil, errors.New("failed to get user by UUID")
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx *app.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	err := ctx.DB.Save(user).Error
	if err != nil {
		ctx.Logger.Error("Failed to update user", err)
		return errors.New("failed to update user")
	}

	return nil
}

func (s *UserService) DeleteUser(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.User{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete user", err)
		return errors.New("failed to delete user")
	}

	return nil
}

// 根据用户名或邮箱获取用户
func (s *UserService) GetUserByUsernameOrEmail(ctx *app.Context, usernameOrEmail string) (*model.User, error) {
	user := &model.User{}
	err := ctx.DB.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		ctx.Logger.Error("Failed to get user by username or email", err)
		return nil, errors.New("failed to get user by username or email")
	}
	return user, nil
}
