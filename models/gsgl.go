package models

import (
	"time"
)

type Tax struct {
	Id    int64
	Name  string
	Style string
	Money float64
	Sort  int
}

type TaxMonth struct {
	Id         int64
	Month      int
	Project    string
	FontRange  time.Time
	AfterRange time.Time
}
