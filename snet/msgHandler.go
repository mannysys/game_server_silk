package snet

import (
	"fmt"
	"game_server_silk/siface"
	"game_server_silk/utils"
	"strconv"
)

/*
	消息管理处理模块的实现
 */
type MsgHandler struct {
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32] siface.IRouter

	//业务工作Worker池的worker数量（表示开启多少个chan做为管道消息队列）
	WorkerPoolSize uint32
	//负责Worker 取任务的消息队列
	TaskQueue []chan siface.IRequest
}

//创建初始化MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32] siface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,//全局配置参数文件中获取
		TaskQueue:make([]chan siface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//调度执行对应的Router的消息处理方法
func (mh *MsgHandler) DoMsgHandler(request siface.IRequest) {
	//1 从Request找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("调度执行的路由API无法找到", request.GetMsgID(), ", 你需要注册该路由")
		return
	}

	//2 根据MsgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router siface.IRouter) {
	//1 判断 当前msg绑定的API处理方法是否已经存在（检查map集合中是否已经存在msgID对应的路由）
	if _, ok := mh.Apis[msgID]; ok {
		//msgID绑定的路由，已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	//2 添加msg与API绑定关系（添加新msgID绑定的路由处理业务）
	mh.Apis[msgID] = router
	fmt.Println("路由添加到map集合中成功！", msgID)

}


/*
	启动一个Worker工作池（开启工作池的方法调用只能发生一次，一个silk框架只能有一个worker工作池）
 */
func (mh *MsgHandler) StartWorkerPool() {
	//根据workerPoolSize 分别开启worker，每个worker用一个goruntine来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 当前worker对应的channel消息队列 开辟空间 第0个worker 就用第0个channel
		mh.TaskQueue[i] = make(chan siface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

/*
	启动一个Worker工作流程
 */
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan siface.IRequest) {
	fmt.Println("[启动工作池消息队列Chan] worker ID = ", workerID, "is started ...")

	//不断的阻塞等待对应消息队列的消息
	for {
		select {
			case request := <-taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}

/*
	将消息交给 TaskQueue，由worker进行处理
 */
func (mh *MsgHandler) SendMsgToTaskQueue(request siface.IRequest) {
	//1 将消息平均分配给不同的worker
	//根据客户端建立的ConnID来进行分配（求余）
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("[添加客户端的连接ID] Add ConnID = ", request.GetConnection().GetConnID(),
		" [请求消息ID]request MsgID = ", request.GetMsgID(),
		" [分配到工作池的ID]to WorkerID = ", workerID)

	//2 将消息发送给对应的worker的 TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
