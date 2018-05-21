package model

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Reservation struct {
	ID            uint        `gorm:"AUTO_INCREMENT" json:"id"`
	User          User        `json:"user"`
	UserID        uint        `json:"-"`
	Meetingroom   Meetingroom `json:"meetingroom"`
	MeetingroomID uint        `json:"-"`
	Begin         string      `gorm:"type:varchar(20)" json:"begin"` // YYYY-MM-DD HH:MM:SS
	End           string      `gorm:"type:varchar(20)" json:"end"`   // YYYY-MM-DD HH:MM:SS
	BeginTime     string      `gorm:"-" json:"-"`                    // HH:MM:SS only used for calculation
	EndTime       string      `gorm:"-" json:"-"`                    // HH:MM:SS only used for calculation
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     *time.Time  `json:"-"`
}

func FillBeginTimeEndTime(reservations []Reservation) error {
	for i := 0; i < len(reservations); i++ {
		t, err := time.Parse("2006-01-02 15:04:05", reservations[i].Begin)
		if err != nil {
			return err
		}
		reservations[i].BeginTime = t.Format("15:04:05")
		t, err = time.Parse("2006-01-02 15:04:05", reservations[i].End)
		if err != nil {
			return err
		}
		reservations[i].EndTime = t.Format("15:04:05")
	}
	return nil
}

func GetReservationsByMeetingroomID(id uint, begin string, end string) []Reservation {
	var reservations []Reservation
	query := Db.Where("meetingroom_id = ?", id).
		Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("users.id, users.username, users.name, users.tel, users.type")
	})
	if len(begin) > 0 {
		query = query.Where("begin >= ?", begin)
	}
	if len(end) > 0 {
		query = query.Where("end <= ?", end)
	}
	query.Order("reservations.begin ASC").Find(&reservations)
	return reservations
}

func GetReservationsByUserIDWithHistorySwitch(uid uint, isHistory bool) []Reservation {
	var reservations []Reservation
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	nowStr := now.Format("2006-01-02 15:04:05")
	query := Db.Where("user_id = ?", uid).
		Preload("Meetingroom.Building")
	if isHistory {
		query = query.Where("end <= ?", nowStr)
	} else {
		query = query.Where("end > ?", nowStr)
	}
	query.Order("reservations.begin ASC").Find(&reservations)
	return reservations
}

//             begin|_________|end
//    r_begin|____________________|r_end
func GetReservationsContainingBeginEnd(begin, end string) []Reservation {
	var reservations []Reservation
	Db.Where("begin <= ?", begin).Where("end >= ?", end).
		Order("reservations.begin ASC").
		Find(&reservations)
	return reservations
}

//      begin|____________________|end
//            r_begin|_____|r_end
func GetReservationsInBeginEndWithMeetingroomID(id uint, begin, end string) []Reservation {
	var reservations []Reservation
	Db.Where("meetingroom_id = ?", id).
		Where("begin >= ?", begin).Where("end <= ?", end).
		Order("reservations.begin ASC").
		Find(&reservations)
	return reservations
}

//            now|  |nowAdvanced
//       r1_begin|_____|r1_end
//         r2_begin|_____|r2_end
//   r3_begin|________|r3_end
func GetReservationsInUseWithMeetingroomIDUserID(mid, uid uint, now, nowAdvanced string) []Reservation {
	var reservations []Reservation
	Db.Where("meetingroom_id = ?", mid).
		Where("user_id = ?", uid).
		Where("begin <= ?", nowAdvanced).
		Where("end >= ?", now).
		Order("reservations.begin ASC").
		Find(&reservations)
	return reservations
}
