package controller

import (
	"bittersweet/logic"
	"bittersweet/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从context取出当前用户的ID
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.业务处理——创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数(从URL中帖子的id)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据ID获取帖子数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表分页
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := getPageInfo(c)
	//获取帖子列表
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	//GET请求参数(query string)：/api/v1/post2?page=1&size=10&order=time
	//1.获取分页参数
	//c.ShouldBind()  根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.业务处理
	//获取帖子列表
	data, err := logic.GetPostListNew(p) // 整合完善后的获取帖子列表函数 根据前端传来的参数动态获取
	if err != nil {
		zap.L().Error("logic.GetPostListNew(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

//// GetCommunityPostListHandler 根据所属社区获取帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	//2.业务处理
//	//获取帖子列表
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("mysql.GetCommunityPostList() failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	//3.返回响应
//	ResponseSuccess(c, data)
//}
