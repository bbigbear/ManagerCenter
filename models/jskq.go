package models

import (
	"time"
)

type Jskq struct {
	Id              int64
	TeacherId       string
	TeacherName     string
	AttendanceStyle string
	StartDate       time.Time
	StartTime       time.Time
	EndDate         time.Time
	EndTime         time.Time
	CreateDate      time.Time
	Num             int
	Remarks         string
}
