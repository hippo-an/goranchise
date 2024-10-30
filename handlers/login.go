package handlers

import (
	"fmt"
	"github.com/hippo-an/goranchise/auth"
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/controller"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/ent/user"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type (
	Login struct {
		controller.Controller
		form LoginForm
	}
	LoginForm struct {
		Email    string `form:"email" validate:"required,email" label:"Email"`
		Password string `form:"password" validate:"required" label:"Password"`
	}
)

func (l *Login) Get(c echo.Context) error {
	p := controller.NewPage(c)

	p.Layout = "auth"
	p.PageName = "login"
	p.Title = "Login"
	p.Data = LoginForm{}
	if form := c.Get(context.FormKey); form != nil {
		p.Data = form.(LoginForm)
	}
	return l.RenderPage(c, p)
}

func (l *Login) Post(c echo.Context) error {
	fail := func(message string, err error) error {
		c.Logger().Errorf("%s: %v", message, err)
		msg.Danger(c, "An error occurred. Please try again.")
		return l.Get(c)
	}

	form := new(LoginForm)
	if err := c.Bind(form); err != nil {
		return fail("unable to parse login form", err)
	}
	c.Set(context.FormKey, *form)

	if err := c.Validate(form); err != nil {
		l.SetValidationErrorMessages(c, err, form)
		return l.Get(c)
	}

	u, err := l.Container.ORM.User.
		Query().
		Where(user.Email(form.Email)).
		First(c.Request().Context())

	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			msg.Danger(c, "Check email and password")
			return l.Get(c)
		default:
			return fail("error querying user during login", err)
		}
	}

	err = auth.CheckPassword(form.Password, u.Password)
	if err != nil {
		msg.Danger(c, "Invalid credentials. Please try again.")
		return l.Get(c)
	}

	err = auth.Login(c, u.ID)
	if err != nil {
		return fail("unable to log in user", err)
	}

	msg.Success(c, fmt.Sprintf("Welcome back, %s. You are now logged in.", u.Email))
	return l.Redirect(c, "home")
}
