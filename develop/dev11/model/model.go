package model

import (
	"time"
)

type ISOtime time.Time

type Db struct {
	//todo mutex
	Storage map[int]Event
	Index   int
}

type Event struct {
	UserId string  `json:"user_id"`
	Id     int     `json:"id"`
	Date   ISOtime `json:"date"`
	Name   string  `json:"name"`
}

type Result struct {
	Res Event `json:"result"`
}

type AggrResult struct {
	Res []Event `json:"result"`
}

type ErrorReuslt struct {
	Err string `json:"error"`
}
