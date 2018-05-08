package handler

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/zhouziqunzzq/MeetingRoomReservation/model"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/yanzay/log"
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
