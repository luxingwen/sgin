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

// 获取用户列表
func (s *UserService) GetUserList(ctx *app.Context, params *model.ReqUserQueryParam) (r *model.PagedResponse, err error) {
	var users []*model.User
	var total int64
	db := ctx.DB.Model(&model.User{})
	if params.Username != "" {
		db = db.Where("username = ?", params.Username)
	}
	if params.Email != "" {
		db = db.Where("email = ?", params.Email)
	}
	if params.Phone != "" {
		db = db.Where("phone = ?", params.Phone)
	}
	if params.Status != 0 {
		db = db.Where("status = ?", params.Status)
	}

	if params.StartTime != "" {
		db = db.Where("created_at >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		db = db.Where("created_at <= ?", params.EndTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get user list", err)
		return nil, errors.New("failed to get user list")
	}
	err = db.Offset(params.GetOffset()).Limit(params.PageSize).Find(&users).Error
	if err != nil {
		ctx.Logger.Error("Failed to get user list", err)
		return nil, errors.New("failed to get user list")
	}
	return &model.PagedResponse{
		Total: total,
		Data:  users,
	}, nil
}
