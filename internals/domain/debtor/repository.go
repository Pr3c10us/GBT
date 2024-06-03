package debtor

type Repository interface {
	AddDebtor(debtor *Debtor) error
	RemoveDebtor(id string) error
	GetDebtorByID(id string) (*Debtor, error)
}
