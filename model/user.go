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
}

func GetUserByUsername(db *gorm.DB, username string) (user User, err error) {
	if err = db.Where(&User{Username: username}).First(&user).Error; err != nil {
		err = errors.Wrap(err, "GetUser")
	}
	return
}
