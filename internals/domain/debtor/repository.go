package debtor

type Repository interface {
	AddDebtor(debtor *Debtor) error
	RemoveDebtor(id string) error
	EditDebtor(id string, data ...struct {
		key   string
		value any
	}) error

	GetDebtor(id string) Debtor
}
