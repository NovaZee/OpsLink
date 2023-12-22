package router

import (
	"context"
	"github.com/denovo/permission/pkg/service/casbin"
	"github.com/denovo/permission/pkg/service/store"
	"github.com/denovo/permission/pkg/util"
	"github.com/denovo/permission/protoc/model"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type PublicHandler struct {
	casbin *casbin.Casbin
	ss     store.StoreService

	middlewares []gin.HandlerFunc
}

func BuildPublic(cb *casbin.Casbin, storeService store.StoreService, middleware ...gin.HandlerFunc) *PublicHandler {
	return &PublicHandler{casbin: cb, ss: storeService, middlewares: middleware}
}

// LogIn 登录
func (rh *PublicHandler) LogIn(ctx *gin.Context) {
	var font model.Role
	var ctx2 context.Context
	if err := ctx.ShouldBind(&font); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": ErrorParamsError, "status": http.StatusBadRequest})
		return
	}
	if font.Name == "" || font.Password == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "参数不能为空")
		return
	}
	r2, err2 := rh.ss.Get(ctx2, font.Name)

	if err2 != nil || r2 == nil {
		ErrorResponse(ctx, http.StatusBadRequest, font.Name+" 不存在")
		return
	}
	if font.Name != r2.Name || font.Password != r2.Password {
		ErrorResponse(ctx, http.StatusBadRequest, r2.Name+" 账号密码不匹配!")
		return
	}
	token, tokenErr := util.GenerateToken(r2.Id, r2.Name)
	if tokenErr != nil {
		ErrorResponse(ctx, http.StatusBadRequest, r2.Name+" 系统错误!")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": token, "status": http.StatusOK})
	return
}

// SignIn 注册
func (rh *PublicHandler) SignIn(ctx *gin.Context) {
	var font model.Role
	var ctx2 context.Context
	if err := ctx.ShouldBind(&font); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, ErrorParamsError)
		return
	}
	if font.Name == "" || font.Password == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "参数不能为空")
		return
	}
	r2, err2 := rh.ss.Get(ctx2, font.Name)
	if err2 != nil || r2 != nil {
		ErrorResponse(ctx, http.StatusBadRequest, r2.Name+" 已存在")
		return
	}

	newRole := &model.Role{
		Name:     font.Name,
		Password: font.Password,
	}
	rand.Seed(time.Now().UnixNano())
	newRole.Id = rand.Int63()
	token, err := util.GenerateToken(newRole.Id, newRole.Name)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "token 生成失败")
		return
	}
	// 成员信息存入
	e := rh.ss.Create(ctx2, newRole)
	if e != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "用户生成失败")
		return
	}
	//// TODO:成员权限初始化
	//_ = rh.casbin.AddGroupingPolicy(newRole.Name, casbin.GroupRead)
	ctx.JSON(http.StatusOK, gin.H{"message": token, "status": http.StatusOK})
	return
}

func (rh *PublicHandler) GetName() string {
	return "public"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (rh *PublicHandler) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	r := g.Use(middle...)
	{
		r.POST("/logIn", func(ctx *gin.Context) { rh.LogIn(ctx) })
		r.POST("/signIn", func(ctx *gin.Context) { rh.SignIn(ctx) })
	}
}
