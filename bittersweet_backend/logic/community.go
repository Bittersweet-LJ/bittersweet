package logic

import (
	"bittersweet/dao/mysql"
	"bittersweet/models"
)

func GetCommunityList() ([]*models.Community, error) {
	//从数据库查找到所有的community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
