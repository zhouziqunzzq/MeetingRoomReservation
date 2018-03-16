package main

import (
	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"github.com/yanzay/log"
	. "github.com/zhouziqunzzq/teacherAssessmentBackend/config"
	"github.com/zhouziqunzzq/teacherAssessmentBackend/handler"
	"github.com/zhouziqunzzq/teacherAssessmentBackend/model"
	"net/http"
	"strconv"
)

// Global var definition
var mux = httprouter.New()

func initDB() {
	sqliteDatabase, err := gorm.Open("sqlite3", GlobalConfig.SQLITE_FILE)
	if err != nil {
		panic(err)
	}
	model.Db = sqliteDatabase
	model.Db.AutoMigrate(&model.User{})
}

func initRouter() {
	// Test
	mux.GET("/api", handler.Pong)
	mux.GET("/api/v1", handler.Pong)
	mux.GET("/api/v1/test/get", handler.Pong)
	mux.POST("/api/v1/test/post", handler.PongPost)
	// Authentication
	mux.POST("/api/v1/auth/login", handler.Login)
	mux.NotFound = http.HandlerFunc(handler.NotFoundHandler)
}

func initCORS() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   GlobalConfig.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
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
	log.Info("Loading config from file...")
	if _, err := toml.DecodeFile("config.toml", &GlobalConfig); err != nil {
		panic(err)
		return
	}
	// Init database
	log.Info("Connecting to Database...")
	initDB()
	defer model.Db.Close()
	// Init Router, CORS, Middleware, OAuth
	log.Info("Initializing server...")
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
