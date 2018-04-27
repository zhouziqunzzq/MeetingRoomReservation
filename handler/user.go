package handler

import (
	"net/http"
	"github.com/zhouziqunzzq/MeetingRoomReservation/model"
	"github.com/yanzay/log"
)

func HandleGetUserInfo(w http.ResponseWriter, req *http.Request) {
	user, err := model.GetUserInfoByID(model.Db, GetUIDFromJWT(req))
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
