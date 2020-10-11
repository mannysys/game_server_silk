package snet

import (
	"fmt"
	"game_server_silk/siface"
	"net"
)

/*
	封装一个自定义的连接模块
 */
type Connection struct {
	//当前连接的socket TCP套接字（存储当前原生的连接对象）
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//当前连接所绑定的处理业务方法API
	handleAPI siface.HandleFunc

	//定义一个告知当前连接已经退出或停止 channel
	ExitChan chan bool
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api siface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback_api,
		isClosed:  false, //表示当前连接是否处于关闭状态，false表示连接是开启的状态
		ExitChan:  make(chan bool, 1),
	}
	return c
}

//连接模块中的读数据的业务方法
func (c *Connection) StartReader() {
	fmt.Println("[读取当前协程中客户端连接中发送的数据]")
	defer fmt.Println("[客户端建立的连接ID]",c.ConnID ,"[链接读取协程退出,远程客户端的地址]", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端连接的数据到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("[读取客户端连接中的数据<失败>]", err)
			continue
		}

		//调用当前连接所绑定的HandleAPI（调用c.handleAPI属性，就是调用了赋值给handleAPI属性的方法）
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("[当前连接ID]", c.ConnID, "[当前连接绑定方法执行失败]", err)
			break
		}

	}

}

//启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("[启动一个连接模块 ID]", c.ConnID)

	//开启一个协程，去处理从当前连接中读数据的业务
	go c.StartReader()

	//TODO 启动从当前连接写数据的业务

}

//停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("[关闭当前连接模块 ID]", c.ConnID)

	//如果当前连接已经关闭，直接返回
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()

	//回收资源（关闭管道chan）
	close(c.ExitChan)
}

//获取当前原生连接的对象绑定socket conn句柄
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的 TCP状态 IP地址和port端口
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}