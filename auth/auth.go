package auth

import (
	"errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionName             = "ua"
	sessionKeyUserId        = "user_id"
	sessionKeyAuthenticated = "authenticated"
)

func Login(c echo.Context, userId int) error {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return err
	}

	sess.Values[sessionKeyUserId] = userId
	sess.Values[sessionKeyAuthenticated] = true
	return sess.Save(c.Request(), c.Response())
}

func Logout(c echo.Context) error {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return err
	}
	sess.Values[sessionKeyAuthenticated] = false
	return sess.Save(c.Request(), c.Response())
}

func GetUserID(c echo.Context) (int, error) {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return 0, err
	}
	if sess.Values[sessionKeyAuthenticated] == true {
		return sess.Values[sessionKeyUserId].(int), nil
	}
	return 0, errors.New("user not authenticated")
}
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
