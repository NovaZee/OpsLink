package router

// Handler 接口定义了处理程序的方法
type Handler interface {
	Register(router *Router)
}
