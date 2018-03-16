package handler

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
	return
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 Method Not Allowed"))
	return
}

func Pong(w http.ResponseWriter, req *http.Request) {
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "OK",
	}
	responseJson(w, res, http.StatusOK)
	return
}

func PongPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   req.Form,
	}
	responseJson(w, res, http.StatusOK)
	return
}
