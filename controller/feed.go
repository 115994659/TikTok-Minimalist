package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service/video"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	model.Response
	*video.FeedVideoList
}

type ProxyFeedVideoList struct {
	*gin.Context
}

func NewProxyFeedVideoList(c *gin.Context) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{Context: c}
}

func FeedVideoListController(c *gin.Context) {
	p := NewProxyFeedVideoList(c)
	//token --> _
	_, ok := c.GetQuery("token")

	//无登录状态
	if !ok {
		err := p.NoTokenFeed()
		if err != nil {
			p.FeedError(err.Error())
		}
		return
	}
}

func (p *ProxyFeedVideoList) NoTokenFeed() error {
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, intTime*1e6) //时间戳是以ms为单位的
	}
	videoService := InitVideo()
	videoList, err := videoService.Feed(latestTime, 0)
	if err != nil {
		log.Printf("service.video.videoService.Feed(lastTime, userId)出现异常：%v", err)
	}
	if err != nil {
		return err
	}
	p.FeedData(videoList)
	return nil
}

func (p *ProxyFeedVideoList) FeedError(msg string) {
	p.JSON(http.StatusOK, FeedResponse{Response: model.Response{
		StatusCode: 1,
		StatusMsg:  msg,
	}})
}

// Feed same demo video list for every request
func (p *ProxyFeedVideoList) FeedData(videoList *video.FeedVideoList) {
	p.JSON(http.StatusOK, FeedResponse{
		Response:      model.Response{StatusCode: 0},
		FeedVideoList: videoList,
	})
}

func InitVideo() video.VideoServiceImpl {
	var videoService video.VideoServiceImpl
	return videoService
}
