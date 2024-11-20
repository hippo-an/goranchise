package services

import (
	"fmt"
	"github.com/hippo-an/goranchise/config"
	"github.com/labstack/echo/v4"
)

type MailClient struct {
	config    *config.Config
	templates *TemplateRenderer
}

func NewMailClient(config *config.Config, templates *TemplateRenderer) *MailClient {
	return &MailClient{config, templates}
}

func (c *MailClient) SendMail(ctx echo.Context, to, body string) error {
	if c.skipSend() {
		ctx.Logger().Debugf("skipping email sent to: %s", to)
	}
	ctx.Logger().Debug("Mock email sent. To: %s Body: %s", to, body)

	return nil
}

func (c *MailClient) SendTemplate(ctx echo.Context, to, template string, data interface{}) error {
	if c.skipSend() {
		ctx.Logger().Debugf("skipping template email sent to: %s", to)
	}

	if err := c.templates.Parse(
		"mail",
		template,
		template,
		[]string{fmt.Sprintf("mails/%s", template)},
		[]string{},
	); err != nil {
		return err
	}

	_, err := c.templates.Execute("mail", template, template, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *MailClient) skipSend() bool {
	return c.config.App.Environment != config.EnvironmentProd
}
