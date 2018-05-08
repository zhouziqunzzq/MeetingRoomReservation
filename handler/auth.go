package handler

import (
	"net/http"
	"github.com/zhouziqunzzq/MeetingRoomReservation/model"
	"github.com/yanzay/log"
	"github.com/dgrijalva/jwt-go"
	"time"
	. "github.com/zhouziqunzzq/MeetingRoomReservation/config"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/context"
)

func HandleLogin(w http.ResponseWriter, req *http.Request) {
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
				http.StatusOK)
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
			http.StatusOK)
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
	tpl["access_token"] = tokenString
	responseJson(w, tpl, http.StatusOK)
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(GlobalConfig.JWT_KEY), nil
		})
	if err != nil {
		responseJson(w, getErrorTpl(http.StatusUnauthorized, "未授权的访问"),
			http.StatusUnauthorized)
		return
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set(r, "uid", claims["uid"])
			next(w, r)
		} else {
			responseJson(w, getErrorTpl(http.StatusUnauthorized, "无效的Token"),
				http.StatusUnauthorized)
			return
		}
	}
}

func ValidateToken(w http.ResponseWriter, r *http.Request) bool {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(GlobalConfig.JWT_KEY), nil
		})
	if err != nil {
		responseJson(w, getErrorTpl(http.StatusUnauthorized, "未授权的访问"),
			http.StatusUnauthorized)
		return false
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set(r, "uid", claims["uid"])
			return true
		} else {
			responseJson(w, getErrorTpl(http.StatusUnauthorized, "无效的Token"),
				http.StatusUnauthorized)
			return false
		}
	}
}
