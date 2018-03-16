package handler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	return
}

func Pong(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "OK",
	}
	responseJson(w, res, http.StatusOK)
	return
}

func PongPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   req.Form,
	}
	responseJson(w, res, http.StatusOK)
	return
}
