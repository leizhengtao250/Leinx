package lnet

import (
	"fmt"
	"leiTCP/Leinx/liface"
	"leiTCP/Leinx/utils"
	"strconv"
)

/**
消息处理模块的实现
**/

type MsgHandle struct {
	//存放每个msgID 对应的处理逻辑
	Apis map[uint32]liface.IRouter
	//负责Worker读取任务的消息队列
	TaskQueue []chan liface.IRequest
	//业务工作Worker池的数量
	WorkPoolSize uint32
}

//调度执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandle(request liface.IRequest) {
	//1.从request中找到msgID
	msgId := request.GetMsgID()
	//2.根据MsgID调度对应的router业务即可
	handle, ok := mh.Apis[msgId]
	if !ok {
		fmt.Println("api msgId=", request.GetMsgID(), "is NOT FOUND! NEED ADD ROUTER")
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)

}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router liface.IRouter) {
	//1.判断当前msg绑定的api处理方法是否存在
	if _, ok := mh.Apis[msgId]; ok {
		//id 已经注册
		panic("repeat api,msgID=" + strconv.Itoa(int(msgId)))
	}
	//2.添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId=", msgId, "success")
	//

}

/**
初始化/创建MsgHandle方法
*/

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]liface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkPoolSize, //从全局配置中获取
		TaskQueue:    make([]chan liface.IRequest, utils.GlobalObject.WorkPoolSize),
	}

}

//启动一个worker工作池(开启一个工作池只能发生一次，一个框架只能实现一个)
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerpoolsize 分别开启worker，每个worker用一个go承载
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		//一个worker开启
		//1 当前的worker对于的消息channel队列 开辟空间，就用第0个worker 就用第0个channel
		//缓冲通道
		mh.TaskQueue[i] = make(chan liface.IRequest, utils.GlobalObject.MaxWorkTaskNum)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerId int, taskQueue chan liface.IRequest) {
	fmt.Println("Worker ID=", workerId, "is starting")
	//不断阻塞等待对应的消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request
		case request := <-taskQueue:
			mh.DoMsgHandle(request)
		}
	}
}

//将消息交给TaskQueue，由worker来处理
func (mh *MsgHandle) SendMsgToTaskQueue(request liface.IRequest) {
	//1.将消息平均分配给不通过的worker（在分布式系统中，可以根据地域，ip等分配）
	//根据客户端连接的ConnId来进行分配（request分配）
	workerID := request.GetMsgID() % mh.WorkPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		"request MSGID=", request.GetMsgID(), "to WorkerID=", workerID)

	//2.将消息发送到对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request

}
