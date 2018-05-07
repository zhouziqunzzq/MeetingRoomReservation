package model

type Timeplan struct {
	ID       uint    `gorm:"AUTO_INCREMENT" json:"id"`
	Begin    string  `gorm:"type:varchar(10)" json:"begin"` // HH:MM:SS
	End      string  `gorm:"type:varchar(10)" json:"end"`   // HH:MM:SS
	Dayplans Dayplan `gorm:"many2many:dayplan_timeplan" json:"-"`
}
