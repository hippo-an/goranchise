package handlers

import (
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/controller"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/ent/user"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type (
	ResetPassword struct {
		controller.Controller
	}

	ResetPasswordForm struct {
		Password        string `form:"password" validate:"required" label:"Password"`
		ConfirmPassword string `form:"confirm-password" validate:"required,eqfield=Password" label:"Confirm Password"`
	}
)

func (r *ResetPassword) Get(c echo.Context) error {
	p := controller.NewPage(c)
	p.Title = "Reset Password"
	p.Layout = "auth"
	p.PageName = "reset-password"
	return r.RenderPage(c, p)
}

func (r *ResetPassword) Post(c echo.Context) error {
	fail := func(message string, err error) error {
		c.Logger().Errorf("%s: %v", message, err)
		msg.Danger(c, "An error occurred. Please try again.")
		return r.Get(c)
	}

	form := new(ResetPasswordForm)
	if err := c.Bind(form); err != nil {
		return fail("unable to parse reset password form", err)
	}

	if err := c.Validate(form); err != nil {
		r.SetValidationErrorMessages(c, err, form)
		return r.Get(c)
	}

	hash, err := r.Container.Auth.HashPassword(form.Password)
	if err != nil {
		return fail("unable to hash password", err)
	}

	usr := c.Get(context.UserKey).(*ent.User)

	_, err = r.Container.ORM.User.
		Update().
		SetPassword(hash).
		Where(user.ID(usr.ID)).
		Save(c.Request().Context())

	if err != nil {
		return fail("unable to update password", err)
	}

	err = r.Container.Auth.DeletePasswordTokens(c, usr.ID)
	if err != nil {
		return fail("unable to delete password tokens", err)
	}

	msg.Success(c, "Your password has been updated.")

	return r.Redirect(c, "login")
}
