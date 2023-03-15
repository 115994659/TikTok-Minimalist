package model

import (
	"errors"
	"log"
	"sync"
)

type UserInfo struct {
	Id            int64       `json:"id" gorm:"id,omitempty"`
	Name          string      `json:"name" gorm:"name,omitempty"`
	FollowCount   int64       `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64       `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool        `json:"is_follow" gorm:"is_follow,omitempty"`
	User          *UserLogin  `json:"-"`                                     //用户与账号密码之间的一对一
	Videos        []*Video    `json:"-"`                                     //用户与投稿视频的一对多
	Follows       []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos   []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments      []*Comment  `json:"-"`                                     //用户与评论的一对多
}
type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once
)

func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}

func (u *UserInfoDAO) QueryUserInfoById(userId int64, userInfo *UserInfo) error {
	if userInfo == nil {
		return errors.New("空指针异常")
	}
	//DB.Where("id=?",userId).First(userinfo)
	DB.Where("id=?", userId).Select([]string{"id", "name", "follow_count", "follower_count", "is_follow"}).First(userInfo)
	//id为零值，说明sql执行失败
	if userInfo.Id == 0 {
		return errors.New("该用户不存在")
	}
	return nil
}

func (u *UserInfoDAO) AddUserInfo(userinfo *UserInfo) error {
	if userinfo == nil {
		return errors.New("空指针错误")
	}
	return DB.Create(userinfo).Error
}

func (u *UserInfoDAO) IsUserExistById(id int64) bool {
	var userinfo UserInfo
	if err := DB.Where("id=?", id).Select("id").First(&userinfo).Error; err != nil {
		log.Println(err)
	}
	if userinfo.Id == 0 {
		return false
	}
	return true
}

func (v *VideoDAO) QueryFavorVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryFavorVideoListByUserId videoList 空指针")
	}
	//多表查询，左连接得到结果，再映射到数据
	if err := DB.Raw("SELECT v.* FROM user_favor_videos u , videos v WHERE u.user_info_id = ? AND u.video_id = v.id", userId).Scan(videoList).Error; err != nil {
		return err
	}
	//如果id为0，则说明没有查到数据
	if len(*videoList) == 0 || (*videoList)[0].Id == 0 {
		return errors.New("点赞列表为空")
	}
	return nil
}
