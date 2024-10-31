package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/ent/user"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionName             = "ua"
	sessionKeyUserId        = "user_id"
	sessionKeyAuthenticated = "authenticated"
)

type Client struct {
	config *config.Config
	orm    *ent.Client
}
type NotAuthenticatedError struct {
}

func (e NotAuthenticatedError) Error() string {
	return "user not authenticated"
}

func NewClient(cfg *config.Config, orm *ent.Client) *Client {
	return &Client{
		config: cfg,
		orm:    orm,
	}
}

func (c *Client) Login(ctx echo.Context, userId int) error {
	sess, err := session.Get(sessionName, ctx)
	if err != nil {
		return err
	}

	sess.Values[sessionKeyUserId] = userId
	sess.Values[sessionKeyAuthenticated] = true
	return sess.Save(ctx.Request(), ctx.Response())
}

func (c *Client) Logout(ctx echo.Context) error {
	sess, err := session.Get(sessionName, ctx)
	if err != nil {
		return err
	}
	sess.Values[sessionKeyAuthenticated] = false
	return sess.Save(ctx.Request(), ctx.Response())
}

func (c *Client) GetAuthenticatedUserId(ctx echo.Context) (int, error) {
	sess, err := session.Get(sessionName, ctx)
	if err != nil {
		return 0, err
	}
	if sess.Values[sessionKeyAuthenticated] == true {
		return sess.Values[sessionKeyUserId].(int), nil
	}
	return 0, NotAuthenticatedError{}
}

func (c *Client) GetAuthenticatedUser(ctx echo.Context) (*ent.User, error) {
	if userId, err := c.GetAuthenticatedUserId(ctx); err == nil {
		return c.orm.User.Query().
			Where(user.ID(userId)).
			First(ctx.Request().Context())
	}
	return nil, errors.New("user not authenticated")
}

func (c *Client) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (c *Client) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (c *Client) GeneratePasswordResetToken(ctx echo.Context, userId int) (string, *ent.PasswordToken, error) {
	token := c.RandomToken(64)

	hash, err := c.HashPassword(token)
	if err != nil {
		return "", nil, err
	}

	pt, err := c.orm.PasswordToken.
		Create().
		SetHash(hash).
		SetUserID(userId).
		Save(ctx.Request().Context())

	return token, pt, err
}

func (c *Client) RandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
