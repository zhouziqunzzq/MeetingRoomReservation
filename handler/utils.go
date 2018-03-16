package handler

import (
	"net/http"
	"encoding/json"
	"github.com/yanzay/log"
	"github.com/pkg/errors"
)

func responseJson(w http.ResponseWriter, data map[string]interface{}, httpStatusCode int) {
	resJson, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred while encoding json response.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(httpStatusCode)
	w.Write(resJson)
	return
}

func getOKTpl()(map[string]interface{}){
	return map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "success",
	}
}

func getErrorTpl(code int, msg string)(map[string]interface{}){
	return map[string]interface{}{
		"code":   code,
		"result": false,
		"msg":    msg,
	}
}

// Check if given item in given form exists, if not response error message
// in JSON. Error message will use item name if description is empty string.
func getFormItemOrErr(w http.ResponseWriter, req *http.Request,
	item string, description string) (value string, err error) {
	if len(req.Form[item]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
		}
		if description == "" {
			res["msg"] = item + "错误"
		} else {
			res["msg"] = description + "错误"
		}
		responseJson(w, res, http.StatusBadRequest)
		value = ""
		err = errors.New("Fail to get form item")
	} else {
		value = req.Form[item][0]
		err = nil
	}
	return
}
