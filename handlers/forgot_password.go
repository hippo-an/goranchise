package handlers

import (
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/controller"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type (
	ForgotPassword struct {
		controller.Controller
		form ForgotPasswordForm
	}
	ForgotPasswordForm struct {
		Email string `form:"email" validate:"required,email" label:"Email"`
	}
)

func (f *ForgotPassword) Get(c echo.Context) error {
	p := controller.NewPage(c)
	p.Layout = "auth"
	p.PageName = "forgot-password"
	p.Title = "Forgot Password"
	p.Data = ForgotPasswordForm{}

	if form := c.Get(context.FormKey); form != nil {
		p.Data = form.(ForgotPasswordForm)
	}
	return f.RenderPage(c, p)
}

func (f *ForgotPassword) Post(c echo.Context) error {
	fail := func(message string, err error) error {
		c.Logger().Errorf("%s: %v", message, err)
		msg.Danger(c, "An error occurred. Please try again.")
		return f.Get(c)
	}

	form := new(ForgotPasswordForm)
	if err := c.Bind(form); err != nil {
		return fail("unable to parse forgot password form", err)
	}

	c.Set(context.FormKey, *form)

	if err := c.Validate(form); err != nil {
		f.SetValidationErrorMessages(c, err, form)
		return f.Get(c)
	}
	return f.Redirect(c, "home")
}
