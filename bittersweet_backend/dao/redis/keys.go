package redis

//redis key
//使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix          = "bittersweet:"
	KeyPostTimeZSet    = "post:time"   // ZSet:帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  // ZSet:帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" // ZSet:记录用户及投票类型:参数是 post_id

	KeyCommunitySetPF = "community:" // set:保存每个社区下的帖子id
)

// getRedisKey 给redis key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
