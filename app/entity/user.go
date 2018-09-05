package entity

import (
	"chainstack/models"
	"chainstack/repo"
	"chainstack/utilities/uer"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userEntity struct {
	userRepo repo.IUser
}

type User interface {
	Login(username, password string) (*models.User, error)
	ListUsers(nextId, limit int) ([]models.User, error)
}

func NewUser(userRepo repo.IUser) User {
	return &userEntity{
		userRepo: userRepo,
	}
}

func (u userEntity) Login(email, password string) (*models.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, uer.InternalError(err)
	}

	if user == nil {
		// return nil, uer.NotFoundError(errors.New("user not found"))
		return nil, uer.NotAuthorizedError(errors.New("Wrong email or password"))
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Salt+password)); err != nil {
		return nil, uer.NotAuthorizedError(errors.New("Wrong email or password"))
	}

	return user, nil
}

func (u userEntity) ListUsers(nextId, limit int) (users []models.User, err error) {
	users, err = u.userRepo.List(nextId, limit)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	return
}
