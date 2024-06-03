package query

import "github.com/Pr3c10us/gbt/internals/domain/debtor"

type GetDebtorByIDQuery interface {
	Handle(id string) (*debtor.Debtor, error)
}

type getDebtorByIDQuery struct {
	repository debtor.Repository
}

func (query *getDebtorByIDQuery) Handle(id string) (*debtor.Debtor, error) {
	d, err := query.repository.GetDebtorByID(id)
	if err != nil {
		return &debtor.Debtor{}, err
	}
	return d, nil
}

func NewGetDebtorByID(repository debtor.Repository) GetDebtorByIDQuery {
	return &getDebtorByIDQuery{
		repository: repository,
	}
}
