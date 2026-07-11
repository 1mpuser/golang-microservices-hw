package model

import "time"

// OrderPaid — доменное представление входящего события «заказ оплачен».
// Приходит из Kafka (топик order.paid), producer — OrderService.
type OrderPaid struct {
	EventUUID string // для идемпотентности (защита от дубликатов)
	OrderUUID string
	UserUUID  string
}

// ShipAssembled — доменное событие «корабль собран», публикуется в Kafka
// (топик ship.assembled), consumer — OrderService.
type ShipAssembled struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64     // сколько секунд «собирали»
	AssembledAt  time.Time // время завершения сборки
}
