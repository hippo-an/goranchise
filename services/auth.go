package services

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
	authSessionName             = "ua"
	authSessionKeyUserId        = "user_id"
	authSessionKeyAuthenticated = "authenticated"
)

type AuthClient struct {
	config *config.Config
	orm    *ent.Client
}

func NewAuthClient(cfg *config.Config, orm *ent.Client) *AuthClient {
	return &AuthClient{
		config: cfg,
		orm:    orm,
	}
}

func (a *AuthClient) Login(ctx echo.Context, userId int) error {
	sess, err := session.Get(authSessionName, ctx)
	if err != nil {
		return err
	}

	sess.Values[authSessionKeyUserId] = userId
	sess.Values[authSessionKeyAuthenticated] = true
	return sess.Save(ctx.Request(), ctx.Response())
}

func (a *AuthClient) Logout(ctx echo.Context) error {
	sess, err := session.Get(authSessionName, ctx)
	if err != nil {
		return err
	}
	sess.Values[authSessionKeyAuthenticated] = false
	return sess.Save(ctx.Request(), ctx.Response())
}

func (a *AuthClient) GetAuthenticatedUserId(ctx echo.Context) (int, error) {
	sess, err := session.Get(authSessionName, ctx)
	if err != nil {
		return 0, err
	}
	if sess.Values[authSessionKeyAuthenticated] == true {
		return sess.Values[authSessionKeyUserId].(int), nil
	}
	return 0, NotAuthenticatedError{}
}

func (a *AuthClient) GetAuthenticatedUser(ctx echo.Context) (*ent.User, error) {
	if userId, err := a.GetAuthenticatedUserId(ctx); err == nil {
		return a.orm.User.Query().
			Where(user.ID(userId)).
			Only(ctx.Request().Context())
	}
	return nil, NotAuthenticatedError{}
}

func (a *AuthClient) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (a *AuthClient) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *AuthClient) GeneratePasswordResetToken(ctx echo.Context, userId int) (string, *ent.PasswordToken, error) {
	token, err := a.RandomToken(a.config.App.PasswordToken.Length)

	if err != nil {
		return "", nil, err
	}

	hash, err := a.HashPassword(token)
	if err != nil {
		return "", nil, err
	}

	pt, err := a.orm.PasswordToken.
		Create().
		SetHash(hash).
		SetUserID(userId).
		Save(ctx.Request().Context())

	return token, pt, err
}

func (a *AuthClient) RandomToken(length int) (string, error) {
	b := make([]byte, (length/2)+1)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	return token[:length], nil
}

func (a *AuthClient) GetValidPasswordToken(ctx echo.Context, token string, userId int) (*ent.PasswordToken, error) {

	expiration := time.Now().Add(-a.config.App.PasswordToken.Expiration)
	pts, err := a.orm.PasswordToken.
		Query().
		Where(passwordtoken.HasUserWith(user.ID(userId))).
		Where(passwordtoken.CreatedAtGTE(expiration)).
		All(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Error(err)
		return nil, err
	}

	for _, pt := range pts {
		if err := a.CheckPassword(token, pt.Hash); err == nil {
			return pt, nil
		}
	}

	return nil, InvalidPasswordTokenError{}
}

func (a *AuthClient) DeletePasswordTokens(ctx echo.Context, userId int) error {
	_, err := a.orm.PasswordToken.
		Delete().
		Where(passwordtoken.HasUserWith(user.ID(userId))).
		Exec(ctx.Request().Context())
	return err
}
