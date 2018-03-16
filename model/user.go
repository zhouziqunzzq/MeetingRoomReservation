package model

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type User struct {
	// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	// gorm.Model
	ID       uint   `gorm:"AUTO_INCREMENT"`
	Username string `gorm:"unique"`
	Password string
}

func GetUserByUsername(db *gorm.DB, username string)(user User, err error) {
	if err = db.Where(&User{Username: username}).First(&user).Error; err != nil {
		err = errors.Wrap(err, "GetUser")
	}
	return
}
