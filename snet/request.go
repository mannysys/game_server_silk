package snet

import "game_server_silk/siface"

type Request struct {
	//已经和客户端建立好的连接
	conn siface.IConnection
	//客户端请求的数据
	data []byte
}

//得到当前连接模块对象
func (r *Request) GetConnection() siface.IConnection {
	return r.conn
}

//得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.data
}