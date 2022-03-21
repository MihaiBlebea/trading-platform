package pos

import (
	"errors"
)

type PositionRepoMock struct {
	positions map[int]Position
}

func (pr *PositionRepoMock) Save(pos *Position) (*Position, error) {
	id := len(pr.positions) + 1
	pr.positions[id] = *pos
	pos.ID = id

	return pos, nil
}

func (pr *PositionRepoMock) Update(pos *Position) error {
	if _, exists := pr.positions[pos.ID-1]; !exists {
		return errors.New("could not find index")
	}

	pr.positions[pos.ID-1] = *pos

	return nil
}

func (pr *PositionRepoMock) WithAccountAndSymbol(accountId int, symbol string) (*Position, error) {
	for _, pos := range pr.positions {
		if pos.AccountID == accountId && pos.Symbol == symbol {
			return &pos, nil
		}
	}

	return &Position{}, errors.New("could not find record")
}

func (pr *PositionRepoMock) WithAccountId(accountId int) ([]Position, error) {
	positions := []Position{}
	for _, pos := range pr.positions {
		if pos.AccountID == accountId {
			positions = append(positions, pos)
		}
	}

	return positions, nil
}

func (pr *PositionRepoMock) Delete(pos *Position) error {
	delete(pr.positions, pos.ID)

	return nil
}
