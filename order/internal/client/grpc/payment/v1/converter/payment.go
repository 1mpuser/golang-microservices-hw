package converter

type TransactionUUID struct {
	TransactionUUID string
}

func DTOToModel(transactionUUID string) *TransactionUUID {
	return &TransactionUUID{
		TransactionUUID: transactionUUID,
	}
}
