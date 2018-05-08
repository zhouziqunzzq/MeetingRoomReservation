package main

import (
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"github.com/yanzay/log"
	. "github.com/zhouziqunzzq/MeetingRoomReservation/config"
	"github.com/zhouziqunzzq/MeetingRoomReservation/handler"
	"github.com/zhouziqunzzq/MeetingRoomReservation/model"
	"net/http"
	"strconv"
)

// Global var definition
var r = mux.NewRouter().StrictSlash(true)

func initDB() {
	sqliteDatabase, err := gorm.Open("sqlite3", GlobalConfig.SQLITE_FILE)
	if err != nil {
		panic(err)
	}
	model.Db = sqliteDatabase
	model.Db.AutoMigrate(&model.User{}, &model.Building{}, &model.Meetingroom{},
		&model.Weekplan{}, &model.Dayplan{}, &model.Timeplan{}, &model.Reservation{})
}

func initRouter() {
	// Static file
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(GlobalConfig.STATIC_DIR))))
	// public subrouters
	baseApiStr := "/api"
	baseApiVerStr := "/v1"
	baseStr := baseApiStr + baseApiVerStr
	api := r.PathPrefix(baseApiStr).Subrouter()
	v1Api := api.PathPrefix(baseApiVerStr).Subrouter()
	// Test
	api.Methods("GET").Path("/").HandlerFunc(handler.Pong)
	v1Api.Methods("GET").Path("/").HandlerFunc(handler.Pong)
	v1Api.Methods("GET").Path("/test").HandlerFunc(handler.Pong)
	v1Api.Methods("POST").Path("/test").HandlerFunc(handler.PongPost)
	// Authentication
	auth := v1Api.PathPrefix("/auth").Subrouter()
	auth.Methods("POST").Path("/login").HandlerFunc(handler.HandleLogin)
	// User
	userRoutes := handler.GetSubrouterWithMiddlewares(v1Api, baseStr,
		"/user", handler.ValidateTokenMiddleware)
	userRoutes.Methods("GET").Path("/info").HandlerFunc(handler.HandleGetUserInfo)
	// Meetingroom
	mrRoutes := v1Api.PathPrefix("/meetingroom").Subrouter()
	mrRoutes.Methods("GET").Path("/").HandlerFunc(handler.HandleGetMeetingroomList)
	mrRoutes.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(handler.HandleGetMeetingroomByID)
	mrRoutes.Methods("GET").Path("/{id:[0-9]+}/reservation").
		HandlerFunc(handler.HandleGetMeetingroomReservationsByID)
	mrRoutes.Methods("POST").Path("/{id:[0-9]+}/reservation").
		HandlerFunc(handler.HandlePostMeetingroomReservationByID)
	// NotFound
	r.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowedHandler)
}

func initCORS() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   GlobalConfig.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	})
	h := c.Handler(r)
	return h
}

func initGlobalMiddleware(h http.Handler) *negroni.Negroni {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
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
	n := initGlobalMiddleware(h)
	// Start the server
	log.Info("Starting HTTP server...")
	err := http.ListenAndServe(":"+strconv.FormatInt(GlobalConfig.PORT, 10), n)
	if err != nil {
		log.Fatal(err)
	}
}
