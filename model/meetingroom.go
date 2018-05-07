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
	for i := 0; i < int(dayCnt); i++ {
		// TODO: calculate per-day AvlTime
	}
	return
}

// timeplans: a slice of Timeplans in a day
// reservations: a slice of Reservations in a day
// begin: HH:MM:SS
// end: HH:MM:SS
// date: YYYY-MM-DD
func GetAvlTimeOneDay(timeplans []Timeplan, reservations []Reservation,
	begin string, end string, date string) []TimeSlice {
	// TODO: calculate per-day AvlTime
	timeslices := make([]TimeSlice, 0)
	now := 0;
	for i := 0; i < len(reservations); i++ {
		var s1, s2 string
		if reservations[i].EndTime <= timeplans[now].End {
			if timeplans[now].Begin != reservations[i].BeginTime {
				s1 = timeplans[now].Begin
				s2 = reservations[i].BeginTime
			}
			if timeplans[now].End == reservations[i].EndTime {
				now++
			} else {
				timeplans[now].Begin = reservations[i].EndTime
			}
		} else {
			s1 = timeplans[now].Begin
			s2 = timeplans[now].End
			now++
			i--
		}
		if len(s1) != 0 && len(s2) != 0 {
			timeslices = timeslices.append(TimeSlice{Begin: s1, End: s2})
		}
	}
	for i := now; i < len(timeplans); i++ {
		timeslices = timeslices.append(TimeSlice{Begin: timeplans[i].Begin, End: timeplans[i].End})
	}
	st := 0
	ed := len(timeslices)-1
	for i := 0; i < len(timeslices); i++ {
		if begin >= timeslices[i].Begin && begin <= timeslices[i].End {
			if begin == timeslices[i].End {
				st = i+1
			} else {
				st = i
				timeslices[i].Begin = begin
			}
		}
		if end >= timeslices[i].Begin && end <= timeslices[i].End {
			if end == timeslices[i].Begin {
				ed = i-1
			} else {
				ed = i
				timeslices[i].End = end
			}
		}
	}
	return timeslices[st:ed]
}
