package model

import "time"

type Reservation struct {
	ID            uint        `gorm:"AUTO_INCREMENT" json:"id"`
	User          User        `json:"user"`
	UserID        uint        `json:"-"`
	Meetingroom   Meetingroom `json:"meetingroom"`
	MeetingroomID uint        `json:"-"`
	Begin         string      `gorm:"type:varchar(20)" json:"begin"` // YYYY-MM-DD HH:MM:SS
	End           string      `gorm:"type:varchar(20)" json:"end"`   // YYYY-MM-DD HH:MM:SS
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     *time.Time  `json:"-"`
}
