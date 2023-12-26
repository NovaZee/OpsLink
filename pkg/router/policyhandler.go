package router

import (
	"github.com/denovo/permission/pkg/router/kubehandler"
	"github.com/denovo/permission/pkg/service/casbin"
	"github.com/denovo/permission/pkg/service/store"
	"github.com/denovo/permission/pkg/util"
	"github.com/denovo/permission/protoc/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type PolicyHandler struct {
	cb          *casbin.Casbin
	middlewares []gin.HandlerFunc

	ss store.StoreService
}

func BuildRolePolicy(cb *casbin.Casbin, ss store.StoreService, middleware ...gin.HandlerFunc) *PolicyHandler {
	return &PolicyHandler{
		cb:          cb,
		middlewares: middleware,
		ss:          ss,
	}
}

// AddPolicy g
func (ph *PolicyHandler) AddPolicy(ctx *gin.Context) {
	var pModel *model.PModel
	err := ctx.ShouldBindJSON(&pModel)
	if err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, errors.New("实体错误"))
		return
	}
	pModel.PType = "p"
	add, err := ph.cb.AddPolicy(pModel)
	if add == false || err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	kubehandler.KubeSuccessResponse(ctx, http.StatusOK)
	return
}

// bind g
func (ph *PolicyHandler) bind(ctx *gin.Context) {
	var gModel *model.GModel
	err := ctx.ShouldBindJSON(&gModel)
	if err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, errors.New("实体错误"))
		return
	}
	gModel.PType = "g"
	add, err := ph.cb.BindingRoles(gModel)
	if add == false || err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	kubehandler.KubeSuccessResponse(ctx, http.StatusOK)
	return
}

// unBind g
func (ph *PolicyHandler) unBind(ctx *gin.Context) {
	var gModel *model.GModel
	err := ctx.ShouldBindJSON(&gModel)
	if err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, errors.New("实体错误"))
		return
	}
	gModel.PType = "g"
	add, err := ph.cb.UnBindingRoles(gModel)
	if add == false || err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	kubehandler.KubeSuccessResponse(ctx, http.StatusOK)
	return
}

// DeletePolicy  删除权限策略 -manager
func (ph *PolicyHandler) DeletePolicy(ctx *gin.Context) {
	var pModel *model.PModel
	err := ctx.ShouldBindJSON(&pModel)
	if err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	pModel.PType = "p"
	add, err := ph.cb.Delete(pModel)
	if add == false || err != nil {
		kubehandler.KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	kubehandler.KubeSuccessResponse(ctx, http.StatusOK)
	return
}

// UpdatePolicy  删除权限策略 -manager
func (ph *PolicyHandler) UpdatePolicy(ctx *gin.Context) {
	var req *model.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.NewPolicy.PType = "p"
	req.OldPolicy.PType = "p"
	update, err := ph.cb.Update(req)
	if update == false || err != nil {
		ctx.JSONP(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "更新成功", "status": http.StatusOK})
	return
}

func (ph *PolicyHandler) ListALlRoles(ctx *gin.Context) {
	pageNo := ctx.DefaultQuery("pageNo", "1")
	PageSize := ctx.DefaultQuery("PageSize", "10")
	num, err := strconv.Atoi(pageNo)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid parameter: id must be an integer"})
		return
	}
	if num < 1 {
		num = 1
	}
	size, err := strconv.Atoi(PageSize)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid parameter: id must be an integer"})
		return
	}
	role := ph.ss.GetRole()
	length := len(role.Roles)
	if length != 0 {
		start, end := util.PageSlice(length, num, size)
		kubehandler.KubeSuccessMsgResponse(ctx, http.StatusOK, role.Roles[start:end])
		return
	}
}

// SignIn 注册
func (ph *PolicyHandler) SignIn(ctx *gin.Context) {
	var font *model.Role
	if err := ctx.ShouldBindJSON(&font); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, ErrorParamsError)
		return
	}
	if font.Name == "" || font.Password == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "参数不能为空")
		return
	}
	r2, err2 := ph.ss.Get(ctx, font.Name)
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
	e := ph.ss.Create(ctx, newRole)
	if e != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "用户生成失败")
		return
	}
	//// TODO:成员权限初始化
	//_ = rh.casbin.AddGroupingPolicy(newRole.Name, casbin.GroupRead)
	ctx.JSON(http.StatusOK, gin.H{"message": token, "status": http.StatusOK})
	return
}

func (ph *PolicyHandler) ListUses(ctx *gin.Context) {

}
func (ph *PolicyHandler) ListRoles(ctx *gin.Context) {
	role := ctx.Param("role")
	kubehandler.KubeSuccessMsgResponse(ctx, http.StatusOK, ph.cb.ListRoles(role))
	return
}

func (ph *PolicyHandler) ListMyPolicy(ctx *gin.Context) {
	uname := ctx.Param("uname")
	res := ph.cb.ListMyPolicy(uname)
	kubehandler.KubeSuccessMsgResponse(ctx, http.StatusOK, res)
	return
}

func (ph *PolicyHandler) GetName() string {
	return "policy"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	read := g.Use(middle...)
	{
		read.GET("/:uname", func(ctx *gin.Context) { ph.ListMyPolicy(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	admin := g.Use(middle...)
	{
		// policies
		admin.POST("bind", func(ctx *gin.Context) { ph.bind(ctx) })
		admin.POST("unbind", func(ctx *gin.Context) { ph.unBind(ctx) })
		admin.POST("add", func(ctx *gin.Context) { ph.AddPolicy(ctx) })

		admin.DELETE("delete", func(ctx *gin.Context) { ph.DeletePolicy(ctx) })
		admin.POST("update", func(ctx *gin.Context) { ph.UpdatePolicy(ctx) })

		admin.GET("users", func(ctx *gin.Context) { ph.ListUses(ctx) })
		admin.GET("roles/:role", func(ctx *gin.Context) { ph.ListRoles(ctx) })

		// users
		admin.GET("roles/members", func(ctx *gin.Context) { ph.ListALlRoles(ctx) })
		admin.POST("roles/signIn", func(ctx *gin.Context) { ph.SignIn(ctx) })
	}
}
