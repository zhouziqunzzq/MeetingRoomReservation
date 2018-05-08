package model

import "time"

type Weekplan struct {
	ID           uint          `gorm:"AUTO_INCREMENT" json:"id"`
	Name         string        `gorm:"type:varchar(64)" json:"name"`
	Meetingrooms []Meetingroom `json:"meetingrooms"`
	Dayplans     []Dayplan     `gorm:"many2many:weekplan_dayplan" json:"dayplans"`
}

func (weekplan *Weekplan) ConvertToMap() (map[time.Weekday][]Timeplan) {
	rst := make(map[time.Weekday][]Timeplan)
	for i := 0; i < len(weekplan.Dayplans); i++ {
		rst[time.Weekday(weekplan.Dayplans[i].Weekday%7)] = weekplan.Dayplans[i].Timeplans
	}
	return rst
}
