package model

import (
	"time"
	"fmt"
)

type Meetingroom struct {
	ID           uint          `gorm:"AUTO_INCREMENT" json:"id"`
	Building     Building      `json:"building"`
	BuildingID   uint          `json:"-"`
	Floor        int           `json:"floor"`
	Room         string        `gorm:"type:varchar(16)" json:"room"`
	Weekplan     Weekplan      `json:"-"`
	WeekplanID   uint          `json:"-"`
	Reservations []Reservation `json:"reservations"`
	AvlTime      []TimeSlice   `gorm:"-" json:"avl_time"`
}

// GetAvlTime will calculate available time slices from now to the midnight of
// the day after "dayCnt" days. Before calling this func, make sure that Weekplan
// and Dayplans of the Weekplan and Timeplans of the Dayplans are set correctly.
func (m *Meetingroom) GetAvlTime(dayCnt uint) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	nowStr := now.Format("2006-01-02 15:04:05")
	end := now.AddDate(0, 0, int(dayCnt))
	endStr := end.Format("2006-01-02") + " 23:59:59"
	fmt.Println(nowStr)
	fmt.Println(endStr)
	Db.Where("meetingroom_id = ?", m.ID).
		Where("begin > ?", nowStr).Where("end < ?", endStr).
		Order("reservations.begin ASC").
		Find(&m.Reservations)

	fmt.Println(m.Reservations)
	fmt.Println(m.Weekplan.Dayplans[0])
	for i:=0; i < int(dayCnt); i++ {
		// TODO: calculate per-day AvlTime
	}
	return
}
