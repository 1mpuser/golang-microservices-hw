package orderpaid

import (
	"context"
	"sync"

	"github.com/1mpuser/platform/pkg/kafka"
)

// Handler — обработчик входящих Kafka-сообщений OrderPaid.
type Handler struct {
	service AssemblyService
	seen    map[string]struct{}
	mu      sync.Mutex
}

// NewHandler создаёт обработчик события OrderPaid.
func NewHandler(service AssemblyService) *Handler {
	return &Handler{
		service: service,
		seen:    make(map[string]struct{}),
	}
}

// Handle реализует kafka.MessageHandler.
func (h *Handler) Handle(ctx context.Context, msg kafka.Message) error {
	event, err := decodeOrderPaid(msg.Value)

	if err != nil {
		return err
	}

	eventUuid := event.EventUUID

	h.mu.Lock()
	_, ok := h.seen[eventUuid]

	if ok {
		h.mu.Unlock()
		return nil
	}

	h.mu.Unlock()

	err = h.service.Assemble(ctx, event)

	if err != nil {
		return err
	}

	h.mu.Lock()
	h.seen[eventUuid] = struct{}{}
	h.mu.Unlock()

	return nil
}
