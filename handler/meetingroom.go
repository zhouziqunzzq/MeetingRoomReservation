package handler

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/zhouziqunzzq/MeetingRoomReservation/model"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/yanzay/log"
	"time"
	"github.com/pkg/errors"
	"github.com/zhouziqunzzq/MeetingRoomReservation/config"
	"github.com/zhouziqunzzq/MeetingRoomReservation/lockcontroller"
)

func HandleGetMeetingroomList(w http.ResponseWriter, req *http.Request) {
	query := model.Db.Preload("Weekplan.Dayplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("dayplans.weekday ASC")
	}).Preload("Weekplan.Dayplans.Timeplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("timeplans.begin ASC")
	}).Preload("Building")
	req.ParseForm()

	// Build up query with req params
	var isStrict bool
	var buildingID, floor int
	var err error
	var date, begin, end string
	var dateFmtStrLen = len("YYYY-MM-DD")
	var timeFmtStrLen = len("HH:MM:SS")
	// is_strict: required
	isStrictStr, err := getFormItemOrErr(w, req, "is_strict", "是否为严格模式")
	if err != nil {
		return
	}
	isStrict, err = strconv.ParseBool(isStrictStr)
	if err != nil {
		res := getErrorTpl(http.StatusBadRequest, "参数错误:is_strict")
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	// building_id
	if len(req.Form["building_id"]) > 0 {
		buildingID, err = strconv.Atoi(req.Form["building_id"][0])
		if err != nil {
			res := getErrorTpl(http.StatusBadRequest, "参数错误:building_id")
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		query = query.Where("building_id = ?", buildingID)
	}
	// floor
	if len(req.Form["floor"]) > 0 {
		floor, err = strconv.Atoi(req.Form["floor"][0])
		if err != nil {
			res := getErrorTpl(http.StatusBadRequest, "参数错误:floor")
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		query = query.Where("floor = ?", floor)
	}
	// date
	if len(req.Form["date"]) > 0 {
		date = req.Form["date"][0]
	}
	if len(date) != dateFmtStrLen && len(date) != 0 {
		res := getErrorTpl(http.StatusBadRequest, "date格式错误")
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	// begin
	if len(req.Form["begin"]) > 0 {
		begin = req.Form["begin"][0]
	} else {
		begin = "00:00:00"
	}
	if len(begin) != timeFmtStrLen {
		res := getErrorTpl(http.StatusBadRequest, "begin格式错误")
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	// end
	if len(req.Form["end"]) > 0 {
		end = req.Form["end"][0]
	} else {
		end = "23:59:59"
	}
	if len(end) != timeFmtStrLen {
		res := getErrorTpl(http.StatusBadRequest, "end格式错误")
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	if begin > end {
		res := getErrorTpl(http.StatusBadRequest, "参数错误:begin必须在end之前")
		responseJson(w, res, http.StatusBadRequest)
		return
	}

	// Build up meetingroom list and do filter
	var meetingrooms []model.Meetingroom
	var meetingroomsFiltered = make([]model.Meetingroom, 0)
	query.Find(&meetingrooms)
	for i := 0; i < len(meetingrooms); i++ {
		var avlTime []model.TimeSlice
		var err error
		if len(date) == dateFmtStrLen {
			avlTime, err = meetingrooms[i].GetAvlTimeWithDate(date, begin, end)
		} else if len(date) == 0 {
			avlTime, err = meetingrooms[i].GetAvlTimeWithDayCnt(begin, end)
		}
		if err != nil {
			log.Error(err)
			res := getErrorTpl(http.StatusInternalServerError, err.Error())
			responseJson(w, res, http.StatusInternalServerError)
			return
		}
		if isStrict {
			var avlTimeFiltered = make([]model.TimeSlice, 0)
			for j := 0; j < len(avlTime); j++ {
				if begin == avlTime[j].GetBeginTimeStr() &&
					end == avlTime[j].GetEndTimeStr() {
					avlTimeFiltered = append(avlTimeFiltered, avlTime[j])
				}
			}
			avlTime = avlTimeFiltered
		}
		if len(avlTime) > 0 {
			meetingrooms[i].AvlTime = avlTime
			meetingroomsFiltered = append(meetingroomsFiltered, meetingrooms[i])
		}
	}
	res := getOKTpl()
	res["data"] = meetingroomsFiltered
	responseJson(w, res, http.StatusOK)
	return
}

func GetMeetingroomIDFromUriOrError(w http.ResponseWriter, req *http.Request) (id int, err error) {
	vars := mux.Vars(req)
	id, err = strconv.Atoi(vars["id"])
	if err != nil {
		res := getErrorTpl(http.StatusNotFound, "会议室ID非法")
		responseJson(w, res, http.StatusNotFound)
		return
	}
	var meetingrooms []model.Meetingroom
	model.Db.Where("id = ?", id).Find(&meetingrooms)
	if len(meetingrooms) == 0 {
		res := getErrorTpl(http.StatusNotFound, "会议室ID不存在")
		responseJson(w, res, http.StatusNotFound)
		err = errors.New("Invalid Meetingroom ID")
		return
	}
	return
}

func HandleGetMeetingroomByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		res := getErrorTpl(http.StatusNotFound, "会议室ID不存在")
		responseJson(w, res, http.StatusNotFound)
		return
	}
	var meetingrooms []model.Meetingroom
	model.Db.Preload("Weekplan.Dayplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("dayplans.weekday ASC")
	}).Preload("Weekplan.Dayplans.Timeplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("timeplans.begin ASC")
	}).Preload("Building").Where("id = ?", id).Find(&meetingrooms)
	if len(meetingrooms) == 0 {
		res := getErrorTpl(http.StatusNotFound, "会议室ID不存在")
		responseJson(w, res, http.StatusNotFound)
		return
	}
	avlTime, err := meetingrooms[0].GetAvlTimeWithDayCnt("00:00:00", "23:59:59")
	if err != nil {
		log.Error(err)
		res := getErrorTpl(http.StatusInternalServerError, err.Error())
		responseJson(w, res, http.StatusNotFound)
		return
	}
	meetingrooms[0].AvlTime = avlTime
	res := getOKTpl()
	res["data"] = meetingrooms[0]
	responseJson(w, res, http.StatusOK)
	return
}

func HandleGetMeetingroomReservationsByID(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	id, err := GetMeetingroomIDFromUriOrError(w, req)
	if err != nil {
		return
	}

	var begin, end string
	if len(req.Form["begin"]) > 0 {
		begin = req.Form["begin"][0]
	}
	if len(req.Form["end"]) > 0 {
		end = req.Form["end"][0]
	}
	reservations := model.GetReservationsByMeetingroomID(uint(id), begin, end)
	res := getOKTpl()
	res["data"] = reservations
	responseJson(w, res, http.StatusOK)
	return
}

func HandlePostMeetingroomReservationByID(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	id, err := GetMeetingroomIDFromUriOrError(w, req)
	if err != nil {
		return
	}

	var begin, end string
	if _, err = getFormItemOrErr(w, req, "access_token", "token"); err != nil {
		if !ValidateToken(w, req) {
			return
		}
	}
	if begin, err = getFormItemOrErr(w, req, "begin", "开始时间"); err != nil {
		return
	}
	if end, err = getFormItemOrErr(w, req, "end", "结束时间"); err != nil {
		return
	}

	// Check if reservation overlaps with others
	reservations := model.GetReservationsContainingBeginEnd(begin, end)
	beginTime, err := time.Parse("2006-01-02 15:04:05", begin)
	endTime, err := time.Parse("2006-01-02 15:04:05", end)
	if err != nil || beginTime.Year() != endTime.Year() ||
		beginTime.Month() != endTime.Month() || beginTime.Day() != endTime.Day() ||
		len(reservations) > 0 {
		res := getErrorTpl(http.StatusBadRequest, "开始时间或结束时间非法或该时间段已被占用")
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	// Check if reservation is in any of timeplans
	meetingroom := model.GetMeetingroomByID(uint(id))
	ok := false
	weekplan := meetingroom.Weekplan.ConvertToMap()
	for i := 0; i < len(weekplan[beginTime.Weekday()]); i++ {
		if weekplan[beginTime.Weekday()][i].Begin <= beginTime.Format("15:04:05") &&
			weekplan[endTime.Weekday()][i].End >= endTime.Format("15:04:05") {
			ok = true
		}
	}
	if !ok {
		res := getErrorTpl(http.StatusBadRequest, "该会议室在指定的时间段内不开放")
		responseJson(w, res, http.StatusBadRequest)
		return
	}

	// Add new reservation
	uid := GetUIDFromJWT(req)
	reservation := model.Reservation{
		UserID:        uid,
		MeetingroomID: uint(id),
		Begin:         begin,
		End:           end,
	}
	model.Db.Create(&reservation)
	res := getOKTpl()
	res["data"] = map[string]interface{}{
		"reservation_id": reservation.ID,
	}
	responseJson(w, res, http.StatusOK)
	return
}

func HandlePostPictureByMeetingroomID(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	id, err := GetMeetingroomIDFromUriOrError(w, req)
	if err != nil {
		return
	}

	if len(req.Form["picture"]) == 0 {
		responseJson(w, getErrorTpl(http.StatusBadRequest, "图片格式错误"), http.StatusBadRequest)
	}
	pictureEncoded := req.Form["picture"][0]
	log.Debug(id)
	log.Debug(pictureEncoded)
	// TODO: Save picture to STATIC_DIR/meetingroom/picture/{id}
}

func HandlePostMeetingroomUnlockByID(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	mid, err := GetMeetingroomIDFromUriOrError(w, req)
	if err != nil {
		return
	}
	if _, err = getFormItemOrErr(w, req, "access_token", "token"); err != nil {
		if !ValidateToken(w, req) {
			return
		}
	}
	uid := GetUIDFromJWT(req)

	// Check user's reservation
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	nowStr := now.Format("2006-01-02 15:04:05")
	nowAdvanced := now.Add(time.Minute * time.Duration(config.GlobalConfig.LOCK_ADVANCED_MINUTE))
	nowAdvancedStr := nowAdvanced.Format("2006-01-02 15:04:05")
	reservations := model.GetReservationsInUseWithMeetingroomIDUserID(uint(mid), uint(uid),
		nowStr, nowAdvancedStr)

	if len(reservations) <= 0 {
		res := getOKTpl()
		res["result"] = false
		res["msg"] = "您没有预定本时间段的会议室，请检查"
		responseJson(w, res, http.StatusOK)
		return
	} else {
		meetingroom := model.GetMeetingroomByID(uint(mid))
		err := lockcontroller.Unlock(meetingroom.IP)
		if err != nil {
			res := getOKTpl()
			res["result"] = false
			res["msg"] = "会议室解锁失败，请重试"
			responseJson(w, res, http.StatusOK)
		} else {
			res := getOKTpl()
			res["msg"] = "会议室解锁成功"
			responseJson(w, res, http.StatusOK)
		}
	}
	return
}
