package middlewares

import (
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"irisProject/utils"
)

func CheckUserRole(ctx iris.Context){
	_, role := GetJwtParams(ctx)
	if role != "user" {
		err := fmt.Errorf("role '%s' does not have sufficient permissions", role)
		utils.SetResponseError(ctx, iris.StatusForbidden, err.Error(), errors.New("RoleError"))
		ctx.StopExecution()
	}
	ctx.Next()
	return
}

func CheckAdminRole(ctx iris.Context){
	_, role := GetJwtParams(ctx)
	if role != "admin" {
		err := fmt.Errorf("role '%s' does not have sufficient permissions", role)
		utils.SetResponseError(ctx, iris.StatusForbidden, err.Error(), errors.New("RoleError"))
		ctx.StopExecution()
	}
	ctx.Next()
	return

}
