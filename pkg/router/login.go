package router

import (
	"github.com/denovo/permission/pkg/service/casbin"
	"github.com/denovo/permission/pkg/service/store"
	"github.com/denovo/permission/pkg/util"
	"github.com/denovo/permission/protoc/model"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var font *model.Role
	if err := ctx.ShouldBindJSON(&font); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": ErrorParamsError, "status": http.StatusBadRequest})
		return
	}
	if font.Name == "" || font.Password == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "参数不能为空")
		return
	}
	r2, err2 := rh.ss.Get(ctx, font.Name)

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

func (rh *PublicHandler) GetName() string {
	return "public"
}

// ReadRegister ReadRegister
func (rh *PublicHandler) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	r := g.Use(middle...)
	{
		r.POST("/logIn", func(ctx *gin.Context) { rh.LogIn(ctx) })
	}
}
