package model

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	// gorm.Model
	ID       uint `gorm:"AUTO_INCREMENT"`
	Username string
	Password string
}
