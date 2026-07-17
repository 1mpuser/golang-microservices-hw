package model

import "time"

type OrderPaid struct {
	EventUUID string
	OrderUUID string
	UserUUID  string
}

type ShipAssembled struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	AssembledAt  time.Time
}
