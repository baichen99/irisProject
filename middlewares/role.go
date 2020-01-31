package middlewares

import (
	"errors"
	"fmt"
	"irisProject/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type roleConfig struct {
	Role string
}

// RoleMiddleware middleware used to check the JWT's role
type RoleMiddleware struct {
	roleConfig roleConfig
}

// NewRoleMiddleware returns a new user role middleware with given config
func NewRoleMiddleware(cfg ...roleConfig) *RoleMiddleware {
	var c roleConfig
	if len(cfg) == 0 {
		c = roleConfig{}
	} else {
		c = cfg[0]
	}
	return &RoleMiddleware{roleConfig: c}
}

// CheckUserRole will check the user role
func (m *RoleMiddleware) CheckUserRole(ctx context.Context) (err error) {
	_, role := GetJWTParams(ctx)
	if role != "super" && role != m.roleConfig.Role {
		err = fmt.Errorf("role '%s' does not have sufficient permissions", role)
		utils.SetResponseError(ctx, iris.StatusForbidden, err.Error(), errors.New("RoleError"))
		return
	}
	return nil
}

// Serve is http handler
func (m *RoleMiddleware) Serve(ctx context.Context) {
	if err := m.CheckUserRole(ctx); err != nil {
		ctx.StopExecution()
		return
	}
	ctx.Next()
}
