package order

type service struct {
	orderRepository OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	txManager       TxManager
}

func NewService(txManager TxManager, orderRepository OrderRepository, inventoryClient InventoryClient, paymentClient PaymentClient) *service {
	return &service{
		txManager:       txManager,
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
