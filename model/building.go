package model

type Building struct {
	ID           uint          `gorm:"AUTO_INCREMENT" json:"id"`
	Name         string        `gorm:"type:varchar(64)" json:"name"`
	MaxFloor     uint          `gorm:"type:tinyint(3)" json:"max_floor"`
	Meetingrooms []Meetingroom `json:"meetingrooms"`
}
