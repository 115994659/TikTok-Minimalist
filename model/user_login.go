package model

import (
	"errors"
	"sync"
)

type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserLoginDAO struct {
}

var (
	userLoginDao  *UserLoginDAO
	userLoginOnce sync.Once
)

func NewUserLoginDao() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDao = new(UserLoginDAO)
	})
	return userLoginDao
}
func (u *UserLoginDAO) IsUserExistByUsername(username string) bool {
	var userLogin UserLogin
	DB.Where("username=?", username).First(&userLogin)
	if userLogin.Id == 0 {
		return false
	}
	return true
}

func (u *UserLoginDAO) QueryUserLogin(username string, password string, m *UserLogin) error {
	if m == nil {
		return errors.New("结构体指针为空")
	}
	DB.Where("username=? and password=?", username, password).First(m)
	if m.Id == 0 {
		return errors.New("用户不存在，账号或密码出错")
	}
	return nil
}
