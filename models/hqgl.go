package models

import (
	"time"
)

type Qjq struct {
	Id      int64
	Campus  string
	Name    string
	Explain string
}

type Class struct {
	Id          int64
	ClassCode   string
	ClassName   string
	HeadTeacher string
}

//植被
type Tree struct {
	Id   int64
	Area string
	Name string
	Path string
	Date time.Time
}
