package middleware

import (
	"fmt"
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/tests"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadUser(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	ctx.SetParamNames("userId")
	ctx.SetParamValues(fmt.Sprintf("%d", usr.ID))

	_ = tests.ExecuteMiddleware(ctx, LoadUser(c.ORM))
	ctxUsr, ok := ctx.Get(context.UserKey).(*ent.User)
	require.True(t, ok)
	require.Equal(t, usr.ID, ctxUsr.ID)
}
