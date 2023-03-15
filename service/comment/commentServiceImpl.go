package comment

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
)

type CommentServiceImpl struct {
}

type PostCommentFlow struct {
	userId      int64
	videoId     int64
	commentId   int64
	actionType  int64
	commentText string
	comment     *model.Comment
	*MyResponse
}

type QueryCommentListFlow struct {
	userId  int64
	videoId int64

	comments    []*model.Comment
	commentList *List
}

func (commentService CommentServiceImpl) PostComment(userId int64, videoId int64, commentId int64, actionType int64, commentText string) (*MyResponse, error) {
	p := &PostCommentFlow{userId: userId, videoId: videoId, commentId: commentId, actionType: actionType, commentText: commentText}
	if !model.NewUserInfoDAO().IsUserExistById(p.userId) {
		return nil, fmt.Errorf("用户%d不存在", p.userId)
	}
	if !model.NewVideoDAO().IsVideoExistById(p.videoId) {
		return nil, fmt.Errorf("视频%d不存在", p.videoId)
	}
	if p.actionType != CREATE && p.actionType != DELETE {
		return nil, errors.New("actionType参数错误")
	}
	var err error
	switch p.actionType {
	case CREATE:
		p.comment, err = p.CreateComment()
	case DELETE:
		p.comment, err = p.DeleteComment()
	}
	userInfo := model.UserInfo{}
	_ = model.NewUserInfoDAO().QueryUserInfoById(p.comment.UserInfoId, &userInfo)
	p.comment.User = userInfo
	_ = FillCommentFields(p.comment)

	p.MyResponse = &MyResponse{MyComment: p.comment}
	return p.MyResponse, err
}

func (commentService CommentServiceImpl) QueryCommentList(userId, videoId int64) (*List, error) {
	q := &QueryCommentListFlow{userId: userId, videoId: videoId}
	if !model.NewUserInfoDAO().IsUserExistById(q.userId) {
		return nil, fmt.Errorf("用户%d请先登陆！", q.userId)
	}
	if !model.NewVideoDAO().IsVideoExistById(q.videoId) {
		return nil, fmt.Errorf("视频%d不存在或已经被删除", q.videoId)
	}
	err := model.NewCommentDAO().QueryCommentListByVideoId(q.videoId, &q.comments)
	if err != nil {
		return nil, err
	}
	//根据前端的要求填充正确的时间格式
	err = FillCommentListFields(&q.comments)
	if err != nil {
		return nil, errors.New("评论区为空")
	}
	q.commentList = &List{Comments: q.comments}
	return q.commentList, nil
}

// CreateComment 增加评论
func (p *PostCommentFlow) CreateComment() (*model.Comment, error) {
	comment := model.Comment{UserInfoId: p.userId, VideoId: p.videoId, Content: p.commentText}
	err := model.NewCommentDAO().AddCommentAndUpdateCount(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment 删除评论
func (p *PostCommentFlow) DeleteComment() (*model.Comment, error) {
	//获取comment
	var comment model.Comment
	err := model.NewCommentDAO().QueryCommentById(p.commentId, &comment)
	if err != nil {
		return nil, err
	}
	//删除comment
	err = model.NewCommentDAO().DeleteCommentAndUpdateCountById(p.commentId, p.videoId)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func FillCommentFields(comment *model.Comment) error {
	if comment == nil {
		return errors.New("FillCommentFields comments为空")
	}
	comment.CreateDate = comment.CreatedAt.Format("1-2") //转为前端要求的日期格式
	return nil
}

func FillCommentListFields(comments *[]*model.Comment) error {
	size := len(*comments)
	if comments == nil || size == 0 {
		return errors.New("util.FillCommentListFields comments为空")
	}
	dao := model.NewUserInfoDAO()
	for _, v := range *comments {
		_ = dao.QueryUserInfoById(v.UserInfoId, &v.User) //填充这条评论的作者信息
		v.CreateDate = v.CreatedAt.Format("1-2")         //转为前端要求的日期格式
	}
	return nil
}
