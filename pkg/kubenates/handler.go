package kubenates

// ActiveHandler 接口
type ActiveHandler interface {
	AddHandler
	UpdateHandler
	DeleteHandler
}

type AddHandler interface {
	HandleAddEvent(data string)
}

type UpdateHandler interface {
	HandleUpdateEvent(data string)
}

type DeleteHandler interface {
	HandleDeleteEvent(data string)
}
