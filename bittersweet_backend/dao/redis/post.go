package redis

import (
	"bittersweet/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// GetPostIDsInOrder 根据分页结构体中的参数去redis查询id列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的索引起始点
	//3. ZREVRANGE 按分数从高到低的顺序查询指定数量的元素,并返回
	//封装到 getIDsFromKey中，方便调用
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//for _, id := range ids {
	//	data = make([]int64, 0, len(ids))
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	//查找每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	//使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 根据社区去redis查询id列表
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//1.根据用户请求中携带的order参数确定要查询的redis key
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//使用zinterstore 用 社区的帖子set 与 帖子分数的zset 生成一个新的zset
	//对新的zset 按之前的逻辑取数据

	//社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数 !!!!!!!!!!!!!!!!!!!!!!!!!!!
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		//不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) //zinterstore 计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在的话就直接根据key查询ids
	//2.确定查询的索引起始点
	//3. ZREVRANGE 按分数从高到低的顺序查询指定数量的元素
	//封装到 getIDsFromKey中，方便调用
	return getIDsFromKey(key, p.Page, p.Size)
}

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	//确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	//ZREVRANGE 按分数从高到低的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}
