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
)

func HandleGetMeetingroomList(w http.ResponseWriter, req *http.Request) {
	var meetingrooms []model.Meetingroom
	query := model.Db.Preload("Weekplan.Dayplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("dayplans.weekday ASC")
	}).Preload("Weekplan.Dayplans.Timeplans", func(db *gorm.DB) *gorm.DB {
		return db.Order("timeplans.begin ASC")
	}).Preload("Building")
	req.ParseForm()

	// Build up query with req params
	var buildingID, floor int
	var err error
	var begin, end string
	if len(req.Form["building_id"]) > 0 {
		buildingID, err = strconv.Atoi(req.Form["building_id"][0])
		if err != nil {
			res := getErrorTpl(http.StatusBadRequest, "参数错误:building_id")
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		query = query.Where("building_id = ?", buildingID)
	}
	if len(req.Form["floor"]) > 0 {
		floor, err = strconv.Atoi(req.Form["floor"][0])
		if err != nil {
			res := getErrorTpl(http.StatusBadRequest, "参数错误:floor")
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		query = query.Where("floor = ?", floor)
	}
	if len(req.Form["begin"]) > 0 {
		begin = req.Form["begin"][0]
	} else {
		begin = "00:00:00"
	}
	if len(req.Form["end"]) > 0 {
		end = req.Form["end"][0]
	} else {
		end = "23:59:59"
	}
	if begin > end {
		res := getErrorTpl(http.StatusBadRequest, "参数错误:begin必须在end之前")
		responseJson(w, res, http.StatusBadRequest)
		return
	}

	query.Find(&meetingrooms)
	for i := 0; i < len(meetingrooms); i++ {
		err = meetingrooms[i].GetAvlTime(begin, end)
		if err != nil {
			log.Error(err)
			res := getErrorTpl(http.StatusInternalServerError, err.Error())
			responseJson(w, res, http.StatusNotFound)
			return
		}
	}
	res := getOKTpl()
	res["data"] = meetingrooms
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
	err = meetingrooms[0].GetAvlTime("00:00:00", "23:59:59")
	if err != nil {
		log.Error(err)
		res := getErrorTpl(http.StatusInternalServerError, err.Error())
		responseJson(w, res, http.StatusNotFound)
		return
	}
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
	reservations := model.GetReservationsWithBeginEnd(begin, end)
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
