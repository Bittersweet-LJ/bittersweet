package logic

import (
	"bittersweet/dao/mysql"
	"bittersweet/models"
	"bittersweet/pkg/jwt"
	"bittersweet/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询错误
		return err
	}
	// 生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 保存到数据库

	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，就可以通过db.Get拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
