package handlers

import (
	"github.com/hippo-an/goranchise/auth"
	"github.com/hippo-an/goranchise/controller"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type (
	Register struct {
		controller.Controller
		form RegisterForm
	}
	RegisterForm struct {
		Email    string `form:"email" validate:"required,email" label:"Email"`
		Password string `form:"password" validate:"required" label:"Password"`
	}
)

func (r *Register) Get(c echo.Context) error {
	p := controller.NewPage(c)

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

	if err := c.Bind(&r.form); err != nil {
		return fail("unable to parse form values", err)
	}

	if err := c.Validate(r.form); err != nil {
		r.Container.Web.Logger.Errorf("Validation error: %s", err)
		msg.Danger(c, "All fields are required.")
		return r.Get(c)
	}

	hashedPassword, err := auth.HashPassword(r.form.Password)
	if err != nil {
		return fail("unable to hash password", err)
	}

	u, err := r.Container.ORM.User.
		Create().
		SetEmail(r.form.Email).
		SetPassword(hashedPassword).
		Save(c.Request().Context())

	if err != nil {
		c.Logger().Error(err)
		msg.Danger(c, "Check the email and password")
		return r.Get(c)
	} else {
		c.Logger().Infof("user created: %s", u.Email)
	}

	c.Logger().Infof("user created: %s", u.Email)

	err = auth.Login(c, u.ID)
	if err != nil {
		// TODO
	}

	msg.Info(c, "Your account has been created. You are now logged in.")

	return r.Redirect(c, "home")
}
