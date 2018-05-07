package model

type Dayplan struct {
	ID        uint       `gorm:"AUTO_INCREMENT" json:"id"`
	Weekday   uint       `gorm:"type:tinyint(3)" json:"weekday"`
	Weekplans []Weekplan `gorm:"many2many:weekplan_dayplan" json:"weekplans"`
	Timeplans []Timeplan `gorm:"many2many:dayplan_timeplan" json:"timeplans"`
}
