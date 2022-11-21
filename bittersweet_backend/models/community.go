package models

import "time"

// Community 社区分类列表结构体
type Community struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

// CommunityDetail 社区详情信息结构体
type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"` // 注：使用time.Time  sql init时就要加上 parseTime=true
}
