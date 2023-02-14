package comment

import "github.com/RaymondCode/simple-demo/model"

const (
	CREATE = 1
	DELETE = 2
)

type MyResponse struct {
	MyComment *model.Comment `json:"comment"`
}

type List struct {
	Comments []*model.Comment `json:"comment_list"`
}
type CommentService interface {
	PostComment(userId int64, videoId int64, commentId int64, actionType int64, commentText string) (*MyResponse, error)
	QueryCommentList(userId, videoId int64) (*List, error)
}
