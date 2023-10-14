package router

import (
	"errors"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/service/role"
	"github.com/denovo/permission/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddPolicy 新增权限策略 -manager
func AddPolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.DefaultPolicy.Add(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "添加重复", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "添加成功", "status": http.StatusOK})
	return
}

// DeletePolicy  删除权限策略 -manager
func DeletePolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.DefaultPolicy.Delete(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}

// DeletePolicy  删除权限策略 -manager
func UpdatePolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.DefaultPolicy.Delete(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}

// 共享的请求参数处理逻辑 -manager
func processManagerRequestParams(ctx *gin.Context) (*casbin.CasbinModel, error) {
	role := ctx.Query("role")
	source := ctx.Query("source")
	behavior := ctx.Query("behavior")
	if len(role) == 0 || len(source) == 0 || len(behavior) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "params errors", "status": http.StatusBadRequest})
		return nil, errors.New("params errors")
	}
	casbinModel := casbin.NewCasbinModel(role, source, behavior)
	return casbinModel, nil
}

// LogIn 登录
func LogIn(ctx *gin.Context) {
	var font role.FrontRole
	if err := ctx.ShouldBind(&font); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": ErrorParamsError, "status": http.StatusBadRequest})
		return
	}

}

// SignIn 注册
func SignIn(ctx *gin.Context, r *Router) {
	var font role.FrontRole
	if err := ctx.ShouldBind(&font); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": ErrorParamsError, "status": http.StatusBadRequest})
		return
	}
	roles, err2 := r.PermissionEtcdClient.GetPermissionPolicy(font.Name)
	if err2 != nil || len(roles) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": font.Name + " 已存在", "status": http.StatusBadRequest})
		return
	}
	newRole := role.NewRole(font)
	token, err := util.GenerateToken(newRole.Id, newRole.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "token 生成失败", "status": http.StatusBadRequest})
		return
	}
	// 成员信息存入
	e := r.PermissionEtcdClient.SetPermissionPolicy(newRole)
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户 生成失败", "status": http.StatusBadRequest})
		return
	}
	// 成员权限初始化
	_ = r.cb.DefaultPolicy.AddGroupingPolicy(newRole.Name, casbin.GroupRead)
	ctx.JSON(http.StatusOK, gin.H{"message": token, "status": http.StatusOK})
	return

}
