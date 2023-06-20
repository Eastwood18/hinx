package hiface

type IRouter interface {
	// PreHandle hook function before conn business
	PreHandle(request IRequest)
	// Handle main hook function of conn business
	Handle(request IRequest)
	// PostHandle hook function after conn business
	PostHandle(request IRequest)
}
