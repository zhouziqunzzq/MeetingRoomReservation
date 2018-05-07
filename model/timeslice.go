package model

type TimeSlice struct {
	Begin string `json:"begin"` // YYYY-MM-DD HH:MM:SS
	End   string `json:"end"`   // YYYY-MM-DD HH:MM:SS
}
