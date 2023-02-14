package video

import (
	"errors"
	"github.com/RaymondCode/simple-demo/cache"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/util"
	"time"
)

//VideoNum 每次返回的视频流数量
const (
	VideoNum = 30
)

//点赞视频常量
const (
	PLUS  = 1
	MINUS = 2
)

type VideoServiceImpl struct {
}

//参数结构体------FeedFeed(lastTime time.Time, userId int64)参数结构体
type FeedVideoList struct {
	Videos   []*model.Video `json:"video_list,omitempty"`
	NextTime int64          `json:"next_time,omitempty"`
}

//视频流接口实现，限制(limit)条返回
func (videoService VideoServiceImpl) Feed(lastTime time.Time, userId int64) (*FeedVideoList, error) {
	q := &QueryFeedVideoListFlow{userId: userId, latestTime: lastTime}
	//上层通过把userId置零，表示userId不存在或不需要
	if q.userId > 0 {
		//说明用户是登录状态
	}
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
	err := model.NewVideoDAO().QueryVideoListByLimitAndTime(VideoNum, q.latestTime, &q.videos)
	if err != nil {
		return nil, err
	}
	//用户信息注入
	size := len(q.videos)
	userDao := model.NewUserInfoDAO()
	for i := 0; i < size; i++ {
		var userInfo model.UserInfo
		err := userDao.QueryUserInfoById((q.videos)[i].UserInfoId, &userInfo)
		if err != nil {
			continue
		}
		(q.videos)[i].Author = userInfo
	}
	latestTime := q.videos[size-1].PublishTime
	//为最新投稿时间戳赋值
	if &latestTime != nil {
		q.nextTime = (latestTime).UnixNano() / 1e6
	}
	if &latestTime == nil {
		q.nextTime = time.Now().Unix() / 1e6
	}
	q.feedVideo = &FeedVideoList{
		Videos:   q.videos,
		NextTime: q.nextTime,
	}
	return q.feedVideo, nil
}

//参数结构体------PostVideo(userId int64, videoename string, coverName string, title string)
type PostVideoFlow struct {
	videoName string
	coverName string
	title     string
	userId    int64
	video     *model.Video
}

func (videoService VideoServiceImpl) PostVideo(userId int64, videoName string, coverName string, title string) error {
	f := &PostVideoFlow{
		videoName: videoName,
		coverName: coverName,
		userId:    userId,
		title:     title,
	}
	f.videoName = util.GetFileUrl(f.videoName)
	f.coverName = util.GetFileUrl(f.coverName)
	video := &model.Video{
		UserInfoId: f.userId,
		PlayUrl:    f.videoName,
		CoverUrl:   f.coverName,
		Title:      f.title,
	}
	return model.NewVideoDAO().AddVideo(video)
}

type QueryVideoListByUserIdFlow struct {
	userId int64
	videos []*model.Video

	videoList *List
}

func (videoService VideoServiceImpl) QueryVideoListByUserId(userId int64) (*List, error) {
	q := &QueryVideoListByUserIdFlow{userId: userId}
	if !model.NewUserInfoDAO().IsUserExistById(q.userId) {
		return nil, errors.New("用户不存在")
	}
	err := model.NewVideoDAO().QueryVideoListByUserId(q.userId, &q.videos)
	if err != nil {
		return nil, err
	}
	//作者信息查询
	var userInfo model.UserInfo
	err = model.NewUserInfoDAO().QueryUserInfoById(q.userId, &userInfo)
	p := cache.NewProxyIndexMap()
	if err != nil {
		return nil, err
	}
	//填充信息(Author和IsFavorite字段
	for i := range q.videos {
		q.videos[i].Author = userInfo
		q.videos[i].IsFavorite = p.GetVideoFavorState(q.userId, q.videos[i].Id)
	}

	q.videoList = &List{Videos: q.videos}
	return q.videoList, nil
}

type PostFavorStateFlow struct {
	userId     int64
	videoId    int64
	actionType int64
}

func (videoService VideoServiceImpl) PostFavorState(userId, videoId, action int64) error {
	var err error
	p := &PostFavorStateFlow{
		userId:     userId,
		videoId:    videoId,
		actionType: action,
	}
	if !model.NewUserInfoDAO().IsUserExistById(p.userId) {
		return errors.New("用户不存在")
	}
	switch p.actionType {
	case PLUS:
		err = p.PlusOperation()
	case MINUS:
		err = p.MinusOperation()
	}
	return err
}

type QueryFavorVideoListFlow struct {
	userId    int64
	videos    []*model.Video
	videoList *FavorList
}

func (videoService VideoServiceImpl) QueryFavorVideoList(id int64) (*FavorList, error) {
	q := &QueryFavorVideoListFlow{userId: id}
	if !model.NewUserInfoDAO().IsUserExistById(q.userId) {
		return nil, errors.New("用户未登陆！")
	}
	err := model.NewVideoDAO().QueryFavorVideoListByUserId(q.userId, &q.videos)

	for i := range q.videos {
		//作者信息查询
		var userInfo model.UserInfo
		err = model.NewUserInfoDAO().QueryUserInfoById(q.videos[i].UserInfoId, &userInfo)
		if err == nil { //若查询未出错则更新，否则不更新作者信息
			q.videos[i].Author = userInfo
		}
		q.videos[i].IsFavorite = true
	}
	q.videoList = &FavorList{Videos: q.videos}
	return q.videoList, nil
}

// PlusOperation 点赞操作
func (p *PostFavorStateFlow) PlusOperation() error {
	//视频点赞数目+1
	err := model.NewVideoDAO().PlusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
	if err != nil {
		return errors.New("不要重复点赞")
	}
	//对应的用户是否点赞的映射状态更新
	cache.NewProxyIndexMap().UpdateVideoFavorState(p.userId, p.videoId, true)
	return nil
}

// MinusOperation 取消点赞
func (p *PostFavorStateFlow) MinusOperation() error {
	//视频点赞数目-1
	err := model.NewVideoDAO().MinusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
	if err != nil {
		return errors.New("点赞数目已经为0")
	}
	//对应的用户是否点赞的映射状态更新
	cache.NewProxyIndexMap().UpdateVideoFavorState(p.userId, p.videoId, false)
	return nil
}
