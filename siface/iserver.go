package siface

/*
	定义一个服务器模块的接口
*/
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()

	//路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)
	//获取当前server 的连接管理器实例
	GetConnMgr() IConnManager

	//注册OnConnStart 创建连接之后的钩子函数方法
	SetOnConnStart(func(conn IConnection))
	//注册OnConnStop 销毁连接之前的钩子函数方法
	SetOnConnStop(func(conn IConnection))
	//调用OnConnStop 创建连接之后的钩子函数方法
	CallOnConnStart(conn IConnection)
	//调用OnConnStop 销毁连接之前的钩子函数方法
	CallOnConnStop(conn IConnection)
}