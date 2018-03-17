package handler

import (
	"net/http"
	"encoding/json"
	"github.com/yanzay/log"
	"github.com/pkg/errors"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/dgrijalva/jwt-go"
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

func getOKTpl() (map[string]interface{}) {
	return map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "success",
	}
}

func getErrorTpl(code int, msg string) (map[string]interface{}) {
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

/* Generates a negroni handler for the route:
    pathType (base+string)
    which is handled by `f` and passed through
    middleware `mids` sequentially.

    ex:
    NegroniRoute(router, "/api/v1", "/users/update", "POST", UserUpdateHandler, LoggingMiddleware, AuthorizationMiddleware)
*/
func GetSubrouterWithMiddlewares(
	parent *mux.Router,
	base string,
	path string,
	mw ...func(http.ResponseWriter, *http.Request, http.HandlerFunc), // Middlewares
) (*mux.Router) {
	baseRouter := mux.NewRouter()
	n := negroni.New()
	for i := range mw {
		n.Use(negroni.HandlerFunc(mw[i]))
	}
	n.UseHandler(baseRouter)
	parent.PathPrefix(path).Handler(n)
	return baseRouter.PathPrefix(base + path).Subrouter()
}

func GetUIDFromJWT(req *http.Request)(uint) {
	req.ParseForm()
	p := jwt.Parser{
		UseJSONNumber:false,
		SkipClaimsValidation:true,
	}
	t, _, _ := p.ParseUnverified(req.Form["access_token"][0], jwt.MapClaims{})
	return uint(t.Claims.(jwt.MapClaims)["uid"].(float64))
}