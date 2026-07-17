package assemblyconsumer

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/1mpuser/order/internal/model"
	eventsv1 "github.com/1mpuser/shared/pkg/proto/events/v1"
)

// decodeShipAssembled десериализует protobuf-сообщение ShipAssembled из Kafka в доменную модель.
func decodeShipAssembled(data []byte) (model.ShipAssembled, error) {
	var pb eventsv1.ShipAssembled

	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembled{}, fmt.Errorf("десериализовать ShipAssembled: %w", err)
	}

	return model.ShipAssembled{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
		AssembledAt:  pb.GetAssembledAt().AsTime(),
	}, nil
}
