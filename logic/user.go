package logic

import (
	"WebApp/dao/mysql"
	"WebApp/modules"
	"WebApp/pkg/jwt"
	"WebApp/pkg/snowflake"
	"errors"
)

var (
	UerExist      = errors.New("用户已存在")
	UserNameError = errors.New("用户名错误")
	PasswordError = errors.New("密码错误")
)

func SignUp(p *modules.ParamSignUp) (err error) {
	// 用户查询是否存在,避免重名
	var exist bool
	exist, err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return errors.New("数据库查询用户是否存在时出错了")
	}
	if exist {
		return UerExist
	}
	// 生成UID
	userID := snowflake.GenID()
	//构造一个user实例
	user := modules.User{
		UserID:   userID,
		UserName: p.Username,
		Password: p.Password,
	}
	// 保存到数据库
	err = mysql.InsertUser(&user)
	return
}

func SignIn(p *modules.ParamSignIn) (user *modules.User, err error) {
	var right bool
	right, err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return nil, err
	}
	if !right {
		return nil, UserNameError
	}

	user = &modules.User{
		Password: p.Password,
		UserName: p.Username,
	}
	right, err = mysql.PasswordIsRight(user)
	if err != nil {
		return nil, err
	}
	if !right {
		return nil, PasswordError
	}
	//生成jwt
	token, err := jwt.GenToken(user.UserID)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
