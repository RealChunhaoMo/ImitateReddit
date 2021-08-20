package mysql

import (
	"WebApp/modules"
	"crypto/md5"
	"encoding/hex"
)

const secret = "www.bilibili.com"

// CheckUserExist 根据用户名检查用户是否存在
func CheckUserExist(Username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, Username); err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertUser 往数据库中插入一个新用户
func InsertUser(user *modules.User) (err error) {
	//对密码加密
	user.Password = EncryptPassword(user.Password)
	//执行SQL语句
	sqlStr := `insert into  user(user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, user.Password)
	return
}

// EncryptPassword 用md5对密码加密
func EncryptPassword(OriginPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(OriginPassword)))
}

// PasswordIsRight 检查用户输入的密码是否正确
func PasswordIsRight(UserName, Password string) (bool, error) {
	sqlStr := `select password from user where username = ?`

	//因为当时注册的时候，用户的密码加密了那么此时密码验证也要用加密形式验证
	Password = EncryptPassword(Password)
	var OriginPassword string
	if err := db.Get(&OriginPassword, sqlStr, UserName); err != nil {
		return false, err
	}
	return OriginPassword == Password, nil
}
