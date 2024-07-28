package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"sgin/model"
	"sgin/pkg/app"
	"sgin/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(ctx *app.Context, user *model.User) error {
	user.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	user.UpdatedAt = user.CreatedAt
	user.Uuid = uuid.New().String()

	if user.Password == "" {
		user.Password = "123456"
	}

	user.Password = utils.HashPasswordWithSalt(user.Password, ctx.Config.PasswdKey)

	err := ctx.DB.Create(user).Error
	if err != nil {
		ctx.Logger.Error("Failed to create user", err)

		if strings.Contains(err.Error(), fmt.Sprintf("Duplicate entry '%s' for key", user.Username)) {
			return errors.New(user.Username + "用户名已存在")
		}
		if strings.Contains(err.Error(), fmt.Sprintf("Duplicate entry '%s' for key", user.Email)) {
			return errors.New(user.Email + "邮箱已存在")
		}

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
	user.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	if user.Password != "" {
		user.Password = utils.HashPasswordWithSalt(user.Password, ctx.Config.PasswdKey)
	}

	err := ctx.DB.Where("uuid = ?", user.Uuid).Updates(user).Error
	if err != nil {
		ctx.Logger.Error("Failed to update user:", err)
		return errors.New("failed to update user")
	}

	return nil
}

func (s *UserService) DeleteUser(ctx *app.Context, uuid string) error {
	err := ctx.DB.Model(&model.User{}).Where("uuid = ?", uuid).Update("is_deleted", 1).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete user", err)
		return errors.New("failed to delete user")
	}

	return nil
}

// 获取所有可用用户
func (s *UserService) GetAllUsers(ctx *app.Context) ([]*model.User, error) {
	users := make([]*model.User, 0)
	err := ctx.DB.Where("is_deleted = ?", 0).Find(&users).Error
	if err != nil {
		ctx.Logger.Error("Failed to get all users", err)
		return nil, errors.New("failed to get all users")
	}

	return users, nil
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
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", params.Username))
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

	db = db.Where("is_deleted = ?", 0)

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

	for _, user := range users {
		user.Password = ""
	}

	return &model.PagedResponse{
		Total:    total,
		Data:     users,
		Current:  params.Current,
		PageSize: params.PageSize,
	}, nil
}

// 根据用户UUID列表获取用户列表
func (s *UserService) GetUsersByUUIDs(ctx *app.Context, uuids []string) (map[string]*model.User, error) {
	users := make([]*model.User, 0)
	err := ctx.DB.Where("uuid IN (?)", uuids).Find(&users).Error
	if err != nil {
		ctx.Logger.Error("Failed to get users by UUIDs", err)
		return nil, errors.New("failed to get users by UUIDs")
	}

	userMap := make(map[string]*model.User)
	for _, user := range users {
		userMap[user.Uuid] = user
	}

	return userMap, nil
}
