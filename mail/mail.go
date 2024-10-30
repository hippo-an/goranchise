package mail

import (
	"github.com/hippo-an/goranchise/config"
	"github.com/labstack/echo/v4"
)

type Client struct {
	config *config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{config}
}

func (c *Client) SendMail(ctx echo.Context, to, body string) error {
	if c.config.App.Environment != config.EnvironmentProd {

	}
	ctx.Logger().Debug("Mock email sent. To: %s Body: %s", to, body)

	return nil
}
