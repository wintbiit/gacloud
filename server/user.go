package server

import (
	"context"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"regexp"
)

func (s *GaCloudServer) AddUser(ctx context.Context, username, password, email string) (*model.User, error) {
	password = utils.Sha256SumBytes([]byte(password))

	user := &model.User{
		Name:     username,
		Password: password,
		Email:    email,
	}
	err := s.db.WithContext(ctx).Create(user).Error
	return user, err
}

func (s *GaCloudServer) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)
	err := s.db.WithContext(ctx).Where("name = ?", username).First(user).Error
	return user, err
}

func (s *GaCloudServer) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := s.db.WithContext(ctx).Where("email = ?", email).First(user).Error
	return user, err
}

func (s *GaCloudServer) GetUser(ctx context.Context, id uint) (*model.User, error) {
	user := new(model.User)
	err := s.db.WithContext(ctx).First(user, id).Error
	return user, err
}

func (s *GaCloudServer) UpdateUser(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *GaCloudServer) DeleteUser(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (s *GaCloudServer) UserExists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (s *GaCloudServer) UserLogin(ctx context.Context, username, password string) (*model.User, error) {
	if emailPattern.MatchString(username) {
		return s.userLoginByEmail(ctx, username, password)
	}
	return s.userLoginByUsername(ctx, username, password)
}

func (s *GaCloudServer) userLoginByUsername(ctx context.Context, username, password string) (*model.User, error) {
	password = utils.Sha256SumBytes([]byte(password))

	user := new(model.User)
	err := s.db.WithContext(ctx).Where("name = ? AND password = ?", username, password).First(user).Error
	return user, err
}

func (s *GaCloudServer) userLoginByEmail(ctx context.Context, email, password string) (*model.User, error) {
	password = utils.Sha256SumBytes([]byte(password))

	user := new(model.User)
	err := s.db.WithContext(ctx).Where("email = ? AND password = ?", email, password).First(user).Error
	return user, err
}
