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
	SendMsg(msgId uint32, data []byte) error
}
