package model

type TimeSlice struct {
	Begin string `json:"begin"` // YYYY-MM-DD HH:MM:SS
	End   string `json:"end"`   // YYYY-MM-DD HH:MM:SS
}

func (ts *TimeSlice) GetBeginDateStr() string {
	return ts.Begin[0:len("YYYY-MM-DD")]
}

func (ts *TimeSlice) GetEndDateStr() string {
	return ts.End[0:len("YYYY-MM-DD")]
}

func (ts *TimeSlice) GetBeginTimeStr() string {
	return ts.Begin[len("YYYY-MM-DD")+1:]
}

func (ts *TimeSlice) GetEndTimeStr() string {
	return ts.End[len("YYYY-MM-DD")+1:]
}
