package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	userIdCtxKey       = "userId"
	userNameCtxKey     = "userName"
	groupIdCtxKey      = "groupId"
	groupAccountCtxKey = "groupAccount"
	orgIdCtxKey        = "orgId"
	orgNameCtxKey      = "orgName"
)

// UserId get user id from the context. REQUIRES Middleware to have run.
func UserId(ctx *gin.Context) (int64, error) {
	raw := ctx.GetInt64(userIdCtxKey)

	if raw == 0 {
		err := errors.New("userId does not exists in gin.Context")
		return 0, err
	}

	return raw, nil
}

// UserName get user name from the context. REQUIRES Middleware to have run.
func UserName(ctx *gin.Context) string {
	return ctx.GetString(userNameCtxKey)
}

// UserGroupId get user group id from the context. REQUIRES Middleware to have run.
func UserGroupId(ctx *gin.Context) (int64, error) {
	raw := ctx.GetInt64(groupIdCtxKey)

	if raw == 0 {
		err := errors.New("user group id does not exists in gin.Context")
		return 0, err
	}

	return raw, nil
}

// UserGroupName get user group name from the context. REQUIRES Middleware to have run.
func UserGroupName(ctx *gin.Context) string {
	return ctx.GetString(groupAccountCtxKey)
}

// OrgId get organization id from the context. REQUIRES Middleware to have run.
func OrgId(ctx *gin.Context) (int64, error) {
	raw := ctx.GetInt64(orgIdCtxKey)

    /*
	if raw == 0 {
		err := errors.New("orgId does not exists in gin.Context")
		return 0, err
	}
    */

	return raw, nil
}

// OrgName get org name from the context. REQUIRES Middleware to have run.
func OrgName(ctx *gin.Context) string {
	return ctx.GetString(orgNameCtxKey)
}
