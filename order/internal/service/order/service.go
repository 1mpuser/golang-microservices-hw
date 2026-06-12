package order

type service struct {
	orderRepository OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

func NewService(orderRepository OrderRepository, inventoryClient InventoryClient, paymentClient PaymentClient) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
