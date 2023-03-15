package user

import (
	"errors"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/model"
)

type UserServiceImpl struct {
}

func (userService UserServiceImpl) PostUserLogin(username, password string) (*model.LoginResponse, error) {
	q := &PostUserLoginFlow{username: username, password: password}
	if q.username == "" {
		return nil, errors.New("用户名为空")
	}
	if q.password == "" {
		return nil, errors.New("密码为空")
	}
	//准备好userInfo,默认name为username
	userLogin := model.UserLogin{Username: q.username, Password: q.password}
	userinfo := model.UserInfo{User: &userLogin, Name: q.username}

	//判断用户名是否已经存在
	userLoginDAO := model.NewUserLoginDao()
	if userLoginDAO.IsUserExistByUsername(q.username) {
		return nil, errors.New("用户名已存在")
	}

	//更新操作，由于userLogin属于userInfo，故更新userInfo即可，且由于传入的是指针，所以插入的数据内容也是清楚的
	userInfoDAO := model.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&userinfo)
	if err != nil {
		return nil, err
	}

	//颁发token
	token, err := middleware.ReleaseToken(userLogin)
	if err != nil {
		return nil, err
	}
	q.token = token
	q.userid = userinfo.Id
	q.data = &model.LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return q.data, nil
}

func (userService UserServiceImpl) QueryUserLogin(username, password string) (*model.LoginResponse, error) {
	q := &QueryUserLoginFlow{username: username, password: password}
	if q.username == "" {
		return nil, errors.New("用户名为空")
	}
	if q.password == "" {
		return nil, errors.New("密码为空")
	}
	userLoginDAO := model.NewUserLoginDao()
	var login model.UserLogin
	//准备好userid
	err := userLoginDAO.QueryUserLogin(q.username, q.password, &login)
	if err != nil {
		return nil, err
	}
	q.userid = login.UserInfoId

	//准备颁发token
	token, err := middleware.ReleaseToken(login)
	if err != nil {
		return nil, err
	}
	q.token = token
	q.data = &model.LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return q.data, nil
}
