package handlers

import (
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/controller"
	"github.com/hippo-an/goranchise/ent/user"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type (
	Register struct {
		controller.Controller
		form RegisterForm
	}
	RegisterForm struct {
		Name            string `form:"name" validate:"required" label:"Name"`
		Email           string `form:"email" validate:"required,email" label:"Email"`
		Password        string `form:"password" validate:"required" label:"Password"`
		ConfirmPassword string `form:"confirm-password" validate:"required,eqfield=Password" label:"Confirm password"`
	}
)

func (r *Register) Get(c echo.Context) error {
	p := controller.NewPage(c)

	p.Layout = "auth"
	p.PageName = "register"
	p.Title = "Register"
	p.Data = RegisterForm{}
	if form := c.Get(context.FormKey); form != nil {
		p.Data = form.(RegisterForm)
	}

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

	exists, err := r.Container.ORM.User.
		Query().
		Where(user.Email(form.Email)).
		Exist(c.Request().Context())
	switch {
	case err != nil:
		return fail("unable to query to see if email is taken", err)
	case exists:
		msg.Warning(c, "A user with this email address already exists. Please log in.")
		return r.Redirect(c, "login")
	}

	c.Set(context.FormKey, *form)

	if err := c.Validate(form); err != nil {
		r.SetValidationErrorMessages(c, err, form)
		return r.Get(c)
	}

	hashedPassword, err := r.Container.Auth.HashPassword(form.Password)
	if err != nil {
		return fail("unable to hash password", err)
	}

	u, err := r.Container.ORM.User.
		Create().
		SetName(form.Name).
		SetEmail(form.Email).
		SetPassword(hashedPassword).
		Save(c.Request().Context())

	if err != nil {
		c.Logger().Error(err)
		msg.Danger(c, "Check the email and password")
		return r.Get(c)
	}
	c.Logger().Infof("user created: %s", u.Email)

	err = r.Container.Auth.Login(c, u.ID)
	if err != nil {
		c.Logger().Errorf("unable to log in: %v", err)
		msg.Info(c, "Your account has been created.")
		return r.Redirect(c, "login")
	}

	msg.Info(c, "Your account has been created. You are now logged in.")

	return r.Redirect(c, "home")
}
