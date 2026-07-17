package orderpaid

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/1mpuser/assembly/internal/model"
	eventsv1 "github.com/1mpuser/shared/pkg/proto/events/v1"
)

func decodeOrderPaid(data []byte) (model.OrderPaid, error) {
	var pb eventsv1.OrderPaid

	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaid{}, fmt.Errorf("десериализовать свойства: %w", err)
	}

	return model.OrderPaid{
		UserUUID:  pb.GetUserUuid(),
		EventUUID: pb.GetEventUuid(),
		OrderUUID: pb.GetOrderUuid(),
	}, nil
}
