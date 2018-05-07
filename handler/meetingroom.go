package handler

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/zhouziqunzzq/MeetingRoomReservation/model"
	"strconv"
	"fmt"
)

func HandleGetMeetingroomList(w http.ResponseWriter, req *http.Request) {
	var meetingRooms []model.Meetingroom
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
			res := getErrorTpl(http.StatusBadRequest, "building_id参数错误")
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		fmt.Println(buildingID)
		query = query.Where("building_id = ?", buildingID)
	}
	if len(req.Form["floor"]) > 0 {
		floor, err = strconv.Atoi(req.Form["floor"][0])
		if err != nil {
			res := getErrorTpl(http.StatusBadRequest, "floor参数错误")
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		query = query.Where("floor = ?", floor)
	}
	if len(req.Form["begin"]) > 0 {
		begin = req.Form["begin"][0]
	}
	if len(req.Form["end"]) > 0 {
		end = req.Form["end"][0]
	}

	query.Find(&meetingRooms)
	res := getOKTpl()
	res["data"] = meetingRooms
	res["begin"] = begin
	res["end"] = end
	responseJson(w, res, http.StatusOK)
	return
}
