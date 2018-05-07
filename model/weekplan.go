package model

type Weekplan struct {
	ID           uint          `gorm:"AUTO_INCREMENT" json:"id"`
	Name         string        `gorm:"type:varchar(64)" json:"name"`
	Meetingrooms []Meetingroom `json:"meetingrooms"`
	Dayplans     []Dayplan     `gorm:"many2many:weekplan_dayplan" json:"dayplans"`
}
