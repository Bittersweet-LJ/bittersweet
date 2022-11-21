package logic

import (
	"bittersweet/dao/mysql"
	"bittersweet/dao/redis"
	"bittersweet/models"
	"bittersweet/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建新帖子
func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	//3.返回
	return

}

// GetPostByID 根据帖子id查询帖子详情信息
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	//查询并组合我们接口想要的数据

	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	//根据作者id查询用户信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	//根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}

	//组合数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts)) // 根据查到的帖子数 初始化data

	for _, post := range posts {
		//根据作者id查询用户信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("for range posts mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("for range posts mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 获取帖子列表(可排序升级版)
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1.去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("redis.GetPostIDsInOrder(p)", zap.Any("ids", ids))
	// 2.根据id去mysql数据库查询帖子详细信息
	// 要按照给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("mysql.GetPostListByIDs(ids)", zap.Any("posts", posts))
	// 提前查询出每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	// 查询出帖子作者和社区信息并填充
	for idx, post := range posts {
		//根据作者id查询用户信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("for range posts mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("for range posts mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetCommunityPostList  根据所属社区获取帖子列表
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1.去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("redis.GetPostIDsInOrder(p)", zap.Any("ids", ids))
	// 2.根据id去mysql数据库查询帖子详细信息
	// 要按照给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("mysql.GetPostListByIDs(ids)", zap.Any("posts", posts))
	// 提前查询出每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	// 查询出帖子作者和社区信息并填充
	for idx, post := range posts {
		//根据作者id查询用户信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("for range posts mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("for range posts mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew (整合完善后的 获取帖子列表函数) GetPostList2 + GetCommunityPostList
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//根据请求参数不同，执行不同逻辑
	if p.CommunityID == 0 {
		//前端没有传来community_id , 故查询所有帖子
		data, err = GetPostList2(p)
	} else {
		//不为空就按指定社区查询帖子
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
