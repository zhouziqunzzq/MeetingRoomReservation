package main

import (
	"net/http"
	"github.com/zhouziqunzzq/MeetingRoomReservation/handler"
	. "github.com/zhouziqunzzq/MeetingRoomReservation/config"
)

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
	userRoutes.Methods("GET").Path("/reservation").HandlerFunc(handler.HandleGetUserReservationList)
	// Meetingroom
	mrRoutes := v1Api.PathPrefix("/meetingroom").Subrouter()
	mrRoutes.Methods("GET").Path("/").HandlerFunc(handler.HandleGetMeetingroomList)
	mrRoutes.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(handler.HandleGetMeetingroomByID)
	mrRoutes.Methods("POST").Path("/{id:[0-9]+}/unlock").
		HandlerFunc(handler.HandlePostMeetingroomUnlockByID)
	mrRoutes.Methods("GET").Path("/{id:[0-9]+}/reservation").
		HandlerFunc(handler.HandleGetMeetingroomReservationsByID)
	mrRoutes.Methods("POST").Path("/{id:[0-9]+}/reservation").
		HandlerFunc(handler.HandlePostMeetingroomReservationByID)
	// Admin
	adminRoutes := handler.GetSubrouterWithMiddlewares(v1Api, baseStr,
		"/admin", handler.ValidateAdminTokenMiddleware)
	adminRoutes.Methods("GET").Path("/").HandlerFunc(handler.Pong)
	adminRoutes.Methods("POST").Path("/meetingroom/{id:[0-9]+}/picture").HandlerFunc(handler.
		HandlePostPictureByMeetingroomID)
	// NotFound
	r.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowedHandler)
}
