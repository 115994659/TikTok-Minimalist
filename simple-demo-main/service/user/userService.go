package user

import (
	"github.com/RaymondCode/simple-demo/model"
)

type PostUserLoginFlow struct {
	username string
	password string

	data   *model.LoginResponse
	userid int64
	token  string
}
type QueryUserLoginFlow struct {
	username string
	password string

	data   *model.LoginResponse
	userid int64
	token  string
}
type VideoService interface {
	PostUserLogin(username, password string) (*model.LoginResponse, error)
	QueryUserLogin(username, password string) (*model.LoginResponse, error)
}
