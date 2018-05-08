package model

import (
	"time"
	"github.com/getlantern/deepcopy"
	"github.com/zhouziqunzzq/MeetingRoomReservation/config"
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
func (m *Meetingroom) GetAvlTime(begin, end string) (err error) {
	dayCnt := config.GlobalConfig.MAX_QUERY_DAY
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	nowDateStr := now.Format("2006-01-02")
	nowTimeStr := now.Format("15:04:05")

	weekplanMap := m.Weekplan.ConvertToMap()
	calc, err := time.Parse("2006-01-02", now.Format("2006-01-02"))
	if err != nil {
		return
	}
	// Calc first day
	if timeplans, ok := weekplanMap[calc.Weekday()]; ok && nowTimeStr < end {
		var reservations []Reservation
		beginTmp := begin
		if nowTimeStr > begin {
			beginTmp = nowTimeStr
		}
		Db.Where("meetingroom_id = ?", m.ID).
			Where("begin > ?", nowDateStr+" 00:00:00").
			Where("end < ?", nowDateStr+" 23:59:59").
			Order("reservations.begin ASC").
			Find(&reservations)
		err = FillBeginTimeEndTime(reservations)
		if err != nil {
			return
		}
		m.AvlTime = append(m.AvlTime, GetAvlTimeOneDay(timeplans, reservations,
			beginTmp, end, nowDateStr)...)
	}
	// Calc rest of days
	for i := 0; i < dayCnt-1; i++ {
		calc = calc.AddDate(0, 0, 1)
		// Check dayplan of calc.Weekday()
		timeplans, ok := weekplanMap[calc.Weekday()]
		if !ok {
			continue
		}

		// Get reservations of this day
		var reservations []Reservation
		calcDateStr := calc.Format("2006-01-02")
		Db.Where("meetingroom_id = ?", m.ID).
			Where("begin > ?", calcDateStr+" 00:00:00").
			Where("end < ?", calcDateStr+" 23:59:59").
			Order("reservations.begin ASC").
			Find(&reservations)
		err = FillBeginTimeEndTime(reservations)
		if err != nil {
			return
		}

		// Call GetAvlTimeOneDay and append it to final results
		m.AvlTime = append(m.AvlTime, GetAvlTimeOneDay(timeplans, reservations,
			begin, end, calcDateStr)...)
	}
	return
}

// timeplans: a slice of Timeplans in a day
// reservations: a slice of Reservations in a day
// begin: HH:MM:SS
// end: HH:MM:SS
// date: YYYY-MM-DD
func GetAvlTimeOneDay(timeplansOriginal []Timeplan, reservations []Reservation,
	begin string, end string, date string) []TimeSlice {
	// Make a deep copy of timeplansOriginal first
	var timeplans []Timeplan
	deepcopy.Copy(&timeplans, &timeplansOriginal)
	timeslices := make([]TimeSlice, 0)
	now := 0
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
			timeslices = append(timeslices, TimeSlice{Begin: s1, End: s2})
		}
	}
	for i := now; i < len(timeplans); i++ {
		timeslices = append(timeslices, TimeSlice{Begin: timeplans[i].Begin, End: timeplans[i].End})
	}
	st := -1
	stLock := false
	ed := -1
	edLock := false
	for i := 0; i < len(timeslices); i++ {
		if !stLock && begin <= timeslices[i].Begin && begin <= timeslices[i].End {
			st = i
			stLock = true
		} else if begin >= timeslices[i].Begin && begin <= timeslices[i].End {
			if begin == timeslices[i].End {
				st = i + 1
			} else {
				st = i
				timeslices[i].Begin = begin
			}
		}
		if !edLock && end >= timeslices[i].Begin && end >= timeslices[i].End {
			ed = i
		} else if end >= timeslices[i].Begin && end <= timeslices[i].End {
			if end == timeslices[i].Begin {
				ed = i - 1
			} else {
				ed = i
				timeslices[i].End = end
			}
			edLock = true
		}
	}
	if st == -1 || ed == -1 {
		return []TimeSlice{}
	} else {
		for i := st; i < ed+1; i++ {
			timeslices[i].Begin = date + " " + timeslices[i].Begin
			timeslices[i].End = date + " " + timeslices[i].End
		}
		return timeslices[st : ed+1]
	}
}
