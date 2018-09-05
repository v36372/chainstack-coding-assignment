package handler

import (
	"chainstack/app/entity"
	"chainstack/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockUser struct{}

func (mockUser) GetByEmail(email string) (*models.User, error) {
	return &models.User{
		Id:       1,
		Email:    "tin.trng.ng@gmail.com",
		Password: "$2a$10$crGQ4a8.L99zl6Dku14SdO0m3RbHUfaFHOuRPjv5CtvPifK4bQUpG",
		Salt:     "7nWZLcCK0vsPzIM",
		IsAdmin:  true,
	}, nil
}

func (mockUser) Create(*models.User) error {
	return nil
}

func (mockUser) Update(*models.User) error {
	return nil
}

func (mockUser) Delete(*models.User) error {
	return nil
}

func getLoginPayload(email, password string) string {
	data := url.Values{}

	data.Add("email", email)
	data.Add("password", password)

	return data.Encode()
}

func TestLoginHandler(t *testing.T) {
	loginPayload := getLoginPayload("tin.trng.ng@gmail.com", "test")
	req, err := http.NewRequest("POST", "/login", strings.NewReader(loginPayload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	r := gin.New()

	mockUser := mockUser{}
	userHandler := userHandler{
		user: entity.NewUser(mockUser),
	}
	r.POST("/login", userHandler.Login)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fail()
	}
}
