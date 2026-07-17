package input

import "github.com/google/uuid"

type CreateOrderInput struct {
	OrderUUID  uuid.UUID
	HullUUID   uuid.UUID
	UserUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}
