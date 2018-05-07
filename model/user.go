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
	Username string `gorm:"unique;type:varchar(64)" json:"username"`
	Password string `gorm:"type:varchar(64)" json:"password"`
	Name     string `gorm:"type:varchar(64)" json:"name"`
	Tel      string `gorm:"type:varchar(64)" json:"tel"`
	Type     int    `gorm:"type:tinyint(3)" json:"type"` // 0 for admin, 1 for teacher
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
