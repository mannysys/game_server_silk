package snet

import (
	"fmt"
	"game_server_silk/siface"
	"game_server_silk/utils"
	"github.com/pkg/errors"
	"io"
	"net"
)

/*
	封装一个自定义的连接模块
*/
type Connection struct {
	//当前Conn连接模块 隶属于哪个Server
	TcpServer siface.IServer

	//当前连接的socket TCP套接字（存储当前原生的连接对象）
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//定义一个告知当前连接已经退出或停止 channel（由Reader告知Writer退出消息）
	ExitChan chan bool

	//无缓冲的管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//消息管理多路由模块（用来绑定MsgID和对应的处理业务API关系）
	MsgHandler siface.IMsgHandler
}

//初始化连接模块的方法
func NewConnection(server siface.IServer, conn *net.TCPConn, connID uint32, msgHandler siface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false, //表示当前连接是否处于关闭状态，false表示连接是开启的状态
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	//将conn连接模块实例，加入到 ConnManager连接管理器集合中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

//连接模块中的读数据的业务方法
func (c *Connection) StartReader() {
	fmt.Println("[启动读消息的协程]Reader Gortine is running")
	defer fmt.Println("[读消息的协程退出]conn Reader exit!", c.ConnID, c.RemoteAddr().String())
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
			fmt.Println("读取客户端连接中消息头部数据失败：", err)
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
			msg:  msg,
		}

		//检查工作池是否启动
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将得到客户端消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		}else {
			//根据绑定好的MsgID 找到对应处理的api 业务执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}

}

/*
	写消息Goroutine，负责发送给客户端消息的方法
 */
func (c *Connection) StartWriter() {
	fmt.Println("[启动写消息的协程]Writer Gortine is running")
	defer fmt.Println("[写消息的协程退出]conn Writer exit!", c.ConnID, c.RemoteAddr().String())

	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//接收chan管道传过来的数据，写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[写数据失败]Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			//表示Reader已经退出，此时Writer也要退出
			return
		}
	}
}

//启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("[启动一个与客户端建立连接模块] ConnID=", c.ConnID)

	//开启一个协程，去处理从当前连接中读数据的业务
	go c.StartReader()

	//启动从当前连接写数据的业务
	go c.StartWriter()

	//按照在服务端自定义的方法传递进来的 创建连接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("[关闭当前连接模块 ID]", c.ConnID)

	//如果当前连接已经关闭，直接返回
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//调用服务端注册的自定义的方法 销毁连接之前 需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()

	//告知 Writer协程关闭
	c.ExitChan <- true

	//将当前连接 从ConnMgr连接管理器集合中删除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源（关闭管道chan）
	close(c.ExitChan)
	close(c.msgChan)
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

	//将打包的数据通过msgChan管道 发送给负责写消息的协程
	c.msgChan <- binaryMsg

	return nil
}
