package snet

import (
	"fmt"
	"game_server_silk/siface"
	"github.com/pkg/errors"
	"io"
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

	//定义一个告知当前连接已经退出或停止 channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router siface.IRouter
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router siface.IRouter) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		Router:    router,
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
		//读取客户端连接的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("[读取客户端连接中的数据<失败>]", err)
		//	continue
		//}
		//创建一个拆包解包模块对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流,消息头限制8个字节
		headData := make([]byte, dp.GetHeadLen())
		//从连接流中读取数据，将数据读到 headData切片中（切片空间8个字节，把空间读满为止）
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("读取连接流中消息包头部数据失败：", err)
			break
		}

		//拆包，得到消息类型msgID 和 消息数据长度msgDatalen 放在msg消息模块对象中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("消息包头部拆包失败：", err)
			break
		}

		//根据dataLen消息数据长度，再次读取消息数据Data，放在msg.Data中
		//通过拆包得出消息数据长度如果大于0表示有数据内容，就接着读取消息包数据部分
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("读取连接流中消息包数据内容部分失败：", err)
				break
			}
		}
		//将读出来数据data 赋值给message模块
		msg.SetData(data)

		//将链接模块和消息模块，封装在Request模块中，变成一个request请求交给路由模块进行处理
		req := Request{
			conn: c,
			msg: msg,
		}
		//执行注册的路由方法
		go func(request siface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

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

//提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("建立的连接状态已经关闭，无法发送消息")
	}

	//将data进行封包 MsgDataLe|MsgID|Data
	dp := NewDataPack()

	//打包后返回一个已经序列化好的二进制字节类型的数据 binaryMsg
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("消息数据打包失败消息ID：", msgId)
		return errors.New("消息数据进行打包失败")
	}

	//将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("消息包ID：", msgId, "消息包发送失败", err)
		return errors.New("打包好的消息包体发送客户端失败")
	}

	return nil
}