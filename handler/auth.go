package handler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/zhouziqunzzq/teacherAssessmentBackend/model"
	"github.com/yanzay/log"
	"github.com/dgrijalva/jwt-go"
	"time"
	. "github.com/zhouziqunzzq/teacherAssessmentBackend/config"
)

func Login(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	// Parse login form
	var user1 model.User
	var err1 error
	if user1.Username, err1 = getFormItemOrErr(w, req, "username", "用户名"); err1 != nil {
		return
	}
	if user1.Password, err1 = getFormItemOrErr(w, req, "password", "密码"); err1 != nil {
		return
	}
	// Check user in DB
	var user model.User
	user, err := model.GetUserByUsername(model.Db, user1.Username)
	if err != nil {
		if err.Error() == "GetUser: record not found" {
			responseJson(w,
				getErrorTpl(http.StatusNotFound, "用户名不存在"),
				http.StatusNotFound)
		} else {
			log.Error(err)
			responseJson(w,
				getErrorTpl(http.StatusInternalServerError, "登录失败，未知错误"),
				http.StatusInternalServerError)
		}
		return
	}
	// Check password
	if user1.Password != user.Password {
		responseJson(w,
			getErrorTpl(http.StatusNotFound, "密码错误"),
			http.StatusNotFound)
		return
	}
	// Generate JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(GlobalConfig.JWT_KEY))
	if err != nil {
		log.Error(err)
		responseJson(w,
			getErrorTpl(http.StatusInternalServerError, "登录失败，未知错误"),
			http.StatusInternalServerError)
		return
	}
	// Successfully logged in
	tpl := getOKTpl()
	tpl["msg"] = "登陆成功"
	tpl["token"] = tokenString
	responseJson(w, tpl, http.StatusOK)
}
