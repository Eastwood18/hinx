package hnet

import (
	"fmt"
	hiface2 "hinx/hinx-core/hiface"
	"hinx/utils"
	"strconv"
	"time"
)

type MsgHandle struct {
	Apis           map[uint32]hiface2.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan hiface2.IRequest
}

// new msgHandle function
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{Apis: map[uint32]hiface2.IRouter{},
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan hiface2.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandle) DoMsgHandle(request hiface2.IRequest) {
	// find msgId from request
	handle, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is not found! Need register!")
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

func (m *MsgHandle) AddRouter(id uint32, router hiface2.IRouter) {
	//judge
	if _, ok := m.Apis[id]; ok {
		panic("repeat api, msgID= " + strconv.FormatUint(uint64(id), 10))
	}
	//add
	m.Apis[id] = router
	fmt.Println("Add api MsgID= ", id, " success")
}

func (m *MsgHandle) StartWorkerPool() {
	//depend on workerPoolSize, start worker and go it
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan hiface2.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go func(i int) { m.startOneWorker(i, m.TaskQueue[i]) }(i)
	}
}

func (m *MsgHandle) startOneWorker(id int, taskQueue chan hiface2.IRequest) {
	fmt.Println("Worker ID= ", id, " is started...")
	for {
		select {
		case request := <-taskQueue:
			{
				m.DoMsgHandle(request)
			}
		}
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request hiface2.IRequest) {
	//workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	workerID := time.Now().UnixMicro() % int64(m.WorkerPoolSize)

	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request MsgID = ", request.GetMsgID(), " to workerID = ", workerID)
	m.TaskQueue[workerID] <- request
}
