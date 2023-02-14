package cache

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

const (
	favor = "favor"
)

func init() {
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Info.RDB.IP, config.Info.RDB.Port),
			Password: "", //没有设置密码
			DB:       config.Info.RDB.Database,
		})
}

func NewProxyIndexMap() *ProxyIndexMap {
	return &proxyIndexOperation
}

var (
	proxyIndexOperation ProxyIndexMap
)

type ProxyIndexMap struct {
}

// GetVideoFavorState 得到点赞状态
func (i *ProxyIndexMap) GetVideoFavorState(userId int64, videoId int64) bool {
	key := fmt.Sprintf("%s:%d", favor, userId)
	ret := rdb.SIsMember(ctx, key, videoId)
	return ret.Val()
}

// UpdateVideoFavorState 更新点赞状态，state:true为点赞，false为取消点赞
func (i *ProxyIndexMap) UpdateVideoFavorState(userId int64, videoId int64, state bool) {
	key := fmt.Sprintf("%s:%d", favor, userId)
	if state {
		rdb.SAdd(ctx, key, videoId)
		return
	}
	rdb.SRem(ctx, key, videoId)
}
