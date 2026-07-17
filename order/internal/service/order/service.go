package order

type service struct {
	orderRepository OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	txManager       TxManager
	orderProducer   OrderProducer
}

func NewService(txManager TxManager, orderRepository OrderRepository, inventoryClient InventoryClient, paymentClient PaymentClient, orderProducer OrderProducer) *service {
	return &service{
		txManager:       txManager,
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderProducer:   orderProducer,
	}
}
