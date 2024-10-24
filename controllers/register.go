package controllers

import (
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	Register struct {
		Controller
		form RegisterForm
	}
	RegisterForm struct {
		Username string `form:"username" validate:"required"`
		Password string `form:"password" validate:"required"`
	}
)

func (r *Register) Get(c echo.Context) error {
	p := NewPage(c)

	p.Layout = "auth"
	p.PageName = "register"
	p.Title = "Register"
	p.Data = r.form

	return r.RenderPage(c, p)
}

func (r *Register) Post(c echo.Context) error {
	fail := func(message string, err error) error {
		c.Logger().Errorf("%s: %v", message, err)
		msg.Danger(c, "An error occurred. Please try again.")
		return r.Get(c)
	}

	form := new(RegisterForm)
	if err := c.Bind(form); err != nil {
		return fail("unable to parse form values", err)
	}
	r.form = *form

	if err := c.Validate(form); err != nil {
		r.Container.Web.Logger.Errorf("Validation error: %s", err)
		msg.Danger(c, "All fields are required.")
		return r.Get(c)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return fail("unable to hash password", err)
	}

	u, err := r.Container.ORM.User.
		Create().
		SetUsername(form.Username).
		SetPassword(string(hashedPassword)).
		Save(c.Request().Context())

	if err != nil {
		c.Logger().Error(err)
		msg.Danger(c, "Check the username and password")
		return r.Get(c)
	} else {
		c.Logger().Infof("user created: %s", u.Username)
	}

	c.Logger().Infof("user created: %s", u.Username)
	msg.Info(c, "Your account has been created. You are now logged in.")

	return r.Redirect(c, "home")
}
