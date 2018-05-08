package model

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

const (
	ADMIN   = 0
	TEACHER = 1
)

type User struct {
	// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	// gorm.Model
	ID       uint   `gorm:"AUTO_INCREMENT" json:"id"`
	Username string `gorm:"unique;type:varchar(64)" json:"username"`
	Password string `gorm:"type:varchar(64)" json:"-"` // avoid password in json
	Name     string `gorm:"type:varchar(64)" json:"name"`
	Tel      string `gorm:"type:varchar(64)" json:"tel"`
	Type     int    `gorm:"type:tinyint(3)" json:"type"` // 0 for admin, 1 for teacher
}

func GetUserByUsername(username string) (user User, err error) {
	if err = Db.Where(&User{Username: username}).First(&user).Error; err != nil {
		err = errors.Wrap(err, "GetUser")
	}
	return
}

func GetUserInfoByID(uid uint) (user User, err error) {
	if err = Db.Select("id, username, name, tel, type").
		Where(&User{ID: uid}).First(&user).Error; err != nil {
		err = errors.Wrap(err, "GetUserInfo")
	}
	return
}

func CheckAdmin(uid uint) bool {
	user, err := GetUserInfoByID(uid)
	if err != nil || user.Type != ADMIN {
		return false
	}
	return true
}
