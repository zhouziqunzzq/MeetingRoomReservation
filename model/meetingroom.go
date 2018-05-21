package model

import (
	"time"
	"github.com/getlantern/deepcopy"
	"github.com/zhouziqunzzq/MeetingRoomReservation/config"
	"github.com/jinzhu/gorm"
)

type Meetingroom struct {
	ID           uint                        `gorm:"AUTO_INCREMENT" json:"id"`
	Building     Building                    `json:"building"`
	BuildingID   uint                        `json:"-"`
	Floor        int                         `json:"floor"`
	Room         string                      `gorm:"type:varchar(16)" json:"room"`
	Weekplan     Weekplan                    `json:"-"`
	WeekplanID   uint                        `json:"-"`
	IP           string                      `gorm:"type:varchar(64)" json:"ip"`
	WeekplanMap  map[time.Weekday][]Timeplan `gorm:"-" json:"-"`
	Reservations []Reservation               `json:"reservations"`
	AvlTime      []TimeSlice                 `gorm:"-" json:"avl_time"`
}

// GetAvlTimeWithDate will calculate available time slices from date-begin to date-end
// Before calling this func, make sure that Weekplan
// and Dayplans of the Weekplan and Timeplans of the Dayplans are set correctly.
// date: YYYY-MM-DD
// begin: HH:MM:SS
// end: HH:MM:SS
func (m *Meetingroom) GetAvlTimeWithDate(date, begin, end string) (avlTime []TimeSlice, err error) {
	avlTime = make([]TimeSlice, 0)
	// parse date
	calc, err := time.Parse("2006-01-02", date)
	if err != nil {
		return
	}
	// get timeplans from weekplanMap
	if m.WeekplanMap == nil {
		m.WeekplanMap = m.Weekplan.ConvertToMap()
	}
	timeplans, ok := m.WeekplanMap[calc.Weekday()]
	if !ok {
		return
	}
	// get reservations on that day
	reservations := GetReservationsInBeginEndWithMeetingroomID(m.ID, date+" 00:00:00", date+" 23:59:59")
	err = FillBeginTimeEndTime(reservations)
	if err != nil {
		return
	}

	avlTime = GetAvlTimeOneDay(timeplans, reservations, begin, end, date)
	return
}

// GetAvlTimeWithDayCnt will calculate available time slices from now to the midnight of
// the day after "dayCnt" days. Before calling this func, make sure that Weekplan
// and Dayplans of the Weekplan and Timeplans of the Dayplans are set correctly.
// begin: HH:MM:SS
// end: HH:MM:SS
func (m *Meetingroom) GetAvlTimeWithDayCnt(begin, end string) (avlTime []TimeSlice, err error) {
	avlTime = make([]TimeSlice, 0)
	var tmpAvlTime []TimeSlice
	dayCnt := config.GlobalConfig.MAX_QUERY_DAY
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	nowDateStr := now.Format("2006-01-02")
	nowTimeStr := now.Format("15:04:05")

	if m.WeekplanMap == nil {
		m.WeekplanMap = m.Weekplan.ConvertToMap()
	}
	calc, err := time.Parse("2006-01-02", now.Format("2006-01-02"))
	if err != nil {
		return
	}
	// Calc first day
	if nowTimeStr < end {
		tmpAvlTime, err = m.GetAvlTimeWithDate(nowDateStr, nowTimeStr, end)
		if err != nil {
			return
		}
		avlTime = append(avlTime, tmpAvlTime...)
	}
	// Calc rest of days
	for i := 0; i < dayCnt-1; i++ {
		calc = calc.AddDate(0, 0, 1)
		calcDateStr := calc.Format("2006-01-02")
		tmpAvlTime, err = m.GetAvlTimeWithDate(calcDateStr, begin, end)
		if err != nil {
			return
		}
		avlTime = append(avlTime, tmpAvlTime...)
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
			stLock = true
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

func GetMeetingroomByID(id uint) (m Meetingroom) {
	Db.Preload("Weekplan.Dayplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("dayplans.weekday ASC")
	}).Preload("Weekplan.Dayplans.Timeplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("timeplans.begin ASC")
	}).Preload("Building").Where("id = ?", id).Find(&m)
	return
}
