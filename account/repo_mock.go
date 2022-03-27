package account

import "errors"

type AccountRepoMock struct {
	accounts []Account
}

func (ar *AccountRepoMock) Save(account *Account) (*Account, error) {
	account.ID = len(ar.accounts) + 1
	ar.accounts = append(ar.accounts, *account)

	return account, nil
}

func (ar *AccountRepoMock) Update(account *Account) error {
	if len(ar.accounts) < account.ID-1 {
		return errors.New("could not find index")
	}

	ar.accounts[account.ID-1] = *account

	return nil
}

func (ar *AccountRepoMock) WithToken(token string) (*Account, error) {
	for _, acc := range ar.accounts {
		if acc.ApiToken == token {
			return &acc, nil
		}
	}

	return &Account{}, errors.New("could not find record")
}

func (ar *AccountRepoMock) WithId(id int) (*Account, error) {
	if len(ar.accounts) < id-1 {
		return &Account{}, errors.New("could not find index")
	}

	return &ar.accounts[id-1], nil
}
