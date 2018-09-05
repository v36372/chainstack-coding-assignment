package repo

import (
	"chainstack/infra"
	"chainstack/models"

	"github.com/jinzhu/gorm"
)

type user struct {
	base
}

var User IUser

func init() {
	User = user{}
}

type IUser interface {
	GetByEmail(email string) (*models.User, error)
	Create(*models.User) error
	Update(*models.User) error
	Delete(*models.User) error
}

func (u user) Create(user *models.User) error {
	return u.create(user)
}

func (u user) Update(user *models.User) error {
	return u.save(user)
}

func (u user) Delete(user *models.User) error {
	return u.delete(user)
}

func (user) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := infra.PostgreSql.Model(models.User{}).
		Where("email = ?", email).
		Limit(1).
		Find(&user).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &user, err
}
