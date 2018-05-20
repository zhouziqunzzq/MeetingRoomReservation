package handler

import (
	"net/http"
	"github.com/zhouziqunzzq/MeetingRoomReservation/model"
	"github.com/yanzay/log"
	"strconv"
)

func HandleGetUserInfo(w http.ResponseWriter, req *http.Request) {
	user, err := model.GetUserInfoByID(GetUIDFromJWT(req))
	if err != nil {
		log.Error(err)
		responseJson(w,
			getErrorTpl(http.StatusInternalServerError, "获取用户信息失败，未知错误"),
			http.StatusInternalServerError)
	} else {
		tpl := getOKTpl()
		tpl["data"] = user
		responseJson(w, tpl, http.StatusOK)
	}
}

func HandleGetUserReservationList(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	isHistoryStr, err1 := getFormItemOrErr(w, req, "is_history", "历史信息开关")
	if err1 != nil {
		return
	}
	isHistory, err2 := strconv.ParseBool(isHistoryStr)
	if err2 != nil {
		log.Error(err2.Error())
		res := getErrorTpl(http.StatusBadRequest, "is_history格式错误")
		responseJson(w, res, http.StatusBadRequest)
	}
	reservations := model.GetReservationsByUserIDWithHistorySwitch(GetUIDFromJWT(req), isHistory)
	res := getOKTpl()
	res["data"] = reservations
	responseJson(w, res, http.StatusOK)
}
