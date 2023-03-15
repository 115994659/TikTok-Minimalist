package controller

import (
	"errors"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service/comment"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type List2Response struct {
	model.Response
	*comment.List
}
type CommentListResponse struct {
	model.Response
	CommentList []Comment `json:"comment_list,omitempty"`
}
type PostCommentResponse struct {
	model.Response
	*comment.MyResponse
}

type CommentActionResponse struct {
	model.Response
	Comment Comment `json:"comment,omitempty"`
}
type ProxyPostCommentHandler struct {
	*gin.Context
	videoId     int64
	userId      int64
	commentId   int64
	actionType  int64
	commentText string
}
type ProxyCommentListHandler struct {
	*gin.Context

	videoId int64
	userId  int64
}

// CommentAction no practical effect, just check if token is valid
func CommentActionController(c *gin.Context) {

	p := &ProxyPostCommentHandler{Context: c}
	rawUserId, _ := p.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		errors.New("userId解析出错")
	}
	p.userId = userId
	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		errors.New("video_id解析出错")
	}
	p.videoId = videoId
	//根据actionType进行不同赋值（1-发布评论，2-删除评论）
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	switch actionType {
	case comment.CREATE:
		p.commentText = p.Query("comment_text")
	case comment.DELETE:
		p.commentId, err = strconv.ParseInt(p.Query("comment_id"), 10, 64)
		if err != nil {
			errors.New("comment_id解析出错")
		}
	}
	p.actionType = actionType
	commentService := InitComment()
	commentRes, err := commentService.PostComment(p.userId, p.videoId, p.commentId, p.actionType, p.commentText)
	p.JSON(http.StatusOK, PostCommentResponse{
		Response:   model.Response{StatusCode: 0},
		MyResponse: commentRes,
	})
}

// CommentList all videos have same demo comment list
func CommentListController(c *gin.Context) {
	p := &ProxyCommentListHandler{Context: c}
	rawUserId, _ := p.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		errors.New("userId解析出错")
	}
	p.userId = userId

	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		errors.New("videoId解析出错")
	}
	p.videoId = videoId
	commentService := InitComment()
	commentList, err := commentService.QueryCommentList(p.userId, p.videoId)
	p.JSON(http.StatusOK, List2Response{Response: model.Response{StatusCode: 0},
		List: commentList,
	})
}

func InitComment() comment.CommentServiceImpl {
	var commentService comment.CommentServiceImpl
	return commentService
}
