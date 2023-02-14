package util

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"log"
	"path/filepath"
)

//通过用户Id以及用户的视频投稿数量生成文件名
func GenerateFileName(userId int64) string {
	var count int64
	err := model.NewVideoDAO().QueryVideoCountByUserId(userId, &count)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d", userId, count)
}

//通过截屏工具生成照片
func CoverImageVideo(name string, isDebug bool) error {
	v := NewVideoCoverImage()
	if isDebug {
		v.Debug()
	}
	v.InputPath = filepath.Join(config.Info.StaticSourcePath, name+defaultVideoSuffix)
	v.OutputPath = filepath.Join(config.Info.StaticSourcePath, name+defaultImageSuffix)
	v.FrameCount = 1
	queryString, err := v.GetQueryString()
	if err != nil {
		return err
	}
	return v.ExecCommand(queryString)
}

func GetFileUrl(fileName string) string {
	base := fmt.Sprintf("http://%s:%d/static/%s", config.Info.IP, config.Info.Port, fileName)
	return base
}
