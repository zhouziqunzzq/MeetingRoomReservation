package model

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type User struct {
	// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	// gorm.Model
	ID       uint   `gorm:"AUTO_INCREMENT" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Tel      string `json:"tel"`
	Type     int    `json:"type"` // 0 for admin, 1 for teacher
}

func GetUserByUsername(db *gorm.DB, username string) (user User, err error) {
	if err = db.Where(&User{Username: username}).First(&user).Error; err != nil {
		err = errors.Wrap(err, "GetUser")
	}
	return
}

func GetUserInfoByID(db *gorm.DB, uid uint) (user User, err error) {
	if err = db.Select("id, username, name, tel, type").Where(&User{ID: uid}).First(&user).Error; err != nil {
		err = errors.Wrap(err, "GetUserInfo")
	}
	return
}
