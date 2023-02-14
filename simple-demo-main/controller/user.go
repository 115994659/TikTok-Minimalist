package controller

import (
	"errors"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	model.Response
	*model.LoginResponse
}

type UserRegisterResponse struct {
	model.Response
	*model.LoginResponse
}

type UserResponse struct {
	model.Response
	User *model.UserInfo `json:"user"`
}

func RegisterController(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userService := InitUser()
	registerResponse, err := userService.PostUserLogin(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response:      model.Response{StatusCode: 0},
		LoginResponse: registerResponse,
	})
}

func LoginController(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userService := InitUser()
	userLoginResponse, err := userService.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:      model.Response{StatusCode: 0},
		LoginResponse: userLoginResponse,
	})
}

//业务逻辑较为简单，不走service层，在controller层解决需求
func UserInfoController(c *gin.Context) {
	p := NewProxyUserInfo(c)
	userid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//userId, ok := c.Get("userId")
	//if !ok {
	//	p.UserInfoError("解析userId出错")
	//	return
	//}
	//err := p.QueryUserInfoByUserId(userId)
	err := p.QueryUserInfoByUserId(userid)
	if err != nil {
		p.UserInfoError(err.Error())
	}
}

func InitUser() user.UserServiceImpl {
	var userService user.UserServiceImpl
	return userService
}

type ProxyUserInfo struct {
	c *gin.Context
}

func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{c: c}
}
func (p *ProxyUserInfo) QueryUserInfoByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("解析userId失败")
	}
	//由于得到userinfo不需要组装model层的数据，所以直接调用model层的接口
	userinfoDAO := model.NewUserInfoDAO()

	var userInfo model.UserInfo
	err := userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, UserResponse{
		Response: model.Response{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyUserInfo) UserInfoOk(user *model.UserInfo) {
	p.c.JSON(http.StatusOK, UserResponse{
		Response: model.Response{StatusCode: 0},
		User:     user,
	})
}
