package main

import (
	"github.com/BurntSushi/toml"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"github.com/yanzay/log"
	. "github.com/zhouziqunzzq/teacherAssessmentBackend/config"
	"github.com/zhouziqunzzq/teacherAssessmentBackend/handler"
	"net/http"
	"strconv"
)

var mux = httprouter.New()

func initRouter() {
	mux.GET("/api", handler.Pong)
	mux.GET("/api/test/get", handler.Pong)
	mux.POST("/api/test/post", handler.PongPost)
	mux.NotFound = http.HandlerFunc(handler.NotFoundHandler)
}

func initCORS() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   GlobalConfig.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
	})
	h := c.Handler(mux)
	return h
}

func initMiddleware(h http.Handler) *negroni.Negroni {
	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(h)
	return n
}

func main() {
	// Load config from toml
	if _, err := toml.DecodeFile("config.toml", &GlobalConfig); err != nil {
		panic(err)
		return
	}
	// Init Router, CORS and Middleware
	initRouter()
	h := initCORS()
	n := initMiddleware(h)
	// Start the server
	log.Info("Starting HTTP server...")
	err := http.ListenAndServe(":"+strconv.FormatInt(GlobalConfig.PORT, 10), n)
	if err != nil {
		log.Fatal(err)
	}
}
