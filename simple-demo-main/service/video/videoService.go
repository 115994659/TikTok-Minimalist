package video

import (
	"github.com/RaymondCode/simple-demo/model"
	"time"
)

type QueryFeedVideoListFlow struct {
	userId     int64
	latestTime time.Time
	videos     []*model.Video
	nextTime   int64
	feedVideo  *FeedVideoList
}

type List struct {
	Videos []*model.Video `json:"video_list,omitempty"`
}
type FavorList struct {
	Videos []*model.Video `json:"video_list"`
}
type VideoService interface {
	//视频流接口
	Feed(lastTime time.Time, userId int64) (*FeedVideoList, error)
	//上传视频（支持多个视频同时上传）接口
	PostVideo(userId int64, videoname string, coverName string, title string) error
	//查询视频，根据用户id查询其下所有视频
	QueryVideoListByUserId(userId int64) (*List, error)
	//点赞视频
	PostFavorState(userId, videoId, actionType int64) error
	//查询用户点赞视频
	QueryFavorVideoList(userId int64) (*FavorList, error)
}
