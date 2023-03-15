package controller

import (
	"errors"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service/video"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type ListResponse struct {
	model.Response
	*video.List
}
type ProxyQueryVideoList struct {
	c *gin.Context
}

func NewProxyQueryVideoList(c *gin.Context) *ProxyQueryVideoList {
	return &ProxyQueryVideoList{c: c}
}

type VideoListResponse struct {
	model.Response
	VideoList []Video `json:"video_list"`
}

//上传视频后缀限制
var (
	videoSuffixMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
)

// 选择视频上传，封面采用ffmpeg工具截屏
func PublishController(c *gin.Context) {

	//获取userId  查看颁发token时的对应key 示例userId,
	//context  *gin.Context   context.Set("userId", token.Id)
	rawId, _ := c.Get("userId")
	userId, _ := rawId.(int64)
	//userId测试数据，默认为1
	//userId := int64(1)
	title := c.PostForm("title")
	form, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}

	files := form.File["data"]
	for _, file := range files {
		suffix := filepath.Ext(file.Filename)
		if _, ok := videoSuffixMap[suffix]; !ok { //判断是否为视频格式
			PublishVideoError(c, "不支持的视频格式")
			continue
		}
		name := util.GenerateFileName(userId)
		filename := name + suffix
		savePath := filepath.Join("./static", filename)
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		//封面选取
		err = util.CoverImageVideo(name, true)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		videoService := InitVideo()
		err := videoService.PostVideo(userId, filename, name+util.GetDefaultImageSuffix(), title)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		PublishVideoOk(c, file.Filename+"上传成功")
	}
}

// PublishList all users have same publish video list
func PublishListController(c *gin.Context) {
	p := NewProxyQueryVideoList(c)
	rawId, _ := c.Get("userId")
	//userid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	err := p.QueryVideoListByUserId(rawId)
	//err := p.QueryVideoListByUserId(userid)
	if err != nil {
		p.QueryVideoListError(err.Error())
	}
}

func (p *ProxyQueryVideoList) QueryVideoListByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	videoService := InitVideo()
	videoList, err := videoService.QueryVideoListByUserId(userId)
	if err != nil {
		return err
	}

	p.QueryVideoListOk(videoList)
	return nil
}

func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, model.Response{StatusCode: 1,
		StatusMsg: msg})
}

func PublishVideoOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: msg})
}

func (p *ProxyQueryVideoList) QueryVideoListError(msg string) {
	p.c.JSON(http.StatusOK, ListResponse{Response: model.Response{
		StatusCode: 1,
		StatusMsg:  msg,
	}})
}

func (p *ProxyQueryVideoList) QueryVideoListOk(videoList *video.List) {
	p.c.JSON(http.StatusOK, ListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		List: videoList,
	})
}
