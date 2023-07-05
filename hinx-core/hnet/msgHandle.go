package hnet

import (
	"hinx/hinx-core/hconf"
	"hinx/hinx-core/hiface"
	"hinx/hinx-core/hlog"
	"strconv"
	"time"
)

type MsgHandle struct {
	Apis           map[uint32]hiface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan hiface.IRequest
}

// new msgHandle function
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{Apis: map[uint32]hiface.IRouter{},
		WorkerPoolSize: hconf.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan hiface.IRequest, hconf.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandle) DoMsgHandle(request hiface.IRequest) {
	// find msgId from request
	handle, ok := m.Apis[request.GetMsgID()]
	if !ok {
		hlog.Ins().InfoF("api msgID = %d is not found! Need register!", request.GetMsgID())
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

func (m *MsgHandle) AddRouter(id uint32, router hiface.IRouter) {
	//judge
	if _, ok := m.Apis[id]; ok {
		panic("repeat api, msgID= " + strconv.FormatUint(uint64(id), 10))
	}
	//add
	m.Apis[id] = router
	hlog.Ins().InfoF("Add api MsgID= %d success", id)
}

func (m *MsgHandle) StartWorkerPool() {
	//depend on workerPoolSize, start worker and go it
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan hiface.IRequest, hconf.GlobalObject.MaxWorkerTaskLen)
		go func(i int) { m.startOneWorker(i, m.TaskQueue[i]) }(i)
	}
}

func (m *MsgHandle) startOneWorker(id int, taskQueue chan hiface.IRequest) {
	hlog.Ins().InfoF("Worker ID= %d is started...", id)
	for {
		select {
		case request := <-taskQueue:
			{
				m.DoMsgHandle(request)
			}
		}
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request hiface.IRequest) {
	//workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	workerID := time.Now().UnixMicro() % int64(m.WorkerPoolSize)

	hlog.Ins().InfoF("Add ConnID = %d request MsgID = %d to workerID = %d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	m.TaskQueue[workerID] <- request
}
