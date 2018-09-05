package entity

import (
	"chainstack/models"
	"chainstack/repo"
	"chainstack/utilities/uer"
	"errors"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

var letterBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type userEntity struct {
	userRepo repo.IUser
}

type User interface {
	Login(username, password string) (*models.User, error)
	ListUsers(nextId, limit int) ([]models.User, error)
	Create(*models.User, *models.UserQuota, int) error
	Delete(userId int, currentUserId int) error
}

func NewUser(userRepo repo.IUser) User {
	return &userEntity{
		userRepo: userRepo,
	}
}

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]

	}
	return string(b)
}

func (u userEntity) Delete(userId int, currentUserId int) error {
	user, err := u.userRepo.GetById(userId)
	if err != nil {
		return uer.InternalError(err)
	}

	if user == nil {
		// return nil, uer.NotFoundError(errors.New("user not found"))
		return uer.NotFoundError(errors.New("User not found"))
	}

	err = u.userRepo.DeleteUserAndQuota(userId)
	if err != nil {
		return uer.InternalError(err)
	}

	return nil
}

func (u userEntity) Create(user *models.User, quota *models.UserQuota, currentUserId int) error {
	randomSalt := RandStringBytesRmndr(15)

	hash, err := bcrypt.GenerateFromPassword([]byte(randomSalt+user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uer.InternalError(err)
	}

	user.Salt = randomSalt
	user.Password = string(hash)
	user.CreatedBy = currentUserId

	err = u.userRepo.CreateUserAndQuota(user, quota)
	if err != nil {
		return uer.InternalError(err)
	}

	return nil
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
