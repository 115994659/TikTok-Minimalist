package controller

import (
	"errors"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProxyPostFavorHandler struct {
	*gin.Context
	userId     int64
	videoId    int64
	actionType int64
}

type ProxyFavorVideoListHandler struct {
	*gin.Context
	userId int64
}

type FavorVideoListResponse struct {
	model.Response
	*video.FavorList
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteActionController(c *gin.Context) {
	p := &ProxyPostCommentHandler{Context: c}
	rawUserId, _ := p.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		errors.New("userId解析出错")
	}
	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	p.videoId = videoId
	p.actionType = actionType
	p.userId = userId

	//正式调用
	videoService := InitVideo()
	err2 := videoService.PostFavorState(p.userId, p.videoId, p.actionType)
	if err2 != nil {
		p.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	//成功返回
	p.JSON(http.StatusOK, model.Response{StatusCode: 0})
}

// FavoriteList all users have same favorite video list
func FavoriteListController(c *gin.Context) {
	p := &ProxyFavorVideoListHandler{Context: c}
	userid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	p.userId = userid
	videoService := InitVideo()
	favorVideoList, err := videoService.QueryFavorVideoList(p.userId)
	if err != nil {
		p.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	p.JSON(http.StatusOK, FavorVideoListResponse{Response: model.Response{StatusCode: 0},
		FavorList: favorVideoList,
	})
}
