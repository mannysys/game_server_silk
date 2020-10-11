package siface

import "net"

/*
	定义连接模块的抽象层
 */
type IConnection interface {
	//启动连接 让当前的连接准备开始工作
	Start()

	//停止连接 结束当前连接的工作
	Stop()

	//获取当前原始连接的对象绑定socket conn句柄
	GetTCPConnection() *net.TCPConn

	//获取当前连接模块的连接ID
	GetConnID() uint32

	//获取远程客户端的 TCP状态 IP地址和port端口
	RemoteAddr() net.Addr

	//发送数据，将数据发送给远程的客户端
	Send(data []byte) error
}

//定义一个处理连接业务的方法（声明一个 HandleFunc函数，函数类型的）
type HandleFunc func(*net.TCPConn, []byte, int) error