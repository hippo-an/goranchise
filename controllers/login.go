package controllers

import (
	"github.com/hippo-an/goranchise/ent/user"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Controller
}

func (l *Login) Get(c echo.Context) error {
	p := NewPage(c)

	p.Layout = "auth"
	p.PageName = "login"
	p.Title = "Login"
	p.Data = "Login Page"

	return l.RenderPage(c, p)
}

func (l *Login) Post(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		msg.Warning(c, "All fields are required to login")
		return l.Get(c)
	}

	u, err := l.Container.ORM.User.
		Query().
		Where(user.Username(username)).
		First(c.Request().Context())

	if err != nil {
		c.Logger().Errorf("error querying user during login: %v", err)
		msg.Danger(c, "Check username and password")
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			msg.Danger(c, "Invalid credentials. Please try again.")
		}
	}

	msg.Info(c, "You are now logged in.")

	return l.Redirect(c, "home")
}
