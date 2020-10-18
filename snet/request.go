package snet

import "game_server_silk/siface"

type Request struct {
	//已经和客户端建立好的连接
	conn siface.IConnection

	//客户端请求的数据，封装到message模块中
	msg siface.IMessage
}

//得到当前连接模块对象
func (r *Request) GetConnection() siface.IConnection {
	return r.conn
}

//得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取到消息包类型ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}