package part

type service struct {
	partRepository PartRepository
	checker        CompatibilityChecker
	txManager      TxManager
}

func NewService(partRepository PartRepository, checker CompatibilityChecker, txManager TxManager) *service {
	return &service{
		partRepository,
		checker,
		txManager,
	}
}
