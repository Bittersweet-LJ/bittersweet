package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天
/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票    --> 更新分数和投票记录  0 - 1 = -1   +432
	2. 之前投反对票，现在改投赞成票    --> 更新分数和投票记录  -1 - 1 = -2  +432*2
direction=0时，有两种情况：
	1. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录 -1 - 0 = -1    +432
	2. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录 1 - 0 = 1      -432
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票    --> 更新分数和投票记录 0 - (-1) = 1   -432
	2. 之前投赞成票，现在改投反对票    --> 更新分数和投票记录 1 - (-1) = 2   -432*2

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = float64(7 * 24 * 3600)
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	// pipeline 使用redis事务机制
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//把帖子id添加到社区的set中
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2 和 3 需要放到一个pipeline事务中操作
	//2.更新帖子的分数
	// 先查当前用户 给此帖子之前的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	// 如果此次投票的值与之前保存的值一致，直接返回，提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var direction float64
	if value > ov {
		//当前值大，是正向投票 加分
		direction = 1
	} else {
		//当前值小，是反向投票 减分
		direction = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), direction*diff*scorePerVote, postID)
	//3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, // 用户最新投的 赞成票或反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
