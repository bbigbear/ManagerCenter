package models

import (
	"time"
)

type Czlx struct {
	Id      int64
	Name    string
	Explain string
}

type Jzdj struct {
	Id             int64
	Number         string
	Name           string
	Origin         string
	HandlePeople   string
	Sum            float64
	IncomeData     time.Time
	Remarks        string
	RegisterPeople string
	RegisterData   time.Time
}

type Czdj struct {
	Id             int64
	Number         string
	Name           string
	Style          string
	HandlePeople   string
	Sum            float64
	PayData        time.Time
	Remarks        string
	RegisterPeople string
	RegisterData   time.Time
}
