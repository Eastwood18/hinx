package hiface

type IMsgHandle interface {
	DoMsgHandle(request IRequest)
	AddRouter(id uint32, router IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}
