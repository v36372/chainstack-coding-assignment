package repo

import (
	"chainstack/infra"
	"errors"
)

type base struct{}

func (base) create(value interface{}) (interface{}, error) {
	if value != nil {
		err := infra.PostgreSql.Create(value).Error
		return value, err
	}
	return nil, errors.New("Create failed. Value is nil")
}

func (base) save(value interface{}) error {
	if value != nil {
		return infra.PostgreSql.Save(value).Error
	}
	return errors.New("Update failed. Value is nil")
}

func (base) delete(value interface{}) error {
	if value != nil {
		return infra.PostgreSql.Delete(value).Error
	}
	return errors.New("Delete failed. Value is nil")
}
