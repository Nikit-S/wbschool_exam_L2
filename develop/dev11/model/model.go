package model

import (
	"time"
)

type ISOtime time.Time

type Db struct {
	Storage map[int]Event
	Index   int
}

type Event struct {
	UserId string  `json:"user_id" validate:"required"`
	Id     int     `json:"id"`
	Date   ISOtime `json:"date" validate:"required"`
	Name   string  `json:"name" validate:"required"`
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
