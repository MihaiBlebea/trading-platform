package pie

import (
	"errors"

	"gorm.io/gorm"
)

type PieRepo struct {
	conn *gorm.DB
}

func NewPieRepo(conn *gorm.DB) (*PieRepo, error) {
	if err := conn.AutoMigrate(&Pie{}, &PieSlice{}); err != nil {
		return &PieRepo{}, err
	}

	return &PieRepo{conn: conn}, nil
}

func (pr *PieRepo) Save(pos *Pie) (*Pie, error) {
	resp := pr.conn.Create(pos)

	if resp.Error != nil {
		return &Pie{}, resp.Error
	}
	return pos, nil
}

func (pr *PieRepo) WithId(id int) (*Pie, error) {
	var pie Pie
	err := pr.conn.Where("id = ?", id).Find(&pie).Error
	if err != nil {
		return &Pie{}, err
	}

	return &pie, err
}

// Hooks

func (p *Pie) AfterCreate(tx *gorm.DB) error {
	if len(p.Slices) == 0 {
		return errors.New("no slices in this pie")
	}

	for _, slice := range p.Slices {
		if slice.ID != 0 {
			continue
		}

		slice.PieID = p.ID

		if err := tx.Create(&slice).Error; err != nil {
			return err
		}
	}

	return nil
}

func (p *Pie) AfterFind(tx *gorm.DB) error {
	var slices []PieSlice
	if err := tx.Where("pie_id = ?", p.ID).Find(&slices).Error; err != nil {
		return err
	}
	p.Slices = slices

	return nil
}
