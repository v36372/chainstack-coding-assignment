package middleware

import (
	"fmt"

	"net/http"

	"github.com/gorilla/securecookie"
)

type SecCookie struct {
	secureCookie *securecookie.SecureCookie
}

func NewSetCookie(blockKeyStr, hashKeyStr string) *SecCookie {
	hashKey := []byte(hashKeyStr)
	blockKey := []byte(blockKeyStr)
	sec := securecookie.New(hashKey, blockKey)
	sec.MaxAge(3600 * 24 * 365 * 3)
	return &SecCookie{sec}
}

func (s *SecCookie) GetCurrentUserID(request *http.Request) (string, error) {
	token := request.Header.Get("X-Access-Token")
	if token != "" {
		value := ""
		err := s.secureCookie.Decode("auth", token, &value)
		return value, err
	}
	return s.DecodeToken("auth", request)
}

func (s *SecCookie) GetUserIDFromToken(token string) (string, error) {
	var (
		userIDStr string
		err       error
	)
	err = s.secureCookie.Decode("auth", token, &userIDStr)
	return userIDStr, err
}

func (s *SecCookie) SetAuthorizationToken(cookieName, value, path string, w http.ResponseWriter) (string, error) {
	s.secureCookie.MaxAge(3600 * 24 * 365 * 3)
	if encoded, err := s.secureCookie.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			Path:     path,
			MaxAge:   3600 * 24 * 365 * 3,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		return encoded, nil
	} else {
		return "", err
	}
}

func (s *SecCookie) DecodeToken(cookieName string, r *http.Request) (string, error) {
	var cookie *http.Cookie
	var err error
	var value string
	if cookie, err = r.Cookie(cookieName); err != nil {
		return "", err
	}
	if err = s.secureCookie.Decode(cookieName, cookie.Value, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (s *SecCookie) ClearCookie(w http.ResponseWriter, cookieName, cookiePath string) {
	ignoredContent := "BLABLA" // random string
	cookie := fmt.Sprintf("%s=%s; path=%s; expires=Thu, 01 Jan 1970 00:00:00 GMT", cookieName, ignoredContent, cookiePath)
	SetHeader(w, "Set-Cookie", cookie, true)
}

func SetHeader(w http.ResponseWriter, hdr, val string, unique bool) {
	if unique {
		w.Header().Set(hdr, val)
	} else {
		w.Header().Add(hdr, val)
	}
}
