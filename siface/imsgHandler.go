package siface

/*
	消息管理模块抽象层
 */
type IMsgHandler interface {
	//调度执行对应的Router的消息处理方法
	DoMsgHandler(request IRequest)

	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}