package auth

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/ent/passwordtoken"
	"github.com/hippo-an/goranchise/ent/user"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	sessionName             = "ua"
	sessionKeyUserId        = "user_id"
	sessionKeyAuthenticated = "authenticated"
	passwordTokenLength     = 64
)

type InvalidTokenError struct{}

func (e InvalidTokenError) Error() string {
	return "invalid token"
}

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
			Only(ctx.Request().Context())
	}
	return nil, NotAuthenticatedError{}
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
	token, err := c.RandomToken(passwordTokenLength)

	if err != nil {
		return "", nil, err
	}

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

func (c *Client) RandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (c *Client) GetValidPasswordToken(ctx echo.Context, token string, userId int) (*ent.PasswordToken, error) {

	expiration := time.Now().Add(-c.config.App.PasswordTokenExpiration)
	pts, err := c.orm.PasswordToken.
		Query().
		Where(passwordtoken.HasUserWith(user.ID(userId))).
		Where(passwordtoken.CreatedAtGTE(expiration)).
		All(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Error(err)
		return nil, err
	}

	for _, pt := range pts {
		if err := c.CheckPassword(token, pt.Hash); err == nil {
			return pt, nil
		}
	}

	return nil, InvalidTokenError{}
}

func (c *Client) DeletePasswordTokens(ctx echo.Context, userId int) error {
	_, err := c.orm.PasswordToken.
		Delete().
		Where(passwordtoken.HasUserWith(user.ID(userId))).
		Exec(ctx.Request().Context())
	return err
}
